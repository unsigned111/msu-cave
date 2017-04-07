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

all: lighting

run: lighting
	./lighting

profile: lighting
	go build -o lighting -race
	./lighting -test.cpuprofile=lighting.prof
	go tool pprof lighting lighting.prof

test: lighting
	go get github.com/stretchr/testify/assert
	go test ./

lighting: pkg
	go build -o lighting

arm: pkg
	env GOOS=linux GOARCH=arm go build -o lighting.arm

check:
	go fmt ./*.go
	go tool vet ./

bundle: arm
	tar -czf lighting.tar.gz lighting.arm settings.json install.sh lighting.service

# depend is just an alias for pkg.
depend: pkg

# get, build, and install the dependencies.  Note that this requires
# an internet connection.  You will not be able to do this successfully
# while off line.
pkg:
	go get github.com/hypebeast/go-osc/osc github.com/dordille/artnet

# standard go clean
clean:
	go clean

# clean up whatever go compiler creates and all of the dependencies.
# Note that if you do this while off line you will no longer be able to
# re-recompile.
veryclean: clean
	-rm -rf ${PWD}/src/
	-rm -rf ${PWD}/pkg/
	-rm -rf ./github.com/
	-rm -rf ./golang.org/
	-rm ./lighting
	-rm ./lighting.arm
	-rm ./lighting.tar.gz
