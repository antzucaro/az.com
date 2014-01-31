---
title: "2013 - The Year in Music"
date: "2014-01-02"
categories: 
  - "geekery"
---
Every December Kristin and I talk about our favorite albums and songs for the year, and it always leads to great conversations. We debate the merits of the artists, or how that one song really annoyed us in an otherwise spectacular album, or how that one band member has really lost their touch. You get the idea. 2013 was no exception to this, and we recently had our music chat two days before I had to go back to work after the holiday. After talking for a while, it occurred to me that our scrobbled music data was probably sitting on the net, ready to be analyzed. With the spare time I had left before getting back to business, it would be pretty cool to feed those discussions with some hard data to back up our feelings!

We'd been using Last.fm for quite a while, so I reckoned that we'd be able to get it back out easily from the service. What I needed was the artist, album, song, and timestamp for all the songs we listened to during 2013. Having this information would enable me to break down things by intervals of time quite nicely. A quick check on the Last.fm API page confirmed my suspicion: the [getRecentTracks](http://www.last.fm/api/show/user.getRecentTracks) API endpoint had all the pieces of info I was looking for (side note: if only all companies were like this!). I readied myself to fetch the data so I could fuel some more of those interesting discussions!

After signing up for an API key, I checked out the example URL to see what format I'd have to parse. I had two options: XML and JSON. Having a good enough familiarity with JSON already, I decided I'd try out XML to see how difficult it was. I figured that it was a good use of the [golang](http://golang.org/) standard library either way. I whipped up a quick and dirty script to fetch and save the data, fueled by simple command line options which mapped straight into the API's parameters. I ended up with the small [Golang program](#go-script) at the bottom of this post. Also included at the bottom are the corresponding [steps](#steps) on how to use the data file once you have it.

Having fetched my data using this program, I next imported it to PostgreSQL so I could query it until my heart was content! Below are some of the questions that I answered with my new-found knowledge.

### Which artists did we listen to the most?

1. Daft Punk (468)
2. Nine Inch Nails (342)
3. Tame Impala (239)
4. Carbon Leaf (209)
5. How to Destroy Angels (171)
6. Dave Matthews Band (148)
7. Pink Floyd (144)
8. Cliff Martinez (138)
9. Trent Reznor and Atticus Ross (128)
10. Pearl Jam (125)

A few surprises here. We discovered (har har) Daft Punk's <i>Discovery</i> album pretty early on and we didn't let up on it all year, so that was expected. Nine Inch Nails has been a staple also, so nothing there either. What legitimately surprised me was how much Tame Impala showed up. I checked them out when <i>Lonerism</i> came out and was hooked. I guess I got really hooked! That's a lot of plays for only having two albums. Here's to them creating more awesome music in the future!

Cliff Martinez was another pleasant surprise. At the beginning of the year we watched <i>Contagion</i> and really admired the score. I grabbed it and have been enjoying its goodness ever since. I suppose I really do prefer that soundtrack for hacking more than anything else, which most likely explains the number of plays.

### How many different artists did we listen to?

We listened to 591 artists throughout the year. Whoa! This is a bit skewed since many of them come from SomaFM, Spotify, or Last.fm's streaming capabilities. Still cool, though.

### How many tracks did we listen to per month?

* January - 774
* February - 444
* March - 951
* April - 362
* May - 651
* June - 586
* July - 931
* August - 560
* September - 810
* October - 597
* November - 656
* December - 851

I can only speculate about most of the numbers here. We went on vacation in August, so that explains the lower number of plays. That vacation followed a month where I was in the office for over 55 hours a week. Since I listen to a lot of music while in the office, it follows that the totals would be high. Closing out the year, December saw us listening to the Christmas channel on SomaFM a lot. That plus a lot of time off meant more time for music! Plus, cold weather meant we were inside more often than not.

### Which artists did we listen to the most by month?

* January: Daft Punk (61)
* February: Cliff Martinez (44)
* March: Carbon Leaf (80)
* April: Dave Matthews Band (46)
* May: Daft Punk (121)
* June: The Beatles (37)
* July: Pearl Jam (45)
* August: Nine Inch Nails (83)
* September: Tame Impala (41)
* October: Dave Matthews Band (38)
* November: Trent Reznor and Atticus Ross (35)
* December: Nine Inch Nails (45)

These numbers show how attached we get to new album releases. Most of the top hits are from album releases of the previous month. Carbon Leaf's <i>Ghost Dragon Attacks Castle</i> came out in Feb and Ninch Inch Nails' <i>Hesitation Marks</i> came out in November. In each case we listened to the physical CDs for a while, then I finally got around to putting the music on our server where it was much easier to access (thus more plays). Daft Punk's <i>Random Access Memories</i> came out in May, but I got the digital copy and immediately transferred it for our listening enjoyment.  

I'm pretty happy with this data, and the fact that I could get it back out of Last.fm pretty easily. Now that we're into 2014 and I am graced with a Spotify subscription, I'll try to be more conscious of diversity for the next 12 months. When I do this process next December, I want to see many more artists. Bring on the music!

<a id="go-script"><script src="https://gist.github.com/antzucaro/8212369.js?file=summarizr.go"></script></a>

<a id="steps"><script src="https://gist.github.com/antzucaro/8212369.js?file=steps.txt"></script></a>
