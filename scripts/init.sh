# Run this script under project root
export THISDIR=`pwd`
echo "Add ${THISDIR} to GOPATH"
export GOPATH="${GOPATH}:${THISDIR}/"

# Go mod requires below as of Aug 2019
export GO111MODULE=on
