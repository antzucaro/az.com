---
title: "Backing up the Acer Revo 3700"
date: "2011-01-23"
categories: 
  - "blog"
  - "geekery"
---
In my previous post I told you about my Acer Aspire Revo 3700 - awesome things come in small packages! Unfortunately, as great as the overall package is the default operating system included was Windows 7. That's a definite no-go, as I'm not a fan of anything coming out of Microsoft. Obviously that had to go! Being the practical person I am, I didn't want to blow the Windows software <em>completely</em> away. I wanted to keep a backup copy "just in case." What case, you ask? <em>I don't really know, but I like to keep my options open anyway. </em>Here's how I did it.

Normally I'd just fire up my handy copy of <a href="http://www.sysresccd.org/Main_Page">SystemRescueCD</a> to do the backup, but the Revo 3700 doesn't have a CD drive. Bummer. I'd have to boot from a USB image instead. On first glance SystemRescueCD's USB disk creation process seemed like too much work, so I looked for another Linux distribution that might work better. After a small bit of research I ended up settling on <a href="http://grml.org/">grml</a>, which is a Debian-based distribution intended for textmode-savvy users. I followed the instructions I found on their <a href="http://grml.org/download/">download</a> page to install the disk image onto my 2GB USB drive:

<pre class="brush:bash;">dd if=grml_2010.12.iso of=/dev/sdb</pre>

Be careful with that "of=" line if you do this - your USB <em>might not be /dev/sdb</em>! I only have one hard disk on my laptop, so my thumb drives always show up as /dev/sdb when I plug them in. You can verify which letter your thumb drive is by running "dmesg" in your terminal just after you insert the drive. You should see something like this, the bottom of which will tell you the drive identifier of your USB:

<pre class="brush:bash;">
ant@longstreet:~% dmesg
...
[lots of lines snipped]
...
[16937.696105] usb 1-3: new high speed USB device using ehci_hcd and address 2
[16938.005504] Initializing USB Mass Storage driver...
[16938.005646] scsi2 : usb-storage 1-3:1.0
[16938.005760] usbcore: registered new interface driver usb-storage
[16938.005762] USB Mass Storage support registered.
[16939.005085] scsi 2:0:0:0: Direct-Access     PNY      USB 2.0 FD       8.02 PQ: 0 ANSI: 0 CCS
[16939.006616] sd 2:0:0:0: Attached scsi generic sg2 type 0
[16939.007781] sd 2:0:0:0: [sdb] 15695871 512-byte logical blocks: (8.03 GB/7.48 GiB)
[16939.008264] sd 2:0:0:0: [sdb] Write Protect is off
[16939.008273] sd 2:0:0:0: [sdb] Mode Sense: 45 00 00 08
[16939.008280] sd 2:0:0:0: [sdb] Assuming drive cache: write through
[16939.012131] sd 2:0:0:0: [sdb] Assuming drive cache: write through
[16939.012145]  sdb: sdb1
[16939.015272] sd 2:0:0:0: [sdb] Assuming drive cache: write through
[16939.015282] sd 2:0:0:0: [sdb] Attached SCSI removable disk
</pre>

Now back to the backup process. After my dd command finished, I popped my newly-created disk into my Revo. I hit f12 during boot to tell it to boot from USB and I was quickly in business. Speaking of business, my first order was to mount the device on which I would store the backup itself. In my case I was using a USB external hard drive that I mounted with the following command (using the same "dmesg" trick I mentioned above):

<pre class="brush:bash;">mount /dev/sdb1 /media/backup</pre>

My next task was to back up my partition table. In the event that I had to restore what I was about to back up I would need to know how it was physically laid out. I ran the following to save that sweet, sweet partition table:

<pre class="brush:bash;">ï»¿ï»¿sfdisk -d /dev/sda &gt; /media/backup/revo_parttable.bak</pre>

My next step was to mount the media needing to be backed up. There were three partitions on the stock Revo image, so I created a mount directory for each partition. Keeping things simple I just named these directories after the partitions themselves (being the simple man I am):

<pre class="brush:bash;">for i in 1 2 3
do
mkdir -p /media/sda${i}
mount /dev/sda${i} /media/sda${i}
done</pre>

Having those partitions mounted, next up was backing up the actual data. For this task I used my favorite backup utility DAR (For Disk ARchive). DAR is awesome - really, it deserves its own post (which I am now writing down on my TODO list)! For now, though, I'll just list the commands I used:

<pre class="brush:bash;">dar  -R /media/sda1 --gzip 1 -s 2G -c /media/backup/revo_sda1_201001
dar  -R /media/sda2 --gzip 1 -s 2G -c /media/backup/revo_sda2_201001
dar  -R /media/sda3 --gzip 1 -s 2G -c /media/backup/revo_sda3_201001</pre>

Okay, okay. I'll at least give you a breakdown of those command line switches even if I can't get to talking about DAR in depth (who knows if I will, right?).
<ul>
	<li>-R defines the source of your backup. In this case I've mounted each one of the Revo's partitions to a similarly named mount point. Each of those mount points gets backed up.</li>
	<li>--gzip 1 tells DAR to compress the backup using the gzip algorithm. You can specify the level of compression (from 1 to 10, with 10 yielding the smallest image). Here I've used the lowest level of gzip compression.</li>
	<li>-s splits the backup into pieces (slices). I've specified 2GB slices, meaning when the current backup slice reaches 2GB it will start working on a new slice with a different name (it appends an incrementing number to the filename before the suffix).</li>
	<li>-c is the name of the actual backup. It is called a "basename" because the different slices will take on the incrementing number I mentioned above.</li>
</ul>
After those dar commands ran above I was left with the following in my backup directory (I mounted this on my laptop, thus /media/usb and not /media/backup). Sweet!

<pre class="brush:bash;">
ant@longstreet:~% ls -lthrg /media/usb/*dar
-rwxrwxrwx 1 root 2.0G 2011-01-08 08:01 /media/usb/revo_sda1_201001.1.dar
-rwxrwxrwx 1 root 2.0G 2011-01-08 08:08 /media/usb/revo_sda1_201001.2.dar
-rwxrwxrwx 1 root 2.0G 2011-01-08 08:14 /media/usb/revo_sda1_201001.3.dar
-rwxrwxrwx 1 root 2.0G 2011-01-08 08:20 /media/usb/revo_sda1_201001.4.dar
-rwxrwxrwx 1 root 2.0G 2011-01-08 08:26 /media/usb/revo_sda1_201001.5.dar
-rwxrwxrwx 1 root 2.0G 2011-01-08 08:33 /media/usb/revo_sda1_201001.6.dar
-rwxrwxrwx 1 root 420M 2011-01-08 08:34 /media/usb/revo_sda1_201001.7.dar
-rwxrwxrwx 1 root 8.2M 2011-01-08 08:34 /media/usb/revo_sda2_201001.1.dar
-rwxrwxrwx 1 root 2.0G 2011-01-08 08:43 /media/usb/revo_sda3_201001.1.dar
-rwxrwxrwx 1 root 2.0G 2011-01-08 08:50 /media/usb/revo_sda3_201001.2.dar
-rwxrwxrwx 1 root 2.0G 2011-01-08 08:57 /media/usb/revo_sda3_201001.3.dar
-rwxrwxrwx 1 root 2.0G 2011-01-08 09:06 /media/usb/revo_sda3_201001.4.dar
-rwxrwxrwx 1 root 2.0G 2011-01-08 09:15 /media/usb/revo_sda3_201001.5.dar
