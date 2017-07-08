#configure-hifiberry.sh

This tool will automatically configure a Raspberry Pi to use the HiFiBerry DAC+ sound card and will need to be run as superuser.

The script will:
	-Install csound
	-Update /boot/config.txt
	-Create /etc/asound.conf
	-Add user to "audio" group

Once script has completed, a reboot will be necessary. After rebooting, to test if the installation was successful, run "aplay -l"

If successful, your output should resemble the following:

**** List of PLAYBACK Hardware Devices ****
card 0: sndrpihifiberry [snd_rpi_hifiberry_dac], device 0: HifiBerry DAC HiFi pcm5102a-hifi-0 []
Subdevices: 1/1
Subdevice #0: subdevice #0
