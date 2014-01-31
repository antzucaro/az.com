---
title: "Last.fm and MOC"
date: "2011-02-02"
categories: 
  - "blog"
  - "geekery"
---
I'm a big fan of <a href="http://www.last.fm">Last.fm</a>. For those not aware, it is a music recommendation service. It does its recommendations based on your own listening patterns; it watches the tracks you play on your computer and compares it to those played by other people sharing common artists. As an example, Last.fm recommends that I check out the band Catch 22 (who I do not have in my music library) based on the number of times I've played songs by Less Than Jake and because other people who have listened to Less Than Jake listen to Catch 22. Pretty cool.

One thing that is critical to the entire Last.fm experience is the connection between the music player you use and the Last.fm servers that make the recommendations. If you don't have that connection, you won't get any recommendations! Forming the bridge between the two is something called a <em>scrobbler; </em>the scrobbler is the piece of software that watches what you're playing on your music player and sends that information over to the recommendation servers. Last.fm provides a scrobbler client for both Windows and Mac as a free download. Although not officially supported, Linux users can still participate with plugins for many popular players (Amarok, Rhythmbox, Audacious).

Even Linux music players without ready-made plugins can use the service through a software package called <a href="http://www.red-bean.com/decklin/lastfmsubmitd/">lastfmsubmitd</a>. This software package sets up a generic way for players to send data over to Last.fm via the standard client/server model familiar to most Linux users. Here's how I set it up with MOC, a popular command-line music player.
<h3>Get the lastfmsubmitd package</h3>
I run Ubuntu on my home theater PC, which is where I play most of my music at home. To install lastfmsubmitd there I ran the following:

<pre class="brush:bash;">sudo apt-get install lastfmsubmitd</pre>

During the installation you'll be asked several straightforward questions (username, password), after which you'll be left with a running lastfmsubmitd daemon. Check to make sure it is up and running with:

<pre class="brush:bash;">ps -ef | grep lastfmsubmit[d]</pre>

You should see a process running under that name. If not, the following should get you up and running:

<pre class="brush:bash;">sudo /etc/init.d/lastfmsubmitd start</pre>

The next step is hooking up your player to submit information to lastfmsubmitd.
<h3>Change your MOC configuration</h3>
MOC is my music player of choice, and unfortunately it is one that doesn't have a prebuilt plugin for Last.fm. Instead I have to add the following in my ~/.moc/config file to make MOC submit the songs I play to lastfmsubmitd:

<pre class="brush:bash;">OnSongChange = &quot;/usr/lib/lastfmsubmitd/lastfmsubmit --artist %a --title %t --length %d --album %r&quot;</pre>

What this does is submit each song to lastfmsubmitd when you start playing it. lastfmsubmitd then submits the song data to the actual Last.fm servers using its own built in submission client (that's the /usr/lib/lastfmsubmitd/lastfmsubmit part). Take note that this method sends data <em>right when the song changes</em>. If you are a habitual song-changer and don't feel like spamming Last.fm with songs that you really haven't listened to fully, check out Luke Plant's submission script <a href="http://lukeplant.me.uk/blog/posts/moc-and-last-fm/">here</a>, as it may suit you better. I rarely skip around in songs, so I just stick with the default submission client (the one in the OnSongChange config above) .
<h3>Restart and Verify!</h3>
You have to completely restart MOC before the configuration option goes into effect. Take care to kill the process if it is still running after you've exited to the terminal, as MOC is a server process. Once you've verified that the process is gone, fire up 'mocp' again and play a song or two. Check your Last.fm profile page shortly after you start the song. If you see it showing up in the list, you're in business!

<div class="wp-caption aligncenter" style="width: 630px"><a href="/uploads/2011/02/verified.png"><img class="size-full wp-image-650" title="Verified!" src="/uploads/2011/02/verified.png" alt="Verified!" width="630" height="435" /></a><p class="wp-caption-text">Verified!</p></div>
