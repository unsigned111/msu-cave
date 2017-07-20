#!/usr/bin/env bash

#Update Bash Profile and Sources
#Montana State University NeruoCAVE

function copy_as_hidden {
    file_name=$1
    cp "$file_name" "$HOME/.$file_name"
}

copy_as_hidden bash_profile
copy_as_hidden bashrc
copy_as_hidden vimrc
copy_as_hidden profile

echo "Successfully updated env files"
