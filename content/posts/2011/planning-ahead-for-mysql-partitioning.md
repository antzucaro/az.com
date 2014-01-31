---
title: "Planning ahead for MySQL partitioning"
date: "2011-01-30"
categories: 
  - "blog"
  - "geekery"
---
<h3>The Scenario</h3>
Let's suppose you are modeling a database in MySQL 5.5. Based on your knowledge of the requirements for the application <em>behind </em>this database you can guess which of the tables within it are going to be huge. For those huge tables you want to plan ahead to make sure that you can scale out your data such that you have less headaches in the future. Here's the catch: you want to keep it simple until it comes time to scale out. That sounds like an honorable goal, right? Now how exactly do we do that...hmm. This post is how I would approach this problem. I hope you can gain something from it. I certainly learned some things in trying to get this stuff to work!

To ground this problem scenario into something concrete let's assume that we're modeling a simple relationship - players and games in a multiplayer game. Here's the model that we're after:

<div class="wp-caption aligncenter" style="width: 559px"><span style="font-size: x-small;"><a href="/uploads/2011/01/normal_relationship.png"><img class="size-full wp-image-613" title="Multiplayer game ER diagram" src="/uploads/2011/01/normal_relationship.png" alt="Multiplayer game ER diagram" width="559" height="149" /></a></span><p class="wp-caption-text">Multiplayer game ER diagram</p></div>

Here is the DDL to create this model:

<pre class="brush:sql;">
CREATE  TABLE IF NOT EXISTS `parttest`.`player` (
`player_id` INT NOT NULL ,
`name` VARCHAR(30) NULL ,
PRIMARY KEY (`player_id`) )
ENGINE = InnoDB;

CREATE  TABLE IF NOT EXISTS `parttest`.`game` (
`game_id` INT NOT NULL ,
`duration` TIME NULL ,
`start_dt` TIMESTAMP NOT NULL ,
PRIMARY KEY (`game_id`) )
ENGINE = InnoDB;

CREATE  TABLE IF NOT EXISTS `parttest`.`player_game` (
`player_id` INT NOT NULL ,
`game_id` INT NOT NULL ,
`create_dt` TIMESTAMP NOT NULL ,
PRIMARY KEY (`player_id`, `game_id`) ,
INDEX `player_game_fk01` (`player_id` ASC) ,
INDEX `player_game_fk02` (`game_id` ASC) ,
CONSTRAINT `player_game_fk01`
FOREIGN KEY (`player_id` )
REFERENCES `parttest`.`player` (`player_id` )
ON DELETE NO ACTION
ON UPDATE NO ACTION,
CONSTRAINT `player_game_fk02`
FOREIGN KEY (`game_id` )
REFERENCES `parttest`.`game` (`game_id` )
ON DELETE NO ACTION
ON UPDATE NO ACTION)
ENGINE = InnoDB;
</pre>

In this model we have players and games. Many players can participate in a single game and many games can be played by a single player, thus the two tables "player" and "game" have a many-to-many relationship. To represent that in the database we have an associative table ("player_game") between the two. Also included in two of the tables are columns that track when something happened.

For the sake of this example let's assume that there won't be too many players, but those limited amount of players will play the game feverishly (it's an addictive open source game with a small community). That means that the "game" table and the associative "player_game" table will be pretty large as time goes by - those are the two tables that we'll want to make scalable for the future. We'll want the option to easily <em>partition </em>them at some point down the line.
<h3>Enter Partitioning</h3>
Wait, wait, wait. <em>Partitioning</em>? What the heck is that for? I won't bore you with too many details (you can find out more <a href="http://en.wikipedia.org/wiki/Partition_(database)">here</a> and <a href="http://dev.mysql.com/doc/refman/5.5/en/partitioning.html">here</a>), so let's just say that partitioning is a way to split up your database tables in to smaller pieces so that operations upon them can be faster (in most situations that is the goal, at least). Partitions provide the following benefits (this is some, not all):
<ul>
	<li>Maintenance - when you split up 	tables into partitions you can often do maintenance on them such 	that you only affect a small amount of your users at a time. Having 	partitions also allows you to easily implement a pruning or 	archiving strategies depending on how you split the data. If you 	split up your tables based upon create timestamp, for example, you 	can easily drop or move older data (the ones that aren't queried as 	often as the newer ones) via a singleo command.</li>
	<li>Query performance - running 	queries against large datasets most often causes problems. When the 	data is partitioned the database can â€œautomagicallyâ€ determine 	which partitions aren't involved in queries; it can tell which 	partitions aren't involved in the WHERE clause of your SQL 	statements, thus it doesn't need to scan them. That means you 	traverse less data to get the information your query needs, which in 	turn equals faster performance! Of course this only applies if the 	WHERE clause of the queries in question can avoid partitions. If you 	are grabbing everything from a table partitioning obviously won't 	help you.</li>
	<li>Availability - partitions can be 	split onto different physical media, which means your table is still 	partially available even if one of those media sources goes down. 	This also has impacts with respect to maintenance (see above).</li>
</ul>
<h3>Implemention (some hurdles to jump)</h3>
Back to the example. Since we know that the "game" and "player_game" tables are going to be big, let's see how to set ourselves up for easy partitioning in the future. The goal here is to get to a partitioned table with a minimum number of ALTER TABLE statements. Ideally we'd have only one ALTER TABLE - to add partitioning! Let's project ourselves into the future to the point in time when our data is so large that we are forced to consider partitioning. Fortunately for us our tables have no data in them, which makes our tests quick and relatively painless.

Given the initial DDL, the first step would be to try the obvious: alter the table to add partitioning! Let's give that a try:

<pre class="brush:sql;">ALTER TABLE game

PARTITION BY RANGE(unix_timestamp(start_dt)) (

PARTITION part0 VALUES LESS THAN (unix_timestamp('2011-01-01 00:00:00')),

PARTITION part1 VALUES LESS THAN (unix_timestamp('2011-01-02 00:00:00')),

PARTITION part2 VALUES LESS THAN (unix_timestamp('2011-01-03 00:00:00')),

PARTITION part3 VALUES LESS THAN (unix_timestamp('2011-01-04 00:00:00')),

PARTITION part4 VALUES LESS THAN (unix_timestamp('2011-01-05 00:00:00')));

</pre>

Okay, that failed with the following. Bummer:

[code]ERROR 1217 (23000): Cannot delete or update a parent row: a foreign key constraint fails[/code]

The only table in our model with foreign keys is "player_game," so let's remove those constraints to move forward. Let's do that:

<pre class="brush:sql;">alter table player_game drop foreign key player_game_fk01;
alter table player_game drop foreign key player_game_fk02;</pre>

We can now try to partition again. We run the same alter as before and get the following:

[code]ERROR 1503 (HY000): A PRIMARY KEY must include all columns in the table's partitioning function[/code]

Argh! What in the world does this mean? Looking up this error leads us to <a href="http://dev.mysql.com/doc/refman/5.5/en/partitioning-limitations-partitioning-keys-unique-keys.html">this entry</a> in the manual which tells us why. In short: you can't partition a table using a column that is not included in the primary key for that table. Since the "start_dt" column isn't in our primary key, we couldn't proceed with our alter. Okay, so let's add it:

<pre class="brush:sql;">ALTER TABLE game DROP PRIMARY KEY;

ALTER TABLE game ADD PRIMARY KEY (game_id, start_dt);</pre>

Once more with the alter to partition, as before. Cross your fingers...

[code]Query OK, 0 rows affected (0.05 sec)[/code]

Success! Now all we need to do is the same ALTERs to the player_game table (using the same process we just did) and add back those foreign keys that we dropped earlier. We'll then have our model looking like before but with the added benefit of 5 partitions each for game and player_game.

But wait...

<pre class="brush:sql;">ALTER TABLE player_game ADD CONSTRAINT player_game_fk01 FOREIGN KEY (game_id) REFERENCES game (game_id);

ERROR 1005 (HY000): Can't create table 'parttest.#sql-52b_3e' (errno: 150)</pre>

Again - argh! What is wrong now? Going back to the limitations link from the last error gives us the answer: partitioned tables can't have foreign keys. Yep, you heard me correctly - <strong>partitioned tables can't be involved in any foreign keys. </strong>Really, who needs referential integrity anyway? Just kidding, as this <em>really</em> <em>big </em>limitation in my opinion.

Seeing as we've already dropped the foreign keys of the player_game table to get to this point, we don't have anything further to do. To get where we are we had to remove the foreign keys from the tables that we'd be partitioning, and we had to move the column we were partitioning on into the primary key.
<h3>My recommendations</h3>
What does all this mean in terms of how to design your tables now, when they aren't in a partitioned state? Here's my recommendations based on what I've seen with partitioning thus far:
<ul>
	<li>Build the column you will 	be partitioning on in to 	the primary key from the get-go. Make 	sure the ordering of the primary key (if you will have multiple 	columns at all) makes sense for your query load on the table. Put 	the columns most used in WHERE clauses first.</li>
	<li>As an alternative to the 	above you can keep your primary key the same but add 	a unique key including all the columns of the primary key plus the 	column you'll partition on. This will help you avoid duplicate 	entries when you want to bring the partitioning column into the 	primary key, but you'll have a few additional steps later (dropping 	and recreate the primary key).</li>
	<li>Keep the foreign keys in 	the database in place as 	long as you can while it is still available, but plan ahead in the 	application by coding as if they weren't there. 	When the time comes to 	partition you can then be confident that all the data up to that 	point is clean in terms of its referential integrity. You'll also 	have debugged your application to flush out the errors by that time 	(hopefully).</li>
</ul>
So those are my recommendations for dealing with MySQL 5.5's partitioning. Obviously I think they can do a lot better with their implementation. Ideally I'd like to have the ability to alter the table and add the partitions regardless of its referential integrity or primary keys. The fact that I have to sacrifice the integrity of the database tables (one of the key strong points of having a relational database in the first place) makes me very hesitant to use it at all.
<h3>Some Caveats</h3>
I'm an Oracle DBA, so take this advice with a grain of salt. The stuff in this post is just one solution to the "big data" problem. Furthermore this is just how I would handle a partitioning situation given the example provided. I'd love to hear some people with more experience chime in on what they do in similar situations.
