#!/bin/bash
cd
pwd
sudo apt-get update && sudo apt-get upgrade
sudo apt install git -y
sudo apt-get install gcc -y
sudo apt-get install make -y
sudo apt-get install autoconf -y
sudo apt install libsecret-1-0 libsecret-1-dev libglib2.0-dev -y
sudo make --directory=/usr/share/doc/git/contrib/credential/libsecret
git config --global credential.helper /usr/share/doc/git/contrib/credential/libsecret/git-credential-libsecret
