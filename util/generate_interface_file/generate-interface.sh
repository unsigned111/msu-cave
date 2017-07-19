#Montana State University
#NeuroCAVE
#Interface Updater

#!/bin/bash

function write_file {
	echo "Generating interfaces file..."
	address=$1
	echo "# This file describes the network interfaces available on your system" >> interfaces
	echo "# and how to activate them. For more information see interfaces(5)" >> interfaces
	echo "" >> interfaces
	echo "# The loopback network interface" >> interfaces
	echo "auto lo" >> interfaces
	echo "iface lo inet loopback" >> interfaces
	echo "" >> interfaces
	echo "# The primary network interface" >> interfaces
	echo "auto eth0" >> interfaces
	echo "iface eth0 inet static" >> interfaces
	echo "	address $address" >> interfaces
	echo "	netmask 255.255.255.0" >> interfaces
	echo "	network 10.0.0.0" >> interfaces
	echo "	broadcase 10.0.0.255" >> interfaces
	echo "	gateway 10.0.0.1" >> interfaces
}

echo "Fetching data..."
mac_address=$(cat /sys/class/net/wlan0/address)
echo "Wireless MAC address: $mac_address"

static_ip=$(curl https://msu-cave.firebaseio.com/config/pi_list/$mac_Address/ip_internal.json)
dynamic_ip=$(ifconfig | grep \"inet\" | grep -v 127.0.0.1)
echo "Current DHCP address: $dynamic_ip"
echo "Static IP to be assigned: $static_ip"

write_file $static_ip

#echo "Replacing interfaces file"
#mv interfaces /etc/network/
