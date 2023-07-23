+++ 
draft = false
date = 2023-03-11T13:15:00-04:00
title = "Checking out DuckDB"
description = ""
slug = ""
authors = []
tags = []
categories = []
externalLink = ""
series = []
+++


I came across [DuckDB][DuckDB] the other day on Hacker News. It's like SQLite but for OLAP use cases, which made me think about the "[by the numbers][by-the-numbers]" blog posts I used to do for Xonotic. These involve connecting to a PostgreSQL database and running expensive analytical queries via Python's pandas library, eventually rendering them to matplotlib charts that people can _ooo_ and _aah_ over. DuckDB appears to be a good fit for at least the first portion of that (the querying bits), so let's see how it handles a reasonable volume of data.

I'll start with taking a couple of the bigger tables from [XonStatDB] and bringing them into DuckDB for a quick-and-dirty performance comparison. Note that I'm not out to make any categorical statements or even recommendations here. It's purely to get a feel for DuckDB and what place it could occupy in my toolbox, if any.

I happen to have a copy of the XonStat PostgreSQL database locally from which all the data can be sourced. It can be copied directly from Postgres using [DuckDB's scanner][DuckDB Scanner] interface, so let's do that. This helps avoid an extra step of converting it to some other format like CSV as an intermediary step. The more steps involved, the more annoyed I'm likely to get with the process, so it's nice that we can reduce steps here.

First we'll create the destination database and with it the empty tables where eventually all the data will be copied. The schema of the table is almost a carbon copy of the DDL from postgres:

```
CREATE TABLE games
(
    game_id BIGINT NOT NULL,
    start_dt TIMESTAMP NOT NULL,
    game_type_cd TEXT NOT NULL,
    server_id INT NOT NULL,
    map_id INT NOT NULL,
    duration BIGINT,
    winner INT,
    create_dt TIMESTAMP NOT NULL,
    match_id TEXT
);


CREATE TABLE player_weapon_stats
(
  player_weapon_stats_id BIGINT NOT NULL,
  player_id INT NOT NULL,
  game_id INT NOT NULL,
  player_game_stat_id BIGINT NOT NULL,
  weapon_cd TEXT NOT NULL,
  actual INT NOT NULL default 0,
  max INT NOT NULL default 0,
  hit INT NOT NULL default 0,
  fired INT NOT NULL default 0,
  frags INT NOT NULL default 0,
  create_dt TIMESTAMP
);
``` 

Creating the DuckDB database with these tables in it is simple:

```
duckdb weaponstats.db < schema.sql
```

To confirm the tables were created successfully and weren't silently ignored (hey, it could happen), we can open up the DB using the CLI and ask it to print out the schema:

```
$ duckdb weaponstats.db 
v0.7.1 b00b93f0b1
Enter ".help" for usage hints.
D .schema
CREATE TABLE games(game_id BIGINT NOT NULL, start_dt TIMESTAMP NOT NULL, game_type_cd VARCHAR NOT NULL, server_id INTEGER NOT NULL, map_id INTEGER NOT NULL, duration BIGINT, winner INTEGER, create_dt TIMESTAMP NOT NULL, match_id VARCHAR);;
CREATE TABLE player_weapon_stats(player_weapon_stats_id BIGINT NOT NULL, player_id INTEGER NOT NULL, game_id INTEGER NOT NULL, player_game_stat_id BIGINT NOT NULL, weapon_cd VARCHAR NOT NULL, actual INTEGER NOT NULL DEFAULT(0), max INTEGER NOT NULL DEFAULT(0), hit INTEGER NOT NULL DEFAULT(0), fired INTEGER NOT NULL DEFAULT(0), frags INTEGER NOT NULL DEFAULT(0), create_dt TIMESTAMP);;

```

Next we'll load up the scanner which allows us to "see" into Postgres. Punching in the commands from the documentation:

```
D INSTALL postgres;
D LOAD postgres;
```

We then attach to postgres and tell DuckDB which schema to inspect. It will create views for each table it finds, with the view names matching the table names in the attached database.

```
CALL POSTGRES_ATTACH('dbname=xonstatdb user=xonstat host=127.0.0.1 password=some_password', source_schema='xonstat');
```

The nice thing is that after successful attachment, you can query the tables as though they are native to DuckDB. This feels natural and leads to a nice developer experience:

```
D .timer on
D select count(*) from player_weapon_stats;
100% ▕████████████████████████████████████████████████████████████▏ 
┌──────────────┐
│ count_star() │
│    int64     │
├──────────────┤
│     33964270 │
└──────────────┘
Run Time (s): real 7.619 user 10.654525 sys 0.181006

```

You can also pick and choose which tables to scan using POSTGRES_SCAN, which does **not** create views for each table in the attached database. This is useful when you want to bring in tables to DuckDB with their original names. Otherwise you'll need to name your DuckDB tables with some suffix to ensure uniqueness, adding that detail each time you want to run the same queries between the two systems. I'm lazy and would prefer to copy and paste my queries, so I'll opt for scanning table-by-table instead of attaching. That approach looks as follows:

```
SELECT * FROM POSTGRES_SCAN('dbname=xonstatdb user=xonstat host=127.0.0.1 password=some_password', 'xonstat', 'games');
```

All of this looking into PostgreSQL through the window of DuckDB is nice, but ultimately we don't want **just** a pass-through to Postgres. That doesn't change the format of the data at rest, which means it's still tuple (row) based on disk. Why query that format through DuckDB if we can simply query Postgres directly with its own tooling? What we really want to see is how things change when we actually import the data into DuckDB's local database and thus into its columnar disk format. That's where we can really get a feel for whether or not our efforts are worthwhile. 

To do that, we can simply `INSERT..SELECT` from attached PostgreSQL tables into real tables in DuckDB using the same scanner interface mentioned above. We need only a couple of modifications due to data type differences (the interval data type in particular). Here's what that looks like:

```
.timer on


INSERT INTO games
SELECT game_id, start_dt, game_type_cd, server_id, map_id, extract('epoch' from duration), winner, create_dt, match_id
FROM POSTGRES_SCAN('dbname=xonstatdb user=xonstat host=127.0.0.1 password=xonstat', 'xonstat', 'games');


INSERT INTO player_weapon_stats
SELECT * FROM POSTGRES_SCAN('dbname=xonstatdb user=xonstat host=127.0.0.1 password=xonstat', 'xonstat', 'player_weapon_stats');
```

With timing on we can see how long this takes for each table. The player_weapon_stats table is the big one, so it is no surprise it takes the longest:

```
Run Time (s): real 2.632 user 2.161559 sys 0.131486
Run Time (s): real 37.319 user 69.627807 sys 2.154641
```

With that complete all the data is now available locally and is packed into a DuckDB database - a single file on the filesystem. How portable! We can now get to the fun stuff of comparing queries.

Running a few queries between the two, it's really striking how fast DuckDB is, especially when considering how there are no indexes defined on it. 

First, a simple query to count the number of records:

```
-- PostgreSQL
xonstatdb=# select count(*) from player_weapon_stats;
  count   
----------
 33964270
(1 row)

Time: 1558.775 ms (00:01.559)


-- DuckDB
D select count(*) from player_weapon_stats;
┌──────────────┐
│ count_star() │
│    int64     │
├──────────────┤
│     33964270 │
└──────────────┘
Run Time (s): real 0.016 user 0.053499 sys 0.000000

```

Next, a simple aggregate over values within the table:

```
-- PostgreSQL
xonstatdb=# select weapon_cd, count(*)
xonstatdb-# from player_weapon_stats
xonstatdb-# group by weapon_cd
xonstatdb-# order by count(*) desc;
<output snipped>
Time: 3823.059 ms (00:03.823)


-- DuckDB
D select weapon_cd, count(*)   
> from player_weapon_stats     
> group by weapon_cd           
> order by count(*) desc; 
<output snipped>
Run Time (s): real 0.426 user 1.646932 sys 0.023292

```

How about something with joins and aggregations? Here's where things start to get bonkers:

```
-- PostgreSQL
xonstatdb=# select game_type_cd, count(distinct games.game_id), sum(frags), sum(frags)/count(distinct games.game_id)
from games join player_weapon_stats on games.game_id = player_weapon_stats.game_id
group by game_type_cd order by 4 desc;
 game_type_cd | count  |   sum    | ?column? 
--------------+--------+----------+----------
 ctf          | 721001 | 73749208 |      102
 ka           |   6881 |   481208 |       69
 dm           | 820513 | 52134650 |       63
 tdm          |  71499 |  4427671 |       61
 kh           |   9311 |   567298 |       60
 freezetag    |     22 |     1306 |       59
 rune         |    118 |     6268 |       53
 ft           |  28321 |  1417903 |       50
 ca           |  69217 |  3451420 |       49
 ons          |     63 |     2619 |       41
 as           |   4510 |   162361 |       36
 dom          |   4421 |   159426 |       36
 nexball      |      4 |      131 |       32
 duel         | 198317 |  5613840 |       28
 nb           |   1173 |    26124 |       22
 arena        |     12 |       54 |        4
(16 rows)

Time: 189633.479 ms (03:09.633)


-- DuckDB
D select game_type_cd, count(distinct games.game_id) game_count, sum(frags) total_frags, sum(frags)/count(distinct games.game_id) frags_per_game from games join player_weapon_stats on games
┌──────────────┬────────────┬─────────────┬────────────────┐
│ game_type_cd │ game_count │ total_frags │ frags_per_game │
│   varchar    │   int64    │   int128    │     int128     │
├──────────────┼────────────┼─────────────┼────────────────┤
│ ctf          │     721001 │    73749208 │            102 │
│ ka           │       6881 │      481208 │             69 │
│ dm           │     820513 │    52134650 │             63 │
│ tdm          │      71499 │     4427671 │             61 │
│ kh           │       9311 │      567298 │             60 │
│ freezetag    │         22 │        1306 │             59 │
│ rune         │        118 │        6268 │             53 │
│ ft           │      28321 │     1417903 │             50 │
│ ca           │      69217 │     3451420 │             49 │
│ ons          │         63 │        2619 │             41 │
│ dom          │       4421 │      159426 │             36 │
│ as           │       4510 │      162361 │             36 │
│ nexball      │          4 │         131 │             32 │
│ duel         │     198317 │     5613840 │             28 │
│ nb           │       1173 │       26124 │             22 │
│ arena        │         12 │          54 │              4 │
├──────────────┴────────────┴─────────────┴────────────────┤
│ 16 rows                                        4 columns │
└──────────────────────────────────────────────────────────┘
Run Time (s): real 1.797 user 6.853171 sys 0.091531

```

This last case is eye-opening for me: what took DuckDB **under two seconds** took **over three minutes** for Postgres. It really drives home just how performant a columnar database can be when running typical analytical workloads like this. That the tool doesn't have any heavy installation procedure or service that requires care and feeding is nice. It's like icing on an already nice cake. You know...the good kind, not the so-sweet-it-hurts kind. I'll certainly keep DuckDB in mind whenever I'm doing similar analytical tasks in the future.

[by-the-numbers]: https://xonotic.org/posts/2018/2017-by-the-numbers/
[DuckDB]: https://duckdb.org/
[DuckDB Scanner]: https://duckdb.org/docs/extensions/postgres_scanner
[XonStatDB]: https://gitlab.com/xonotic/xonstatdb/-/blob/master/tables/player_weapon_stats.tab
