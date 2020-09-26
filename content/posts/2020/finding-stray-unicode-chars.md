---
title: "Finding Stray Unicode Chars"
date: 2020-09-26T10:58:05-04:00
---

I maintain a rendering system that uses the [Mako][mako] templating language. The other day I received a pull request that looked innocent at a glance, but blew up with an ASCII decoding error during testing. 

_Side note: Mako has a requirement for all input characters to be ASCII? That doesn't seem right - it's probably my code, or that I'm on an old version of [Python][Python2.7]. I'll have to investigate further._

The PR defined variables similar to the following:

```
# the-template.mako snippet
<%
    some_map["foo"] = "value"
    some_map["bar"] = "other value"
%>
```

So we've received a decoding error. That's good, but where is it _exactly_? Fortunately we can find that out easily enough. Let's try loading up the last-loaded template file `the-template.mako` within the Python REPL to get some finer-grained information:

```
$ python3 
>>> with open("the-template.mako", "rb") as f:
...     data = f.read()
... 
>>> data.decode("ascii")
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
UnicodeDecodeError: 'ascii' codec can't decode byte 0xc2 in position 6: ordinal not in range(128)

```

So this is useful information. We know that the first offending byte value is `0xc2` and it's in byte position 6 in the file. Let's have a look at that position in gVIM. 

After opening up the file in gVIM, some additional information display settings will be useful for honing into the suspect character:

```
# Force display of the status line at the bottom of the vim window.
:set laststatus=2 

# Show the ordinal value of the character under the cursor in the status line.
:set statusline=%b
```

With these in place we could move the cursor throughout the file to see if any ordinal values in the status line are greater than 128. That seems time consuming, so fortunately we can utilize the second piece of information in the error message to quickly move the cursor to the exact byte position in the file (6):

```
# Move the cursor to the 6th byte in the file.
:6go
```

Entering this in gVIM takes the cursor right to that byte position, and from there I can see the statusline highlighting an ordinal value of `160`, bringing the issue to light. Aha! What appeared to be  a whitespace character actually was a unicode [non-breaking space][nbsp] with a hex value of U+00A0. 

Armed with this information, it was then trivial to go through the diffed lines in the pull request to correct all other instances of this sneaky space character to allow rendering to proceed without issue.

[mako]: https://www.makotemplates.org/
[nbsp]: https://en.wikipedia.org/wiki/Non-breaking_space
[Python2.7]: https://docs.makotemplates.org/en/latest/unicode.html