# Running

The project is broken into a collection of small application that communicate
with each other via TCP. Each has a small job.  The motivation for this is that
some of the application are more suited to specific languages.

The file `Procfile` gives a description for how to run each program.
For example, the line
```
    headset_bridge: cd headset_bridge; ./headset_bridge
```
Means that you run headset bridge by running,
```
    > cd headset_bridge; ./headset_bridge
```

If you have [Forman](https://github.com/ddollar/foreman) installed, you can
run all the programs with the command `foreman start`.  It isn't necessary to
install though.  Running the commands individually is easy enough.

# Notes

I wrote the headset bridge program in Go because I could not find a node library
other than Gobot that was independent of the neruosky drivers.  However, I just
found the [cyclonejs](https://cylonjs.com/documentation/platforms/neurosky/)
project.  This may allow us to remove the headset bridge all together.
