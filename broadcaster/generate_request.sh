#!/bin/bash
address="localhost:3000"

# send a valid request
http "$address" timestamp:="$(date +%s)" delta:=1 theta:=2 loAlpha:=3 hiAlpha:=4 \
    loBeta:=5 hiBeta:=6 loGamma:=7 midGamma:="$(date +%s)"

# send an request missing some data
http "$address" delta:=1 theta:=2 loAlpha:=3 hiAlpha:=4 \
    hiBeta:=6 loGamma:=7 midGamma:=8

# and a request with empty payload
http "$address"
