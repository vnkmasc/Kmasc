#
# Copyright contributors to the Hyperledger Fabric project
#
# SPDX-License-Identifier: Apache-2.0
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at:
#
# 	  http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

###############################################################################
# Build image
###############################################################################

ARG UBUNTU_VER=24.04
FROM ubuntu:${UBUNTU_VER} AS builder

ARG TARGETARCH
ARG TARGETOS
ARG FABRIC_VER
ARG GO_VER
ARG GO_TAGS

RUN apt update && apt install -y \
    git \
    gcc \
    curl \
    make \
    libssl-dev

RUN curl -sL https://go.dev/dl/go${GO_VER}.${TARGETOS}-${TARGETARCH}.tar.gz | tar zxf - -C /usr/local
ENV PATH="/usr/local/go/bin:$PATH"

ADD . .

RUN make orderer GO_TAGS=${GO_TAGS} FABRIC_VER=${FABRIC_VER}


###############################################################################
# Runtime image
###############################################################################

ARG UBUNTU_VER=24.04
FROM ubuntu:${UBUNTU_VER}

ARG FABRIC_VER

# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
RUN echo 'hosts: files dns' > /etc/nsswitch.conf

ENV FABRIC_CFG_PATH=/etc/hyperledger/fabric
ENV FABRIC_VER=${FABRIC_VER}

COPY    --from=builder  build/bin/orderer           /usr/local/bin
COPY    --from=builder  sampleconfig/msp            ${FABRIC_CFG_PATH}/msp
COPY    --from=builder  sampleconfig/orderer.yaml   ${FABRIC_CFG_PATH}
COPY    --from=builder  sampleconfig/configtx.yaml  ${FABRIC_CFG_PATH}

VOLUME  /etc/hyperledger/fabric
VOLUME  /var/hyperledger

EXPOSE  7050

CMD     [ "orderer", "start" ]