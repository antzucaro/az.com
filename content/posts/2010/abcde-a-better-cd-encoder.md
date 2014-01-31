---
title: "ABCDE - A Better CD Encoder"
date: "2010-11-04"
categories: 
  - "blog"
  - "geekery"
---
I'm a big fan of automating things. In almost all aspects of my life you can find an example of me trying to free up some of my time by eliminating the mundane, repeatable steps with some automated process (think auto-billpay and programmable coffee makers). This is even more evident in my computer work - I don't want to be typing in the same commands over and over! For one, all of that typing takes time, and time is something of which I have too little. For two, me having to manually type or click things out each time exposes whatever I'm doing to human error.

In this post I'd like to talk about one example of automation that I've used to save a lot of time: archiving of my CDs to MP3. I did this about two summers back when I was trying to archive my entire collection to an external hard drive I'd just bought (over 300 CDs between myself and my wife). My goal was to rip them all with:
<ul>
	<li>The same level of quality (the highest quality variable bitrate MP3)</li>
	<li>The same folder structure (&lt;Artist&gt;/&lt;Album&gt;/&lt;Artist&gt; - &lt;Track Number&gt; - &lt;Track Name&gt;.mp3)</li>
	<li>Meaningful meta-info (IDv3 tags)</li>
	<li>Playlists in the album folder (M3U files)</li>
</ul>
I started with a program called <a href="http://en.wikipedia.org/wiki/K3b">K3B</a> that was more than adequate in terms of my goals above, but it just took too long. After about 10 CDs I realized that the amount time to rip this way wasn't feasible (even with two CD drives to work with in the tower I was using). To Google I went, eventually finding a neat little program called ABCDE that seemed to fit the bill.
<h2>Enter ABCDE</h2>
<a href="http://code.google.com/p/abcde/">ABCDE</a> is a deceptively simple <a href="http://en.wikipedia.org/wiki/Bash_%28Unix_shell%29">Bash</a> script that packs a whole host of features: Mp3/FLAC/Ogg ripping, playlists, CDDB lookups, and file tagging. What made me choose it over everything else wasn't its feature set per se, but rather that it can be customized and simplified such that you don't have to type out your preferences every time; you essentially set up your preferences in a configuration file <em><strong>once</strong></em>, then forget about them! In the end I had my setup customized so that all I needed to do to rip my CDs was one of the following commands (depending on which drive the CD was loaded, of course):

<pre class="brush:bash;">
ripcd
ripdvd
</pre>

How did I get to that point, you ask? Configuration, of course! Read on...
<h2>Configuring ABCDE</h2>
After you install it (which I won't cover here), ABCDE stores its default configuration in the appropriately named <em>/etc/abcde.conf</em>.  I started off by editing this file to make sure that it did things exactly the way I wanted (see my goals above).  Here's an annotated guide of what I changed in the default setup.
<h3>Allow track padding</h3>
The visual presentation of my music files is important to me. This includes how they are sorted when they are sitting in their respective folders - I wanted them to be in the folder in the same order as they are in the album! By default ABCDE doesn't pad their track numbers (should you put them in the filename), which leads to track number "10" being shown before track number "1" if the track number is the first significant difference to the file names. Padding the track names solves this issue (it makes track 1 appear with "01" in the filename) and provides me a consistent representation in my file browser. To enable this, make this change in your config:

<pre class="brush:bash;">PADTRACKS=y</pre>
<h3>MP3 Encoding Parameters</h3>
To get the highest quality encoding out of ABCDE I have the following LAME option (LAME is what ABCDE uses to create MP3s):

<pre class="brush:bash;">LAMEOPTS=&quot;-b 320&quot;</pre>

Let me say here that encoding parameters are extremely subjective. Some folks want lossless,  while others want huge space savings. What I want is the best possible  quality while staying within the constraints of the MP3 file format. I  do that because I think lossless (FLAC) just takes up too much space. I'm also lazy in that I don't  want to convert to MP3 if I have it archived in a different format. If you are interested in exploring other MP3 options for ABCDE here, check out the <a href="http://wiki.hydrogenaudio.org/index.php?title=LAME">Hydrogenaudio Wiki</a> for a pretty succinct breakdown of the other LAME options.
<h3>Output Location</h3>
When I rip a CD, I want all of the output to go to a common location. That way I know where to look when (not if) things go wrong! Change the following options to specify where you'd like your files to go. Here I'm storing them in my home directory.

<pre class="brush:bash;">OUTPUTDIR=/home/ant
WAVOUTPUTDIR=/home/ant/tmp</pre>

You may want to change up your temporary directory to go to a more <em>temporary</em> location like /tmp. For me I just keep it there.
<h3>Output Format</h3>
Here's where things get fun. Everyone has their own file naming scheme, and rightfully so. They are your files after all, right? I like to have the artist name as a top level directory, followed by the album name in another directory under that. For the actual music files (which get placed in the album folder), I like to - again - put the artist in the name in case I want to copy just those files to a USB drive. Here is the setting to make ABCDE create this structure:

<pre class="brush:bash;">OUTPUTFORMAT='${ARTISTFILE}/${ALBUMFILE}/${ARTISTFILE} - ${TRACKNUM} - ${TRACKFILE}'

VAOUTPUTFORMAT='${ARTISTFILE}/${ALBUMFILE}/${ARTISTFILE} - ${TRACKNUM} - ${TRACKFILE}'</pre>

Note that there are two lines there. The first line is for normal albums with only one artist. The second is for music compilations where you have several different artists on the same CD. This is one area that I wish ABCDE did a little better - ideally I'd like to be able to specify my own "album artist" name instead of having folders for each artist on the CD, but I have another solution for that which I'll cover in another post. For now I'm content with that flaw because I rarely purchase compilations.
<h3>Playlists</h3>
As I said before, I want M3U playlists automatically whenever I rip a CD. This gives me a simple, one-click solution to playing an album without having to mess with bad sorting in whatever music player I'm using at the time (they all seem to mess it up). First, use this option to define the format of the M3U playlist. Here I am saying to create the file with the album as it's name, to be stored in the album directory:

<pre class="brush:bash;">PLAYLISTFORMAT='${ARTISTFILE}/${ALBUMFILE}/${ALBUMFILE}.m3u'

PLAYLISTFORMAT='${ARTISTFILE}/${ALBUMFILE}/${ALBUMFILE}.m3u'

</pre>

Next, you need to tell ABCDE to actually create the playlist when you rip a CD (it doesn't do this by default). Add this to do that:

<pre class="brush:bash;">ACTIONS=cddb,read,encode,tag,move,playlist,clean</pre>

Note that there's a lot of other stuff going on before the playlist generation - that's intended.
<h3>Give Me My Space!</h3>
By default ABCDE converts spaces and forward slashes to underscores. I'm fine with the latter, but having underscore in the filename where a space should be is endlessly annoying to me ("Pink Floyd - 03 - Time.mp3" looks much better than "Pink_Floyd_-_03_-_Time.mp3," don't you think?). It's 2010 for goodness sake! Computers have evolved to be able to handle spaces in file names! Replace the "mungefilename" function with this one to keep your spaces:

<pre class="brush:bash;">mungefilename ()
{
 echo &quot;$@&quot; | sed s,:,\ -,g | tr /\* _+ | tr -d \'\&quot;\?\[:cntrl:\]
}</pre>
<h3>Setting Up Aliases</h3>
As I mentioned before, when I was using ABCDE heavily two summers ago I had aliases set up for each physical CD drive on my machine. If you don't have two drives, don't worry about this. If you do, and you think you'd get annoyed or tired of typing out which device from which to rip, define these aliases and put them in your ~/.bashrc. You'll want to substitute your actual device names , of course (mine are /dev/dvd and /dev/cdrw):

<pre class="brush:bash;">alias ripcd='abcde -d /dev/cdrw'

alias ripdvd='abcde -d /dev/dvd'</pre>
<h2>All Together Now</h2>
Okay, so now that you have the configuration, let me show you how it looks in practice. Here I've put a CD (Fall Out Boy's worst CD thus far, IMO)  in my CDRW drive and used my alias to start ABCDE. What appears first is the matching CDDB entries:
<p style="text-align: left;"></p>


<div class="wp-caption aligncenter" style="width: 682px"><a href="/uploads/2010/11/abcde_1.jpg"><img class="size-full wp-image-418 " title="After executing &quot;ripcd&quot;" src="/uploads/2010/11/abcde_1.jpg" alt="After executing &quot;ripcd&quot;" width="682" height="400" /></a><p class="wp-caption-text">After executing &quot;ripcd&quot;</p></div>
<p style="text-align: left;">After reviewing which entry looks like the best fit I choose #2 (mainly because it is listed as "rock" and not "misc"):</p>
<p style="text-align: left;"></p>


<div class="wp-caption aligncenter" style="width: 682px"><a href="/uploads/2010/11/abcde_2.jpg"><img class="size-full wp-image-419 " title="Choosing the second entry" src="/uploads/2010/11/abcde_2.jpg" alt="Choosing the second entry" width="682" height="400" /></a><p class="wp-caption-text">Choosing the second entry</p></div>
<p style="text-align: left;">I then answer a couple of questions - do I want to edit the CDDB entry (not this time, but you can with your editor of choice) and is this CD multi-artist (no again) - and it starts ripping:</p>
<p style="text-align: center;"></p>


<div class="wp-caption aligncenter" style="width: 682px"><a href="/uploads/2010/11/abcde_3.jpg"><img class="size-full wp-image-420 " title="Ripping started" src="/uploads/2010/11/abcde_3.jpg" alt="Ripping started" width="682" height="400" /></a><p class="wp-caption-text">Ripping started</p></div>
<p style="text-align: left;">When everything is done ABCDE will happily drop you back to your command prompt. Here I'm inspecting the files created after the rip finished. All looks well!</p>
<p style="text-align: center;"></p>


<div class="wp-caption aligncenter" style="width: 717px"><a href="/uploads/2010/11/abcde_4.jpg"><img class="size-full wp-image-421 " title="Rip/encode done!" src="/uploads/2010/11/abcde_4.jpg" alt="Rip/encode done!" width="717" height="441" /></a><p class="wp-caption-text">Rip/encode done!</p></div>
<h2>Conclusion</h2>
