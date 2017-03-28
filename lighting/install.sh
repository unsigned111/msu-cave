#!/bin/bash
export CONFIGPATH=/etc/lighting
export INSTALLPATH=/usr/bin
export SERVICEPATH=/etc/systemd/system

echo "Installing configuration file"
mkdir $CONFIGPATH
cp ./settings.json $CONFIGPATH/settings.json
echo "Installing binary"
cp ./lighting.arm $INSTALLPATH/lighting
echo "Installing systemd binary"
cp ./lighting.service $SERVICEPATH/lighting.service
echo "Enabling lighting systemd service"
systemctl daemon-reload
systemctl enable lighting.service
systemctl start lighting.service
echo "Completed!"
