#Montana State University NeuroCAVE
#Raspberry Pi Startup Script

#!/bin/bash

if [ $(hostname -I) == "10.0.1.5"]
then
	#Start ssh stuff nonsense thing
else
	cd ~/msu-cave/
	./start
fi	
	
