# Rest API

The Public API - Rest API will allow external customers to access to their data via REST endpoints. The data will have been stored by [Consumers](../consumers/README.md).

Additionally some translation/augmentation of API payloads might be performed to provide an API optimised for customer preferences.

## Tutorials

* Project structure: https://github.com/golang-standards/project-layout
* Go Basics: https://tour.golang.org/welcome/1
* Go-Swagger: https://github.com/go-swagger/go-swagger/blob/master/docs/tutorial/todo-list.md

## Installation

Note: This project can utilise Docker dependences only from the root directory of this project.  The below installation instructions are if you would like to use native tooling, and assumes Mac operating system.

### Go

* Install Go: https://golang.org/doc/install
* Add GOPATH environment variable and add to your path in your config file of choice (.profile, .bashrc, .zshrc) like so:
```
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$PATH
```

### Go Swagger

* Install dependencies 'go-swagger,' etc. as follows:
  * NOTE: golangci-lint@1.16.0 has an issue that requires it to be installed via brew (go get will fail). This should be reviewed with subsequent versions.
```
brew tap go-swagger/go-swagger 
brew install go-swagger
```

### Golangci Lint

```
brew install golangci/tap/golangci-lint
```

## Generating Models from the swagger

```
make generate-models
```
