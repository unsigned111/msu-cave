#!/bin/bash

while :
do
	cat /sys/class/thermal/thermal_zone0/temp
	sleep 3
done
