---
title: "Retropie Setup"
date: "2014-11-24 21:50:00-05:00"
categories:
   - "geekery"
---
I've been fiddling a lot with [RetroPie](http://blog.petrockblock.com/retropie/) lately. It satisfies a lot of the nostalgic gaming urges I get from time to time. There's just nothing like a bout of Super Metroid to cleanse the soul! With how much fun I was having with the system, I wanted to share the love. 

I will admit that it wasn't exactly a smooth process from end to end. I certainly had a lot of trials and tribulations along the way to the smooth system that I have today. When a friend of mine purchased a [CanaKit B+ kit](http://www.amazon.com/CanaKit-Raspberry-Complete-Original-Preloaded/dp/B008XVAVAW) from Amazon at my recommendation, I took it as an opportunity to document the process from end to end. I gave him my completely-functional SD card and I took his. What follows are the steps I took to get a working system. Here we go!

<!--more-->

First grab the latest [RetroPie SD card image](http://blog.petrockblock.com/download/retropie-project-image/) (v2.3 at the time of this writing). After reading a few forum posts, this seemed like the way to go. It's an extremely quick way to get up and running. Start off by unzipping the raw image into a directory of your choice:

<pre><code>
unzip RetroPieImage_ver2.3.img.zip
</code>
</pre>
	
Next, copy the image onto your SD card of choice with the dd command. My SD card was /dev/sde at the time, so plug in your device letter accordignly:

<pre><code>
sudo dd if=RetroPieImage_ver2.3.img bs=2M of=/dev/sde
</code>
</pre>

Once the copy finishes, take out the SD card and boot up the Pi with your
keyboard, mouse, and game controller connected. I have a [Logitech Dual Action](http://www.amazon.com/Logitech-Dual-Action-Game-Pad/dp/B0000ALFCI/ref=sr_1_1?ie=UTF8&qid=1416884420&sr=8-1&keywords=logitech+dual+action+gamepad&pebp=1416884433304) controller that works fine. Once you have everything all connected, plug in your power supply and wait for your system to boot. You'll go right into emulationstation by default, but we still have some other stuff to do first. Press "F4" to exit out of emulationstation without setting it up.

Once you are back on the command prompt, further configure your Pi with the following:

<pre><code>
sudo raspi-config
</code>
</pre>

This will bring up a textual configuration interface. Follow the prompts to expand your root filesystem (option #1) and change your default password (option #2). If you live in the US, you may also want to set your keyboard locale as well - by default, it ships with a en-GB locale. You can change this with option #4. 

If you feel comfortable overclocking your Pi, now is the time to do it. The overclocking options are located under option #7. I chose to overclock mine to 950MHz, which is the "high" option. I've found this to give reasonable performance without completely maxing out the system.

After you are done configuring to your heart's content with raspi-config, exit and reboot using the options provided. You'll hit the emulationstation menu again. This time go through the steps to configure your game controller. It is nothing more than hitting the buttons it tells you to on the screen. Once done with that, make sure the interface is responsive to your newly-configured buttons. If something has gone awry, you can reconfigure by hitting the "start" button and choosing the "configure input" option. 

Now that the frontend is all set, next up you'll configure your controller for the emulators the frontend will use. The default RetroPie image ships with a ton of them, but thankfully they can all leverage one controller setup. You'll have to run through each of the buttons on your controller again. Run the following to do just that:

<pre><code>
cd /opt/retropie/emulators/RetroArch/installdir/bin/
sudo retroarch-joyconfig -j 0 -o /opt/retropie/configs/all/retroarch.cfg
</code>
</pre>

Last but not least, alter the config file you just created to allow you to exit out of your gaming emulators without needing a keyboard. Do that with the following:

<pre><code>
sudo nano /opt/retropie/configs/all/retroarch.cfg
</code>
</pre>

Add the following lines to establish a hotkey. This hotkey will let you exit emulators by pressing start+select. Note that the actual values may vary here. Be sure to use the code values for the keys you want! They can be found in the same file.

<pre><code>
input_enable_hotkey_btn = "8"
input_exit_emulator_btn = "9"
</code>
</pre>

With these options set up, you can leave your keyboard and mouse at home.
You're now ready from retro-gaming action! Enjoy. 
