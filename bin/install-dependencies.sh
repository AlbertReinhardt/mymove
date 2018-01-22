#! /bin/bash

set -eux -o pipefail

# Install utilities

# Throwing away stdout logs because they were too plentiful and
# Circle couldn't display them in their web interface.
# Errors should still post to the console.
sudo apt-get update
sudo apt-get -qq install -y apt-transport-https > /dev/null
# This is needed to use `psql` to test DB connectivity, until the app itself
# starts making database connections.
sudo apt-get -qq install -y postgresql-client-9.4 > /dev/null

# Install Go

# Throwing away stdout logs because they were too plentiful and
# Circle couldn't display them in their web interface.
# Errors should still post to the console.
curl -L -O https://redirector.gvt1.com/edgedl/go/go1.9.2.linux-amd64.tar.gz
sudo tar xvf go1.9.2.linux-amd64.tar.gz -C /usr/local > /dev/null
sudo ln -s /usr/local/go/bin/go /usr/local/bin/go

# Install dep
go get -u github.com/golang/dep/cmd/dep

# Install Node and Yarn
# Throwing away stdout logs because they were too plentiful and
# Circle couldn't display them in their web interface.
# Errors should still post to the console.
sudo curl -sL https://deb.nodesource.com/setup_6.x -o nodesource_setup.sh
sudo bash nodesource_setup.sh
sudo apt-get -qq install -y nodejs > /dev/null
sudo curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
sudo apt-get update && sudo apt-get install yarn
