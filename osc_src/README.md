# Headset Data Playback Simulator

This tool takes file recordings of headset data back as OSC messages, simulating the headset bridge output.  It's intended to be used as a testing and debugging tool to make sure other software components (lighting/sound/etc) work, without requiring the complex setup of the headset bridge.

## Building
To build the application for the desktop, `cd` into the directory and run `make`.  To build the application for the Pi/ARM, `cd` into the directory and run `make osc_src.arm`.

## Running
To run the application in normal mode:
`./osc_src'

### Available Command Line Flags:
  * `-f`: the headset data file to use as the source (default = './default.csv')
  * `-p`: the port number to send OSC messages to (default = 7770)
  * `-d`: the delay amount (in seconds) (default = 1 second)
  * `-t`: toggle the headset on/off endpoint ('/onoff') (10x the delay amount) (default = off)
  * `-s`: send random valid values to the similarity endpoint ('/similarity') (5x the delay amount) (default = off)
