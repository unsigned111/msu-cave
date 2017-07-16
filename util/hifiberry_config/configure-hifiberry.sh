#HiFiBerry DAC+ Setup Script

#!/bin/bash

user=$(whoami)

sudo apt-get install csound

sudo cp config.txt /bin/
sudo cp asound.conf /etc/

sudo usermod -a -G audio $user

echo ""
echo "HiFiBerry Successfully configured. Reboot now to complete installation."
