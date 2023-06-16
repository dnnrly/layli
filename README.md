# layli

This tools produces diagrams and has 2 simple aims:

1. Define components and connections as code
2. Look pretty

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/dnnrly/layli)](https://github.com/dnnrly/layli/releases/latest)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/dnnrly/layli/release.yml?branch=master)](https://github.com/dnnrly/layli/actions/workflows/release.yml?branch=master)
[![report card](https://goreportcard.com/badge/github.com/dnnrly/layli)](https://goreportcard.com/report/github.com/dnnrly/layli)
[![godoc](https://godoc.org/github.com/dnnrly/layli?status.svg)](http://godoc.org/github.com/dnnrly/layli)

![GitHub watchers](https://img.shields.io/github/watchers/dnnrly/layli?style=social)
![GitHub stars](https://img.shields.io/github/stars/dnnrly/layli?style=social)
[![Twitter URL](https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Fgithub.com%2Fdnnrly%2Flayli)](https://twitter.com/intent/tweet?url=https://github.com/dnnrly/layli)


## Using layli

### Installation

```
$ go install github.com/dnnrly/layli/cmd/layli
```

### A simple examples

Your first layli file:
```yml
node: "Hello World"
```

```bash
$ layli hello-world.layli
```

## layli principles

layli aims to let you specify nodes and edges (boxes and lines) and looks after arranging them in a pleasing way. If you've ever used [plantuml](https://plantuml.com) you'll be familiar with describing the diagrams in a simple to understand text file to generate a pretty diagram. Well, perhaps not as pretty as you would hope. This tool aims to solve this.

Here are some principles that hope to tackle this problem head on:

1. Nodes must be centered points that sit on a "node grid"
2. Edge paths must travel across a "path grid"
3. Node borders must be on the "path grid"
4. Edges must meet nodes at "ports"
5. Ports must sit on the "path grid"
6. Where a port is not specified, layli may select a "default" port on any side of the node
7. Nodes must be layed out so that the total area marked by the outside bounds of every node is as small as possible
8. Nodes must be seperated by at least 1 space on the "path grid"
9.  Edges must must not cross!
10. Edges must follow a grid path (ie. not curved or diagonal)
11. Edges must be as short as possible
12. Edges must have as few corners as possible
13. Edge paths may sit on top of each other at the beginning or at the end

### Defining nodes

Specifying a simple node:

```yml
node: "The node name
```

Complex nodes:
```yml
node:
  id: node-1
  contents: "A\nwith formatting"
```

Connecting nodes:

```yml
edge:
  from: node-1
  to: "The node name"
```

## Developing layli

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them

```bash
Give examples
```

### Installing from source

```bash
$ git clone http://github.com/dnnrly/layli.git
$ cd layli
$ make install
```

### Running Unit Tests

```bash
$ make test
```

### Running Acceptance tests

```bash
$ make acceptance-test
```

## Important `make` targets

* `deps` - downloads all of the deps you need to build, test, and release
* `install` - installs your application
* `build` - builds your application
* `test` - runs unit tests
* `ci-test` - run tests for CI validation
* `acceptance-test` - run the acceptance tests
* `coverage-report` - merge coverage statistics from all sources
* `mocks` - generate mocks for interface
* `lint` -  run linting
* `clean` - clean project dependencies
* `clean-deps` - remove all of the build dependencies too


## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/dnnrly/layli/tags). 

## Authors

* **Pascal Dennerly** - *Initial work* - [dnnrly](https://github.com/dnnrly)

See also the list of [contributors](https://github.com/dnnrly/layli/contributors) who participated in this project.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* There is a blog that I read a couple of years ago that described solving a similar problem. For the life of me, I can't remember find it anywhere to give the appropriate credit. But believe me when I say that a lot of the ideas are based on this blog and a lot of credit belongs to the author.
