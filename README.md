# Sherpa

[![Build Status](https://travis-ci.org/jrasell/sherpa.svg?branch=master)](https://travis-ci.org/jrasell/sherpa) [![Go Report Card](https://goreportcard.com/badge/github.com/jrasell/sherpa)](https://goreportcard.com/report/github.com/jrasell/sherpa) [![GoDoc](https://godoc.org/github.com/jrasell/sherpa?status.svg)](https://godoc.org/github.com/jrasell/sherpa)

Sherpa is a job scaler for [HashiCorp Nomad](https://www.nomadproject.io/) and aims to be highly flexible so it can support a wide range of architectures and budgets.

## Download & Install

* The Sherpa binary can be downloaded from the [GitHub releases page](https://github.com/jrasell/sherpa/releases) using `curl -L https://github.com/jrasell/sherpa/releases/download/0.0.1/sherpa-linux-amd64 -o sherpa`

* A docker image can be found on [Docker Hub](https://hub.docker.com/r/jrasell/sherpa/), the latest version can be downloaded using `docker pull jrasell/sherpa`.

* Sherpa can be built from source by cloning the repository `git clone github.com/jrasell/sherpa.git` and then using the `make build` command. 

## Documentation

Please refer to the [documentation](./docs) directory for guides to help with deploying and using Sherpa in your Nomad setup.

## Contributing

Contributions to Sherpa are very welcome! Please reach out if you have any questions.
