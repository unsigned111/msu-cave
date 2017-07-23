#Montana State University
#NeuroCAVE
#Interface Updater

#!/bin/bash

write_file() {
	echo "Generating interfaces file..."
	cd ~/msu-cave/util/generate_interface_file/
	address=$(echo $1 | tr -d '"')
	echo "# interfaces(5) file used by ifup(8) and ifdown(8)" >> interfaces
	echo "# Include files from /etc/network/interfaces.d:" >> interfaces
	echo "" >> interfaces
	echo "# The loopback network interface" >> interfaces
	echo "auto lo" >> interfaces
	echo "iface lo inet loopback" >> interfaces
	echo "" >> interfaces
	echo "# The primary network interface" >> interfaces
	echo "auto $2" >> interfaces
	echo "iface $2 inet static" >> interfaces
	echo "	address $address" >> interfaces
	echo "	netmask 255.255.255.0" >> interfaces
	echo "	network 10.0.0.0" >> interfaces
	echo "	broadcast 10.0.0.255" >> interfaces
	echo "	gateway 10.1.1.1" >> interfaces
}

echo "Fetching data..."
mac_address=$(cat /sys/class/net/wlan0/address)
echo "Wireless MAC address: $mac_address"

pod=$(curl https://msu-cave.firebaseio.com/config/pi_list/$mac_address/pod.json)
echo "Assigned to Pod $pod"
pod=$(echo $pod | tr -d '"')
static_ip=$(curl https://msu-cave.firebaseio.com/config/pod_list/$pod/ip_address.json)
echo "Assigned IP address $static_ip"
current_ip=$(hostname -I)
echo "Current address: $Current_ip"

cd /sys/class/net/enx*/
interface_mac=$(cat address)
interface=$(echo "enx$interface_mac" | tr -d ':')

cd ~/msu-cave/util/generate_interface_file/

if [ -f "interfaces" ]
then
	rm interfaces
fi

write_file $static_ip $interface
