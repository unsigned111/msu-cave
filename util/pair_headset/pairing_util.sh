#!/bin/bash
#Montana State University NeuroCAVE
#Headset Pairing Utility


echo "Automatic Pairing Utility"

firebase="https://msu-cave.firebaseio.com"

wifi_address=$(cat /sys/class/net/wlan0/address)
echo "Device Wireless MAC address: $wifi_address"

function conn {
	address=$(echo $1 | tr -d '"')
	echo "Attempting connection to $address"
	sudo rfcomm connect "hs1" "$address"
}

override_address=$(curl $firebase/config/pi_list/{$wifi_address}/override_address.json)
override_headset=$(curl $firebase/config/pi_list/{$wifi_address}/override_hs.json)
if [$override_address != "none"]
then
	echo "Address override found. Connecting to $override_address"
	while :
	do
		echo "attempting connection to headset"
		conn $override_Address
		sleep 3
	done
elif [$override_headset != "none"]
then
	echo "Headset Override found. Connecting to headset $override_headset"
	target_address=$(curl $firebase/config/headset_list/$override_headset/address.json)
	echo "Target MAC address: $target_address"
	while :
	do
		echo "attempting connection to headset"
		conn $target_address
		sleep 3
	done
else
	echo "No overrides found. Determining headset..."
	pod=$(curl $firebase/config/pi_list/$wifi_address/pod.json)
	echo "pod $pod"
	target_headset=$(curl $firebase/config/pod_list/$pod/target_headset.json)
	echo "Headset $target_headset"
	target_address=$(curl $firebase/config/headset_list/$target_headset/address.json)
	echo "Target MAC address: $target_address"
	while :
	do
		conn $target_address
		sleep 3
	done
fi
