# Sherpa

[![Build Status](https://travis-ci.org/jrasell/sherpa.svg?branch=master)](https://travis-ci.org/jrasell/sherpa) [![Go Report Card](https://goreportcard.com/badge/github.com/jrasell/sherpa)](https://goreportcard.com/report/github.com/jrasell/sherpa) [![GoDoc](https://godoc.org/github.com/jrasell/sherpa?status.svg)](https://godoc.org/github.com/jrasell/sherpa)

Sherpa is a highly available, fast, and flexible horizontal job scaling for [HashiCorp Nomad](https://www.nomadproject.io/). It is capable of running in a number of different modes to suit different requirements, and can scale based on Nomad resource metrics or external sources.

### Features
* __Scale jobs based on Nomad resource consumption and external metrics:__ The Sherpa autoscaler can use a mixture of Nomad resource checks, and external metric values to make scaling decisions. Both are optional to provide flexibility. Jobs can also be scaled via the CLI and API in either a manual manner, or by using webhooks sent from external applications such as Prometheus Alertmanager.
* __Highly available and fault tolerant:__ Sherpa performs leadership locking and quick fail-over, allowing multiple instances to run safely. During availability issues, or deployment Sherpa servers will gracefully handle leadership changes resulting in uninterrupted scaling. 
* __Operator friendly:__ Sherpa is designed to be easy to understand and work with as an operator. Scaling state in particular can contain metadata, providing insights into exactly why a scaling activity took place. A simple UI is also available to provide an easy method of checking scaling activities.

## Download & Install

* The Sherpa binary can be downloaded from the [GitHub releases page](https://github.com/jrasell/sherpa/releases) using `curl -L https://github.com/jrasell/sherpa/releases/download/v0.4.2/sherpa_0.4.2_linux_amd64 -o sherpa`

* A docker image can be found on [Docker Hub](https://hub.docker.com/r/jrasell/sherpa/), the latest version can be downloaded using `docker pull jrasell/sherpa`.

* Sherpa can be built from source by cloning the repository `git clone github.com/jrasell/sherpa.git` and then using the `make build` command. 

## Documentation

Please refer to the [documentation](./docs) directory for guides to help with deploying and using Sherpa in your Nomad setup.

## Contributing

Contributions to Sherpa are very welcome! Please reach out if you have any questions.

### Contributors

Thanks to everyone who has contributed to this project.

[@jvineet](https://github.com/jvineet) [@josegonzalez](https://github.com/josegonzalez) [@pmcatominey](https://github.com/pmcatominey) [@numiralofe](https://github.com/numiralofe) [@commarla](https://github.com/commarla) [@hobochili](https://github.com/hobochili)
