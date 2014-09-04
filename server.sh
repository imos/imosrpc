#!/bin/bash

set -e -u
export TMPDIR="$(mktemp -d "${TMPDIR:-/tmp}/XXXXXX")"
trap "rm -R '${TMPDIR}'" EXIT

if [ "$#" -lt 1 -o "${1:0:1}" = '-' ]; then
  echo 'Usage: server.sh repository [arguments...]' >&2
  exit 1
fi

target="${TMPDIR}/main.go"
cat << EOM > "${target}"
package main
import "github.com/imos/imosrpc"
import _ "${1}"
func main() { imosrpc.Serve(); }
EOM
shift
go run "${target}" "$@"
