---
title: "Unicode Characters"
date: "2014-05-16 09:55:00-05:00"
categories:
   - "geekery"
tags:
   - "unicode"
   - "python"
   - "xonstat"
---
The other day I was adding support for some unicode characters within [XonStat](http://stats.xonotic.org) and during the process stumbled upon a few helpful things. It turns out that many Linux applications provide for the easy entry of code point values with a simple key press combination: while holding the control and shift keys, enter 'U' and then the unicode code point value to insert that character into whatever text buffer you're currently working with. For example, entering the character "é" (at [code point](http://en.wikibooks.org/wiki/Unicode/Character_reference/0000-0FFF) 00E9) you would hold control and shift and then enter U, 0, 0, E, 9. Nice!

<!--more-->

What's even better is that Firefox is one of the apps that supports this entry method. This is especially useful for playing around with your @font-face declarations, since entry behaves the same way within the developer tools as well. That meant that instead of searching for a Xonotic player who had a special character in their nick (which is what I was testing), I could simply inspect their nickname and insert an arbitrary code point to see if it displayed the font's glyph properly. Similarly, I could also specify any code point in a POST request for incoming data by doing the same thing if needed (at the time I didn't).

The unfortunate thing is that many text editors have had the same functionality for a long time, and each has developed their own entry pattern. Entering these characters in Vim, for example, requires hitting Control-V while in insert mode, then entering U and the code point value. This isn't bad at all, it's just slightly different from how other applications handle it. I guess you can't have it all!

During my testing process I also found out that Python 2 has several flavors of unicode literals developers can use. I was naively assuming that all code points could be entered with one syntax, so I was assembling four and five hex-character code points with the same format: \uXXXXX, with the Xs being the code point value in hexadecimal. My @font-face had its glyphs mostly in the four character range, but it did have a few in the five character range. It was those that weren't working properly, which caused me endless frustration. I was entering "\u1F680" as the code point literal in the code, but I kept getting an omega symbol (U+1F68) and a 0 on my display! I kept scratching my head over this until I looked closely into the [Python documentation](https://docs.python.org/2/howto/unicode.html#unicode-literals-in-python-source-code). I soon found out that there is a separate syntax for code points above the four character range: \UXXXXXXXX. The fix in my code was to change all of the five-character code points to have an upper-case U and to pad them such that they contained 8 hexadecimal characters in total. You can see a few examples of this change in [this line of code](https://github.com/antzucaro/XonStat/blob/master/xonstat/util.py#L58). After making this change all of the glyphs in the font worked perfectly. Phew!
