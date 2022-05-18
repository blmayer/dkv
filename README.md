# DKV

> DKV is a distributed in memory key-value database.


## Using

To run this database you'll need one instance of the *root* program, located in 
`cmd/root`, and 1 or more of the *node*, they are in `cmd/node`.

To compile each you run `go build ./cmd/<root|node>/`, and the binary will be
created for your platform.

By default the root node listens on port 8080, as it is commonly used, and
the port 1234 is made available for leaf instances to connect.

### Example

Open 4 terminals: on the first run `./root`, the root node must run first.
Then on the other terminals run `./node`, they will connect to the root node.
Now you can make HTTP requests for the root node:

`curl http://localhost:8080/x -d 'test'`

Will add the key *x* with value *test* to some instances. To retrieve it back use:

`curl http://localhost:8080/x`

And so on. Try creating some keys and then killing one node process. You can also
delete keys using:

`curl -X DELETE http://localhost:8080/x`


## Roadmap

- [x] get
- [x] post
- [x] delete
- [x] check on post [commit](https://github.com/blmayer/dkv/blob/d160e34976d570c9373090b23ef3901b8e04bcc7/cmd/root/instances.go#L60)
- [ ] tests
- [ ] make replication better
  - [ ] change from cli
  - [ ] make it a function of the number of instances?
  - [ ] use hash function to select instances instead of random
- [ ] benchmark

## Meta

Initial commit: [#7da73ac908f21fda8d2bee8c601d513899a729b4](https://github.com/blmayer/dkv/commit/7da73ac908f21fda8d2bee8c601d513899a729b4)

License: Creative Commons Attribution 4.0
