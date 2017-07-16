Automatic Headset Pairing Script

****IMPORTANT: Must run script as superuser****

pairing_util.sh automatically connects a raspberry pi to a headset that is specified in firebase.

It uses the WiFi MAC address of the pi as the identifier to look up its pod number in firebase, 
looks up the headset that is assigned to that pod, then connects to that headset using rfcomm.

The Raspberry Pis are held at https://msu-cave.firebaseio.com/config/pi_list, and are sorted by MAC address.

Inside, there is the option to set override_hs or override_address to a specific target:
	-override_hs will take a headset number and immediately lookup the address of that headset and connect to it.
	-override_address will take a MAC address and immediately connect to the specified device.

If both override values are set to "none," the script will proceed to operate normally.

If connection is unable to be established or is dropped, the script will wait 3 seconds and attempt to reconnect in an infinite loop.

If the a new headset is assigned to a pod, the script must be stopped and restarted (probably want to reboot the entire pod??)



**This script will not work if your raspberry pi and headset are not indexed in firebase**
