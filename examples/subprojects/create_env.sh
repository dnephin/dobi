#!/bin/sh

set -eu

if [ -f .env -a -n "$(cat .env)" ]; then
    exit
fi

echo "Enter your settings"
echo

echo "Enter the username to use inside Docker.  If not running on the host as root, this should be the same as your host user account."
read -p "Username (default: $HOST_USER_NAME): " username
username=${username:-$HOST_USER_NAME}

read -p "UID (default: $HOST_USER_ID): " uid
uid=${uid:-$HOST_USER_ID}

cat <<EOF > .env
uid=$uid
username=$username
EOF
