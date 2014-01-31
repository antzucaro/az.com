---
title: "The US Capitol in HDR"
date: "2011-08-07"
categories: 
  - "photography"
  - "wandering around"
---
This is the US Capitol as seen in HDR. I've posted this pic before, but that was just the raw image. Tonight I thought I'd take a break from the San Francisco pics and post something out of the ordinary. I've been following Trey Ratcliff's stuff, and that gave me the idea to try the HDR concept. More than just that, I threw in some panoramic stuff in for fun.

<div class='wp-caption aligncenter' style='width: 660px; margin-left: auto; margin-right: auto;'>
<img width='650px' height='323px' alt="The US Capitol, seen in HDR" title='The US Capitol, seen in HDR' src='/uploads/2011/08/USCHDR/capitol_hdr_m.jpg'>
<p class='wp-caption-text'>The US Capitol, seen in HDR <a href='/uploads/2011/08/USCHDR/capitol_hdr_l.jpg'><img alt='See this image in full size' src='/static/fs_img.jpg' /></a></p>
</div>

So here's the genesis of this, if you're interested. I started with a set of 8 Nikon NEFs (that's their RAW format) which I corrected with ufraw. I then stitched those images together with hugin, which gave me a full sized, normally-exposed panorama. Having noted the original exposure value in ufraw and also having saved the hugin project, I then went back and adjusted the exposure on the original 8 raw images, saving them under their original names (overwriting). I did this twice; the first time I adjusted the exposures down one exposure value (to -1 EV), and the second time up one (to +1 EV). Each time I did this I ran the panorama through hugin again using the saved project file, which gave me the same exact orientation, but with differently exposed photos each time. 

In the end I had three panoramas of differing exposure values. This was a suitable set to feed into qtpfsgui, which I did using the "Fattal" flavor of algorithm to produce the result above. A good bit of work, but voila!


