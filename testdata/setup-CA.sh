#!/bin/sh
#
# Set up test CAs.
#
# You are not expected to run this script: all generated files are already
# checked in to the repository. This script is just meant to show how this
# directory was initialized.
#
# These CAs are intended for use in unit tests only.

set -o errexit
set -o noclobber
set -o nounset

readonly testdata="$(dirname $0)"

if [ -e "${testdata}/CA" ]; then
  echo "error: Test CAs seem to be already set up." >&2
  exit 1
fi

set -o xtrace

mkdir -p CA/root/newcerts
touch CA/root/index.txt
echo '01' > CA/root/serial

openssl ecparam -name secp384r1 -genkey -out CA/root/key.pem
openssl req -new -subj '/O=Web Packager Test/CN=Root CA' \
  -key CA/root/key.pem -out CA/root/csr.pem
openssl ca -config openssl.cnf -name CA_root \
  -in CA/root/csr.pem -out CA/root/cert.pem \
  -startdate 20200301000000Z -enddate 20300227000000Z -selfsign -notext -batch

mkdir -p CA/inter/newcerts
touch CA/inter/index.txt
echo '01' > CA/inter/serial

openssl ecparam -name secp384r1 -genkey -out CA/inter/key.pem
openssl req -new -subj '/O=Web Packager Test/CN=Intermediate CA' \
  -key CA/inter/key.pem -out CA/inter/csr.pem
openssl ca -config openssl.cnf -name CA_root \
  -in CA/inter/csr.pem -out CA/inter/cert.pem \
  -startdate 20200301000000Z -enddate 20250228000000Z -notext -batch
