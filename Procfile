pairing:
    cd util/pair_headset
    ./pairing_util.sh
headset_bridge:
    cd headset_bridge
    ./headset_bridge
broadcaster:
    cd broadcaster
    main.js -i installation_id -e eeg_headset_id \
        -o 127.0.0.1:7771 127.0.0.1:7770
sound:
    cd sound
    csound -+alsa -odac -b2048 -B2048 sound.csd
lighting:
    cd lighting
    ./lighting -f settings.json -d
