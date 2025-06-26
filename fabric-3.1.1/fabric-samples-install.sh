#!/bin/bash

# Tải fabric-samples repository
if [ ! -d "fabric-samples" ]; then
  git clone https://github.com/hyperledger/fabric-samples.git
else
  echo "fabric-samples directory already exists. Skipping clone."
fi

# Tải install-fabric.sh
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh

# Cấp quyền thực thi cho install-fabric.sh
chmod +x install-fabric.sh

# Chạy install-fabric.sh với các tham số docker samples binary
./install-fabric.sh docker samples binary 