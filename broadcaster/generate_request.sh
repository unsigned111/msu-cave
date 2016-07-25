#!/bin/bash
address="localhost:3000"

http "$address" delta:=1 theta:=2 loAlpha:=3 hiAlpha:=4 \
    loBeta:=5 hiBeta:=6 loGamma:=7 midGamma:="$(date +%s)"

http "$address" delta:=1 theta:=2 loAlpha:=3 hiAlpha:=4 \
    hiBeta:=6 loGamma:=7 midGamma:=8

http "$address"
