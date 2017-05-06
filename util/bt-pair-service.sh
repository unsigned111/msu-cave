#!/bin/bash

#Bluetooth Pairing Service

echo "Starting Bluetooth Pairing Service"

#firebase= "https://msu-cave.firebaseio.com/headset_list"	#Firebase URL

while :
do
	target_address= "9C:B7:0D:90:ED:E5"			#Insert Target address here
	sleep 5
	echo "\nConnecting to {$target_address}"
	sudo rfcomm connect "hs1" {$target_address} 2
	echo "Reconnecting..."
done
