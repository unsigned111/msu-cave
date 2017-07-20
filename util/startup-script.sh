#Montana State University NeuroCAVE
#Raspberry Pi Startup Script

#!/bin/bash

#Get Script Updates
cd ~/msu-cave/
git pull

mac=$(cat /sys/class/net/wlan0/address)
pod=$(curl https://msu-cave.firebaseio.com/config/pi_list/$mac/pod.json)
ip=$(curl https://msu-cave.firebaseio.com/config/pod_list/$pod/ip_address.json)

if [ $(hostname -I) == $ip ]
then
	if [ $pod == "server" ]
	then
		echo "SSH Server Configured"
	fi
	echo "Network configuration correct."
	echo "Starting Pod..."
	cd ~/msu-cave/
	./start
else
	./util/generate_interface_file/generate-interface.sh
	#echo "Replacing interfaces file"
	#echo raspberry | sudo -kS mv interfaces /sys/network/interfaces
	#echo "Rebooting..."
	#echo raspberry | sudo -kS reboot
fi
