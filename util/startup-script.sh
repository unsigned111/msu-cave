#Montana State University NeuroCAVE
#Raspberry Pi Startup Script

#!/bin/bash

if [ $(hostname -I) == "10.0.1.5"]
then
	cd ~
	./ngrok tcp 22
else
	cd ~/msu-cave/
	./start
fi	
	
