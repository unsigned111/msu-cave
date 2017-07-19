#Montana State University NeuroCAVE
#Raspberry Pi Startup Script

#!/bin/bash

#Get Script Updates
cd ~/msu-cave/
git pull

current_hash=$(cat /sys/network/interfaces | md5sum)
cd ~/msu-cave/util/generate_interface_file/
./generate-interface.sh
new_hash=$(cat interfaces | md5sum)

if["$current_hash" != "$new_hash"]
then
	echo "hashes dont match"
	#echo raspberry | sudo -kS mv interfaces /sys/network/interfaces
	#echo raspberry | sudo -kS reboot
fi

#Start Forego
cd /home/pi/msu-cave
./forego start
