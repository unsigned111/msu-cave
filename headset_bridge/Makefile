# Setup a guard that we are running `make`
# in the current directory. This could be improved, however,
# Setting up the path for  go is much easier if we just run make
# from the expected directory.
REL_PWD := $(dir $(lastword $(MAKEFILE_LIST)))
PWD := $(abspath ${REL_PWD})
ifneq (${REL_PWD}, ./)
  $(error Please run make from the ./${REL_PWD} directory)
endif

# setup the go path
export GOPATH=${PWD}
export GOBIN=${PWD}/bin/

all: headset_bridge

# sometimes it is nice to type make run while doing development
run: headset_bridge
	./headset_bridge

test:
	go test

headset_bridge: pkg headset_bridge.go
	go build -o headset_bridge

headset_bridge.arm: pkg headset_bridge.go
	env GOOS=linux GOARCH=arm go build -o headset_bridge.arm

# depend is just an alias for pkg.
depend: pkg

# get, build, and install the dependencies.  Note that this requires
# an internet connection.  You will not be able to do this successfully
# while off line.
pkg:
	go get -d -u github.com/hybridgroup/gobot/...
	go install github.com/hybridgroup/gobot/platforms/neurosky
	go get github.com/parnurzeal/gorequest
	go get github.com/stretchr/testify/assert

# standard go clean
clean:
	go clean

# clean up whatever go compiler creates and all of the dependencies.
# Note that if you do this while off line you will no longer be able to
# re-recompile.
veryclean: clean
	rm -rf ${PWD}/src/
	rm -rf ${PWD}/pkg/
