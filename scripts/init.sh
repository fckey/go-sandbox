# Run this script under twitter-analysis
export THISDIR=`pwd`
echo "Add ${THISDIR} to GOPATH"
export GOPATH="${GOPATH}:${THISDIR}/go"

# Run this only one time
# go get github.com/dghubble/go-twitter
# go get github.com/dghubble/oauth1
go get -u github.com/golang/dep/cmd/dep
