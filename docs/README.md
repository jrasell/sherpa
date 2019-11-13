# Sherpa Documentation

Sherpa is a highly available, fast, and flexible horizontal job scaling for [HashiCorp Nomad](https://www.nomadproject.io/). It is capable of running in a number of different modes to suit different requirements, and can scale based on Nomad resource metrics or external sources.

### Key Features
* __Scale jobs based on Nomad resource consumption and external metrics:__ The Sherpa autoscaler can use a mixture of Nomad resource checks, and external metric values to make scaling decisions. Both are optional to provide flexibility. Jobs can also be scaled via the CLI and API in either a manual manner, or by using webhooks sent from external applications such as Prometheus Alertmanager.
* __Highly available and fault tolerant:__ Sherpa performs leadership locking and quick fail-over, allowing multiple instances to run safely. During availability issues, or deployment Sherpa servers will gracefully handle leadership changes resulting in uninterrupted scaling. 
* __Operator friendly:__ Sherpa is designed to be easy to understand and work with as an operator. Scaling state in particular can contain metadata, providing insights into exactly why a scaling activity took place. A simple UI is also available to provide an easy method of checking scaling activities.

## Table of contents
1. [API](./api) documentation.
1. [CLI](./commands) documentation.
1. [Sherpa server](./configuration) configuration documentation.
1. [Guides](./guides) provide in-depth information on Sherpa behaviour, configuration, and workflows.
1. [Demos](./demos) provides and number of self contained examples to run through, allowing for better understanding of running Sherpa.
