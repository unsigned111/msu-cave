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

all: osc_src

# sometimes it is nice to type make run while doing development
run: osc_src
	./osc_src

osc_src: pkg osc_src.go
	go build -o osc_src

osc_src.arm: pkg osc_src.go
	env GOOS=linux GOARCH=arm go build -o osc_src.arm

# depend is just an alias for pkg.
depend: pkg

# get, build, and install the dependencies.  Note that this requires
# an internet connection.  You will not be able to do this successfully
# while off line.
pkg:
	go get github.com/hypebeast/go-osc/osc

# standard go clean
clean:
	go clean

# clean up whatever go compiler creates and all of the dependencies.
# Note that if you do this while off line you will no longer be able to
# re-recompile.
veryclean: clean
	rm -rf ${PWD}/src/
	rm -rf ${PWD}/pkg/
