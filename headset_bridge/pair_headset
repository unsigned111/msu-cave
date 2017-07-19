#!/bin/bash
#Montana State University NeuroCAVE
#Headset Pairing Utility


echo "Automatic Pairing Utility"

function conn {
	address=$(echo $1 | tr -d '"')
	echo "Attempting connection to $address"
	echo raspberry | sudo -kS rfcomm connect "hs1" "$address"
}

target_address=$(../fetch_firebase_data.sh | jq .target_address)
echo "Target MAC address: $target_address"
while :
do
    conn $target_address
    sleep 3
done
