# Broadcaster

The boadcaster module is the communication portal for each node.  It allows the
node to publish and receive data to the shared data source

It is built using [firebase](https://console.firebase.google.com/) as the shared
data source.  To set up a new firebase see [setting up a firebase
project](https://firebase.google.com/docs/server/setup)

Broadcaster is written in node.js and dependencies are managed with npm (node
package manager).  See [node install](https://nodejs.org/en/) for direction on
installing node (npm now comes bundled as part of node.js).

When adding a new dependency, to create a reproducible build, please install
with

```
    npm install package-name --save --save-exact
```

When developing, it is sometimes useful to be able to easily send requests to
a running server.  The script `generate_request.sh` sends requests of
various types using [httpie](https://github.com/jkbrzt/httpie)

# Running

This description assumes the setup which I am using on my dev machine.  Firebase
requires a credential file, which is NOT in version control.  Please contact me
if you need credentials.

Provided that you are using the setup in my dev machine, you can run the
broadcaster with the command:
```
    node main.js -i [instllation id]  -e [eeg headset id]
```

There are additional arguments, but these are only needed if you are trying to
run on a different Firebase or dev environment.

# NOTES
* The data is only accessable to admins, so, currently we are connect to the
db as an admin.  We may want to consider in a later iteration adding better
access control but that shouldn't stop us now.

