# Broadcaster

The boadcaster module is the communication portal for each node.  It allows the
node to publish and receive data to the shared data source

It is built using [firebase](https://console.firebase.google.com/) as the shared
data source.  To set up a new firebase see [setting up a firebase
project](https://firebase.google.com/docs/server/setup)

TODO:DLM: Add the command for running it should pass in (at least)
* firebase key
* url of firebase root
* id of instalation
* id of node

The project is written in node.  Dependency management is with npm.  When adding
a new dependency, to create a reproducible build, please install with

```
    npm install package-name --save --save-exact
```

# NOTES
Currently anyone can write to the db, and we are connecting as an admin,
we may want to consider in a later iteration adding some auth
but that shouldn't stop us now

