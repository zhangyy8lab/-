#!/bin/bash
set -ex

if [ "${1:0:1}" = '-' ] || [ "${1:0:1}" = '' ]; then
  set -- /usr/local/server/mainServer "$@"
fi

exec $@
