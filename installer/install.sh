#!/bin/sh

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

temp="temp-$(uuidgen)"
mkdir $temp

exe_file="$temp/exe"

wget -O $exe_file $download_url
chmod +x $exe_file

dest=""

if [ $EUID -ne 0 ]; then
  dest="$HOME/.local/bin"

  yellow_color="\033[0;33m"
  color_off="\033[0m"
  printf "${yellow_color}make sure the PATH variable contains $dest$color_off\n"
else
  dest="/usr/local/bin"
fi

mv $exe_file "$dest/budgie"
rm -rf $temp

echo "budgie has been successfully installed in $dest"
