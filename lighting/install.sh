#!/bin/bash
export CONFIGPATH=/etc/lighting
export INSTALLPATH=/usr/bin
export SERVICEPATH=/etc/systemd/system

mkdir $CONFIGPATH
cp ./settings.json $CONFIGPATH/settings.json
cp ./lighting.arm $INSTALLPATH/lighting
cp ./lighting.service $SERVICEPATH/lighting.service

systemctl daemon-reload
systemctl enable lighting.service
systemctl start lighting.service
