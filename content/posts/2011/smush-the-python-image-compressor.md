---
title: "Smush - The Python Image Compressor"
date: "2011-04-06"
categories: 
  - "geekery"
---
Like the good web citizen I am, I was looking the other day on <a href="http://www.webpagetest.org" title="WebPagetest home page">WebPagetest</a> to see how this very site fared in terms of speed. Being that it is a static site I was expecting to take home a good grade report to my Mom and Dad, but lo and behold I was greeted with the following after waiting for the test to finish:

<div class='wp-caption aligncenter' style='width: 433px; margin-left: auto; margin-right: auto;'>
<img width='423px' height='122px' alt="Ouch" title='Ouch' src='/uploads/2011/04/webpagetest.gif'>
<p class='wp-caption-text'>Ouch!</p>
</div>


Fs almost across the board. Ouch. Unless those Fs meant "fantastic," I wasn't doing a very good job. Looking across those categories I identified one area that I thought was an easy target to improve upon - image compression. Seeing as I'd read about compression after reading <a href="http://stevelosh.com/" title="Steve Losh's home page">Steve Losh's</a> <a href="http://stevelosh.com/blog/2010/01/moving-from-django-to-hyde/" title="Steve's static website post">post</a> on how he did his static site (probably from a <a href="http://news.ycombinator.com" title="Hacker News home page">Hacker News</a> link somewhere), I thought I'd do the same. While the aforementioned post listed a Ruby utility for handling the compression of directories of images (recursively, of course), I'm a Python guy. To Google I went to find a Python compression utility for images. 

Eventually I found something that met my requirements: <a href="https://github.com/thebeansgroup/smush.py" title="Smush's github home page">smush.py</a>. It compresses jpgs, pngs, and gifs recursively in their directories. Nice!  I did have one issue before I could start playing with it, however: Python 2.6 is standard on Ubuntu 10.10 (which I'm running), so I had to go grab Python 2.7. Nothing a little command line action couldn't fix quickly:

<pre class="brush:bash;">
wget http://python.org/ftp/python/2.7.1/Python-2.7.1.tgz
tar -xvzf Python-2.7.1.tgz
cd Python-2.7.1
./configure
make
sudo make altinstall
</pre>

Note that I'm using the "altinstall" option here, as I didn't want to mess with my default (2.6) Python installation on Ubuntu. 

Anyway, once I apt-get install-ed the rest of the requirements listed in the README I fired off the script on my directory of images that I use for this site. Lots of messages scrolled by talking about "converting to progressive" and such, finally terminating in this:

<pre class="brush:bash;">
Smushing Finished

1236 files scanned:
    0 GIFs optimised out of 0 scanned. Saved 0kb
    0 GIFGIFs optimised out of 0 scanned. Saved 0kb
    507 JPEGs optimised out of 1215 scanned. Saved 6479kb
    10 PNGs optimised out of 21 scanned. Saved 77kb
Total time taken: 264.42 seconds
</pre>

Given that my image directory is 118M, a savings of about 6.5M is a 5.5% improvement. Not bad! After rsync-ing those images to my media subdomain I checked webpagetest again:

<div class='wp-caption aligncenter' style='width: 433px; margin-left: auto; margin-right: auto;'>
<img width='423px' height='122px' alt="Ouch. Again." title='Ouch. Again.' src='/uploads/2011/04/webpagetest.gif'>
<p class='wp-caption-text'>Ouch. Again.</p>
</div>

Hrm. Looking further into the details I found out that webpagetest grades the compression of the photos on the page based on how close in size they are to the same image when compressed with a quality level of 50 (for JPEG) in Photoshop. Okay...got it. I tested out some of my photographs (the main content of my site, really) with that quality level and was not impressed. For sure I could do better than the level I'm currently saving at w/ ImageMagick's convert utility, but I'm not willing to go all the way down to 50 - those test images looked blocky and ugly. It was clear that I wasn't ever going to get that 'A' I wanted. Since I'm always going to fail that particular test, I'll just have to be content with saving as much space as I can while retaining the quality level I want for my photographs. Using smush.py I can get that. It's a shame I couldn't have my cake and eat it too, though. C'est la vie!
