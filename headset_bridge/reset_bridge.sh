#!/bin/bash

./kill_bridge
#kill rfcomm socket here

./pair_headse & disown

./connect_bridge & disown
