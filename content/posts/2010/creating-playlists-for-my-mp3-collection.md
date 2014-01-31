---
title: "Creating Playlists for my MP3 Collection"
date: "2010-11-07"
categories: 
  - "blog"
  - "geekery"
---
As I've mentioned <a href="/2010/11/abcde-a-better-cd-encoder/">before</a>, about two summers back I went through my whole CD collection and digitized it to MP3 format using a tool called ABCDE. What I didn't mention in that post was that at the time I'd forgotten to create playlists automatically during the process. Bummer! I had all of these albums digitized, but whenever I wanted to play a particular album I'd have to add the folder through my music player or manually select all of the files to add to my playlist. That wasn't too bad, but if you've had any experience with playlists in Windows Media Player, VLC, or even XMMS you'll know that in practice it is just a pain to handle. For whatever reason it may be, the ordering of the album tracks on the playlist <em><strong>always</strong></em> gets messed up when adding manually or via folders. <em><strong>Always!</strong></em> Now, I know I can fix this simply by sorting by the ID3 metadata (click the "track number" column in most players to sort), but the point is that <em>I shouldn't have to do this</em>. I wanted a simple solution.

I get annoyed with these things pretty quickly, so I started thinking about how I could remedy the situation to save myself some time. Obviously I could investigate <em>why </em>each of the players were failing to sort my properly named (and usually tagged) files, but that seemed like too much work. I thought a simple solution would be to create M3U (a <a href="http://www.google.com/url?sa=t&amp;source=web&amp;cd=1&amp;ved=0CBwQFjAA&amp;url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FM3U&amp;rct=j&amp;q=m3u&amp;ei=sMnWTOWRFMP-8AaXnPmZBg&amp;usg=AFQjCNHVxaaNTKNnCiADQmd33edsMs0o4A&amp;sig2=OU39DvcsU8FPv3XmFXfiiQ&amp;cad=rja">standard</a> playlist format) playlist files in each album directory that would contain the exact track ordering of that particular album. That way to play an album all I'd need to do was double click on the playlist file and <em>voila!</em> it would play it in the right order. This would also give me the ability to drag and drop multiple albums without having to worry about sort issues with having multiple track 1s (not too much of a problem with today's players, but still).

I did a little research and found a tool to help me write the playlists. It's a Python library called <a href="http://www.google.com/url?sa=t&amp;source=web&amp;cd=1&amp;ved=0CBkQFjAA&amp;url=http%3A%2F%2Fcode.google.com%2Fp%2Fmutagen%2F&amp;rct=j&amp;q=mutagen&amp;ei=FsvWTPb9IoT48AblrMjYDA&amp;usg=AFQjCNHzts_jR6qh9C_NtJi1COJQEfqHWQ&amp;sig2=eoawYstMNOwkQafn61uG1g&amp;cad=rja">mutagen</a>, and it basically allowed me to programmatically read the metadata of Mp3 files to get the following information (among the many other things it can retrieve):
<ul>
	<li>Track number</li>
	<li>Artist name</li>
	<li>Album name</li>
	<li>Track length (seconds)</li>
</ul>
After installing mutagen, I wrote the following Python script (m3u.py) to traverse through a set of directories and create playlist files at the lowest level. For example, if I ran it on <em>/home/ant/Music/Coldplay/Parachutes</em> I would end up with a playlist file inside the <em>Parachutes</em> directory:

[python]
#!/usr/bin/python

import os
import sys
import glob
from mutagen.mp3 import MP3
from mutagen.easyid3 import EasyID3

def makem3u(dir=&quot;.&quot;):
    try:
        print &quot;Processing directory '%s'.&quot; % dir
        os.chdir(dir)

        # get ID3 meta objects for each mp3,
        # store in a list
        playlist = ''
        mp3s = []
        for file in glob.glob(&quot;*.[mM][pP]3&quot;):
            if playlist == '':
                playlist = EasyID3(file)['album'][0] + '.m3u'
            meta_info = {
                'filename': file,
                'length': int(MP3(file).info.length),
                'tracknumber': EasyID3(file)['tracknumber'][0].split('/')[0],
            }
            mp3s.append(meta_info)

        if len(mp3s) &gt; 0:
            print &quot;Writing playlist %s.&quot; % playlist

            # write the playlist
            of = open(playlist, 'w')
            of.write(&quot;#EXTM3U\n&quot;)

            # sorted by track number
            for mp3 in sorted(mp3s, key=lambda mp3: int(mp3['tracknumber'])):
                of.write(&quot;#EXTINF:%s,%s\n&quot; % (mp3['length'], mp3['filename']))
                of.write(mp3['filename'] + &quot;\n&quot;)

            of.close()

    except:
        print &quot;Error when trying to process directory '%s'. Ignoring...&quot; % dir
        print &quot;Text:&quot;, sys.exc_info()[0]

def main(argv = None):
    if argv is None:
        argv = sys.argv

    # directories containing music files
    dirs = []

    if len(sys.argv) == 2 and sys.argv[1] == '-':
    # we do not have command line arguments,
    # so read from STDIN
       for line in sys.stdin:
           dirs.append(line.strip())
    else:
    # passed in directories on the command line
        for dir in sys.argv[1:]:
            dirs.append(dir)

    # for each directory passed to us, go
    # to it and make the M3U out of the
    # MP3 files there
    for dir in dirs:
        makem3u(dir)

    return 0

if __name__ == &quot;__main__&quot;:
    sys.exit(main())
[/python]

You can pass directories two ways this script - as arguments or via standard input. This allows me to do this:

<pre class="brush:bash;">m3u.py /home/ant/Music/Coldplay/Parachutes</pre>

...or if I wanted to do a bunch of directories in bulk:

<pre class="brush:bash;">find /home/ant/Music -type d -links 2 | m3u.py - </pre>

Note that I made only a small attempt to make this script bulletproof. It gives a decent attempt to find the metadata in the MP3 files it finds in the directory, but if it encounters an error thrown by mutagen it dies out with an error message telling me the directory that failed and the exception it hit. Nine times out of ten (in my testing) the script will die because the tag information in the MP3s are malformed, which means that you wouldn't want to create a playlist off of them using that info anyway.

