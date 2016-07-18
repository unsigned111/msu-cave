# Setup a guard that we are running
#   > make
# and we are in the current direcotry.  Setting up the path for
# go is much easier if we just run make from the expected directory.
REL_PWD := $(dir $(lastword $(MAKEFILE_LIST)))
PWD := $(abspath ${REL_PWD})
ifneq (${REL_PWD}, ./)
  $(error Please run make from the ./${REL_PWD} directory)
endif

# setup the go path
export GOPATH=${PWD}
export GOBIN=${PWD}/bin/

run: build
	./headset_bridge

build: pkg
	go build

# depend is just an alias for pkg.
depend: pkg

# get, build, and install the dependencies.
pkg:
	go get -d -u github.com/hybridgroup/gobot/...
	go install github.com/hybridgroup/gobot/platforms/neurosky
	go get github.com/parnurzeal/gorequest

clean:
	go clean

veryclean: clean
	rm -rf ${PWD}/src/
	rm -rf ${PWD}/pkg/
