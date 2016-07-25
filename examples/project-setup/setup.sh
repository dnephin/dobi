#!/bin/bash

set -eu

if [ -f .env -a -n "$(cat .env)" ]; then
    exit
fi

echo "Enter your settings"
echo

random="$RANDOM$RANDOM$RANDOM"
read -p "Unique id for the project (default: $random): " unique_id
unique_id=${unique_id:-$random}

read -p "Username (default: $HOST_USER): " username
username=${username:-$HOST_USER}

read -p "Listen port (default: 8080): " port
port=${port:-8080}

cat <<EOF > .env
unique_id=$unique_id
username=$username
port=$port
EOF
