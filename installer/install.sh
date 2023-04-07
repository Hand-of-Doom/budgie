#!/bin/sh

if [ $EUID -ne 0 ]; then
  cat << EOF
it requires a root access to move the budgie executable file into the bin directory
please run as root
EOF
  exit
fi

if [ -z $0 ]; then
  cat << EOF
pass the version you need as the first argument
see https://github.com/Hand-of-Doom/budgie/releases
ex: ./install.sh 1.0.1
EOF
  exit
fi

tag="v$0"
download_url="https://github.com/Hand-of-Doom/budgie/releases/download/$tag/budgie"

status=$(curl --write-out %{http_code} --output /dev/null $download_url)
if [ "$status" = 404 ]; then
  cat << EOF
wrong version
see https://github.com/Hand-of-Doom/budgie/releases
EOF
  exit
fi

dest="temp-$(uuidgen)"
mkdir $dest

exe_file="$dest/exe"

wget -O $exe_file $download_url
chmod +x $exe_file

mv $exe_file /usr/local/bin/budgie
rm -rf $dest

echo "budgie has been successfully installed"
