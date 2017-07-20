#!/usr/bin/env bash

firebase="https://msu-cave.firebaseio.com"

wifi_address=$(cat /sys/class/net/wlan0/address)

pod=$(curl $firebase/config/pi_list/$wifi_address/pod.json)

pod_data=$(curl $firebase/config/pod_list/$pod.json)

target_headset=$(echo $pod_data | jq .target_headset)

target_address=$(curl $firebase/config/headset_list/$target_headset/address.json)

echo $pod_data | jq ". |= .+ { \"pod\": ${pod}, \"target_address\": ${target_address} }"
