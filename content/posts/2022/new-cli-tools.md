+++ 
draft = false
date = 2022-04-13T17:52:51-04:00
title = "New CLI Tools"
description = ""
slug = ""
authors = []
tags = []
categories = []
externalLink = ""
series = []
+++
There was a wonderful post on [new(ish) command line tools][Julia Evans Post] by Julia Evans that was making the rounds
on [Hacker News][Hacker News Thread] this week. It's a joy to see such projects get attention not only because the authors deserve it, but also (selfishly) because it means
I'll likely find something new to play with! I'm always down to optimize my
command line experience. Plus, trying something different is always fun and often *extremely*
rewarding. I guess what I'm trying to say is you'll never know when the "new hotness" will become your "tried and true" unless you leave your comfort zone every once in a while.

Anyway, in similar vein some other
enterprising author clued me into [FASD][FASD] a while back and now I can't live without it. It's one of the
first things I install when I get a new machine. 
For those who aren't already aware or don't want to take another hop to read the project README, FASD is
tiny tool that helps you quickly jump from directory to directory. It does this based on the
paths that you visit the most often, so after you provide it with this "seed" data you can jump
directly between them using fuzzy string matching. This is especially useful if you happen to work
in environments with deeply-nested directories like I do. I typically have projects structured
inside a `Code` directory within my home, then mirror directory structure based on the organization
in source control and the repo's folder itself. This means I'll end up with nested structures like so:

```
/home/ant/Code/some-work-org/org-fancy-tool/src/main/java/package/Foo.java
/home/ant/Code/another-work-org/not-so-fancy-tool/src/main/java/package/Bar.java
```

...and with FASD I can hop between the two of these with just `z Foo` or `z Bar`. Crazy, right?
Think of all the characters saved there. It's literally seconds saved from my life. Seconds, I tell
you!

FASD supports most popular shells (bash, zsh, tcsh) with a teeny bootstrap command that you inject into
your shell's config file. It's really low friction compared with how much joy it brings. Set it,
forget it, and don't regret it. 

While I don't want to ramble on too much about all of the little quality of life CLI tools I use, I do
want to call attention to two more. The first one is [ripgrep][ripgrep]. It's a blazing-fast replacement
for the venerable `grep` UNIX utility with some added niceties like being aware of your `.gitignore` patterns. 
I use it so much that I often find myself starting to type `rg` out in production where (of course)
we don't have such nice things available, so I make a little frowny face at myself before replacing
the command with `grep` and going about my business, the quality of my day ever-so-slightly
decreased.

The other one is [direnv][direnv]. This is another small tool that executes
things you place in a specifically-named shell script (`.envrc`) when you enter or leave directories. This is worth it alone for dealing with the zillion
Python projects that I've managed at work over the years. Of course each of them has their own
set of build requirements, env vars, and the like. So instead of manually starting up the
`virtualenv` each time I need to use it I put the following in a `.envrc` file at the root-level directory of the
project:

```
export SOME_ENV_VAR=some_value
source path/to/my/venv/bin/activate
```

Then the virtual environment and environment var are each "magically" set without me needing to remember
a thing. I know this isn't a terribly fancy use case for this tool, but damn if it doesn't bring
me an inordinate amount of joy. Give it a try and you may see what I mean. 

So, yeah: give some of the tools in Julia's list a try. Maybe even check out the HN thread for more
things that may better suit your fancy. You may find your next "new hotness"!


[direnv]: https://direnv.net/
[FASD]: https://github.com/clvv/fasd
[Julia Evans Post]: https://jvns.ca/blog/2022/04/12/a-list-of-new-ish--command-line-tools/
[Hacker News Thread]: https://news.ycombinator.com/item?id=31009313
[ripgrep]: https://github.com/BurntSushi/ripgrep
