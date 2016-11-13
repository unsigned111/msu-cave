# Daemon for Cave lighting controller
## Package Requirements:
    * Ubuntu/Debian packages:
        * python 2.7.*
        * libncurses5
        * python-dev
        * virtualenv (2.7.*)
    * PIP requirements can be found in requirements.txt

## How to setup the running environment:
####1. Makes sure you have libncurses5, python 2.7, etc installed.
On OS 10.10+ these are already installed.
On ubuntu/debian install these packages with apt.
```shell
apt-get install python-dev libncurses5 virtualenv
```

####2. Create and switch to the virtual environment.
```shell
cd project directory
virtualenv env
source env/bin/activate
````
####3. Install the necessary PIP package requirements.
```shell
pip install -r ./requirements.txt
```

## Configuring the daemon:
Settings for the daemon can be altered by editing the `config.yaml` file located at `src/config/config.yaml` in the project directory.

If you change the configuration while the service is running, the service must be restarted before the settings will be used.

## Running the daemon:
Execute the `lightdaemon.py` script from the command line.
```shell
cd src
python lightdaemon.py
```

## Equipment requirements
The lighting equipment should be physically connected following the diagram found in
![diagram](Diagram.png)

Make sure that the DMX output port on the Enttec box is connected to the DMX input port on the light.

The LED light needs to be configured to be in slave mode.
* To cycle the light through modes press both interface buttons at the same time.
    * In slave mode, the LCD display should display a 3 digit number.
    * Once in slave mode, use the up/down buttons to set it to channel 001
    * If the light was previously in slave mode, it should startup in slave mode
    
The Enttec light controller networking is configured to 10.7.153.129.  Make sure the computer's networking settings are such that it can communicate with this IP.

## Troubleshooting

#### Lighting daemon will not start
Check the settings file to make sure it is binding to the correct network address (and that the interface exists/has the right IP address).

#### Light is cycling colors (through the entire spectrum)
The light is in the wrong mode.  Press both interface buttons at the same time to cycle to the next mode.  Repeat until it reads a 3 digit number

#### Light is not on
* Check the light plugged in and turned on.
* Check the light is in slave mode (3 digits on the display)
    * Press both interface buttons at the same time to cycle to the next mode.
    * Repeat until it reads a 3 digit number
* Check the light is on channel 001
    * Use the up/down interface buttons to change the channel)
* Check the Enttec DMX controller is turned on.
* Check the network between the computer and Enttec box is active (check ethernet indicators on the Enttec box)
    * Try pinging the network address 10.7.153.129 to verify the Enttec box can be reached.
* Check the service output to make sure that the service is receiving headset data
    * If data is not being received, check the network settings between the headset source and the computer hosting the lighting service

#### Light is glowing white (and you want it to change colors)
* The service is in "headset off" mode.  Change the "default_on" setting in the service configuration file to True.

#### Light is bouncing between colors (and you want it to glow white)
* The service is in "headset on" mode.  Change the "default_on" setting in the service configuration file to False.

## Presentation notes
The service listens to a single channel of eeg data (configurable).  The intensity of this channel is mapped to a linear gradient between two colors (also configurable).

In the final installation, the lighting will pulse white when the pod is unoccupied.  Once the headset is put on, the lighting will respond to the eeg readings from the headset.  As covariance is achieved, the intensity of the shared color will increase to intensify the shared experience.

## Other notes
* If you have any other problems and can't figure it out, contact Chris Huvaere
    * chuvaere@montana.edu
    * 406-580-1181