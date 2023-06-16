# Go CLI Template

This is a template for Go CLI tools. Major features are:

1. Setup script
2. Release build action
3. PR validation action
4. Code of Conduct
5. Basic security policy
6. Modules enabled
7. Rudimentary accepance tests

## Setup

1. Create a new repo from this template
2. `$ ./setup.sh`
3. Follow the prompts

Use the `-d` option to see what will be modified without changing any files.

**You can delete everything above this line afterwards.**

# Project Title

One Paragraph of project description goes here

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/dnnrly/goclitem)](https://github.com/dnnrly/goclitem/releases/latest)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/dnnrly/goclitem/Release%20workflow)](https://github.com/dnnrly/goclitem/actions?query=workflow%3A%22Release+workflow%22)
[![report card](https://goreportcard.com/badge/github.com/dnnrly/goclitem)](https://goreportcard.com/report/github.com/dnnrly/goclitem)
[![godoc](https://godoc.org/github.com/dnnrly/goclitem?status.svg)](http://godoc.org/github.com/dnnrly/goclitem)

![GitHub watchers](https://img.shields.io/github/watchers/dnnrly/goclitem?style=social)
![GitHub stars](https://img.shields.io/github/stars/dnnrly/goclitem?style=social)
[![Twitter URL](https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Fgithub.com%2Fdnnrly%2Fgoclitem)](https://twitter.com/intent/tweet?url=https://github.com/dnnrly/goclitem)


## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them

```
Give examples
```

### Installing

```bash
$ git clone http://github.com/dnnrly/goclitem.git
$ cd goclitem
$ make install
```

### Running Unit Tests

```bash
$ make test
```

### Running Acceptance tests

```bash
$ make deps
$ make build acceptance-test
```

## Important `make` targets

* `deps` - downloads all of the deps you need to build, test, and release
* `install` - installs your application
* `build` - builds your application
* `test` - runs unit tests
* `ci-test` - run tests for CI validation
* `acceptance-test` - run the acceptance tests
* `lint` -  run linting
* `update` - update Go dependencies
* `clean` - clean project dependencies
* `clean-deps` - remove all of the build dependencies too


## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/dnnrly/goclitem/tags). 

## Authors

* **Your name here** - *Initial work* - [dnnrly](https://github.com/dnnrly)

See also the list of [contributors](https://github.com/dnnrly/goclitem/contributors) who participated in this project.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* Important people here
