#Montana State University NeuroCAVE
#Firebase Backup tool

#!/bin/bash

cd ~/firebase_backup/
directory_name=$(date)
mkdir "$directory_name"
cd "$directory_name"
wget https://msu-cave.firebaseio.com/config.json
wget https://msu-cave.firebaseio.com/installations.json

echo "Successfully backed up firebase"
