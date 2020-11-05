---
title: "CDDB Woes"
date: 2020-11-04T19:55:00-04:00
---

It's 2020 and I'm still buying CDs. Today I sat down to rip some newly bought discs from Goodwill. Breaking out my tried-and-true `abcde` (which [I wrote][abcde] about exactly ten years ago, whoa!) led to the following, though:

```
$ abcde
Grabbing entire CD - tracks: 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16
CDDB unavailable.
---- Unknown Artist / Unknown Album ----
1: Track 1
2: Track 2
3: Track 3
4: Track 4
5: Track 5
6: Track 6
7: Track 7
8: Track 8
9: Track 9
10: Track 10
11: Track 11
12: Track 12
13: Track 13
14: Track 14
15: Track 15
16: Track 16

Edit selected CDDB data [Y/n]? 
```

No CDDB available? What gives? This is kind of a problem, since the [tagging and cataloging system][beets] I use needs at least a small amount of information with which to fill in gaps. The `abcde` utility is what provides this information, and without it I'd only be able to utilize folder names to find music. Boo.

A couple of Google searches later led me to the underlying [reason][freedb]: the FreeDB compatibility API run by MusicBrainz shut down in March 2019. Users were instructed to utilize the *real* API to get the same information instead of this one, which was only intended to help users transition from a defunct one to begin with. In other words, it was a temporary solution and now that temporary time period was up.

Okay, no problem. ABCDE has support for MusicBrainz lookups directly. Let's give that a try:

```
# ~/.abcde.conf
CDDBMETHOD=musicbrainz
```

```
$ abcde
dasd: 1 16  150 187112 0
11212
19110
33122
47902
59057
69192
80295
91635
106160
120437
142710
152615
163525
168217
183207 
Grabbing entire CD - tracks: 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16
No Musicbrainz match.
---- Unknown Artist / Unknown Album ----
1: Track 1
2: Track 2
3: Track 3
4: Track 4
5: Track 5
6: Track 6
7: Track 7
8: Track 8
9: Track 9
10: Track 10
11: Track 11
12: Track 12
13: Track 13
14: Track 14
15: Track 15
16: Track 16

Edit selected CDDB data [y/N]?
```

Wait, what? No MusicBrainz match either? Hmm. Looking up some information on that, it seems that my version of ABCDE was using this API to fetch information about discs (substitute $DISCID with the actual value obtained from [`cd-discid`][cd-discid] utility):

`http://musicbrainz.org/ws/1/release/?type=xml&discid=$DISCID`

Navigating to that URL and putting in a random discid from my collection tells me why **that** API isn't working: the version being used [is deprecated][api deprecation], and all calls to the API just respond back with the URL to that blog post as the payload. There's a newer version, of course, but ABCDE doesn't speak it. Strike two.

Backing up a bit I came to find out that there are alternate servers for fetching FreeDB information. [One in particular][gnudb] looked promising. Substituting its information into ABCDE finally (thankfully) yielded success:

```
# ~/.abcde.conf
CDDBMETHOD=cddb
CDDBURL="http://gnudb.gnudb.org:80/~cddb/cddb.cgi"

$ abcde
Grabbing entire CD - tracks: 01 02 03 04 05 06 07 08 09 10 11 12 13
Which entry would you like abcde to use (0 for none)? [0-3]: 1
Selected: #1 (Alanis Morissette / Jagged Little Pill)
---- Alanis Morissette / Jagged Little Pill ----
Year: 1995
Genre: Rock
1: All I Really Want
2: You Oughta Know
3: Perfect
4: Hand In My Pocket
5: Right Through You
6: Forgiven
7: You Learn
8: Head Over Feet
9: Mary Jane
10: Ironic
11: Not The Doctor
12: Wake Up
13: You Oughta Know [Alternate Take] & Your House [A Capella-Version] [Hidden Track]
``` 

Ahhhh, success. Finally I can get to actually cataloging the stuff in digital form. Thank you so much, GnuDb! Keep up the great work.

[abcde]: https://antzucaro.com/posts/2010/abcde-a-better-cd-encoder/
[api deprecation]: https://blog.metabrainz.org/2018/02/01/web-service-ver-1-0-ws-1-will-be-removed-in-6-months/
[beets]: https://beets.io/
[cd-discid]: http://manpages.ubuntu.com/manpages/trusty/man1/cd-discid.1.html
[freedb]: https://blog.metabrainz.org/2018/09/18/freedb-gateway-end-of-life-notice-march-18-2019/
[gnudb]: https://gnudb.org/index.php
