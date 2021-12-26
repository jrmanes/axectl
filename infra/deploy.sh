#!/bin/bash
# Deploy to remote Linux server with arm64 architecture
DEF_HOST="192.168.64.2"
HOST=${1:-$DEF_HOST}
DEF_USER="joseramon"
USER=${2:-$DEF_USER}

make build
scp ./bin/piktoctl ${USER}@${HOST}:/home/${USER}/
#scp ./bin/m1/linux/piktoctl ${USER}@${HOST}:/home/${USER}/
