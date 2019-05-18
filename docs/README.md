# Sherpa Documentation

Sherpa is a fast and flexible job scaler for [HashiCorp Nomad](https://www.nomadproject.io/) capable of running in a number of different modes to suit your needs.

## Table of contents
1. [API](./api) documentation
1. [CLI](./commands) documentation
1. [Sherpa server](./configuration) configuration documentation
1. [Guides](./guides) to provide additional information on Sherpa behaviour and workflows

## Server Run Modes

The Sherpa server can be configured in a number of ways, to provide flexibility and scalability across Nomad deployments. For all the configuration options, please take a look at the [server command](./commands/server.md) documentation. Below outlines the Sherpa server run variations and suggests where they could be most viable.

### Autoscaling Run Types

Sherpa can perform autoscaling or act as a proxy for the CLI or other external sources. Each run mode has different pros and cons as detailed below.

#### Scaling Proxy

Sherpa can act as a scaling proxy, taking requests via the scaling API endpoints and performing the required actions. The actions can be triggered by external sources such as [Prometheus AlertManager](https://prometheus.io/docs/alerting/alertmanager/), where rules are configured on telemetry data points. When an alert is triggered, the system can then send a Sherpa API request via webhooks. This solution is the most scalable, delegating the resource analysis to systems designed for this type of work. 


#### Built-in Autoscaler

The built-in autoscaler is ideal for smaller, development or cost limited setups. It runs on an internal time and will asses the resource usage of all job groups which have an active scaling policy. It is important to remember that the internal autoscaler will put additional load onto the Nomad servers. This is caused by the fact that analysing the memory and cpu consumption of a job requires X Nomad API calls.

### Policy Run Types

Policies are a method of controlling how and when job groups are autoscaled. When using the built-in autoscaler, strict checking is enabled which means job group will only be scaled if they have an associated scaling policy.

#### API Policy Engine

Scaling policies can be written, updated and deleted via the API and CLI. These policies are then stored in one of the available [backends](./guides/policies.md#policy-storage-backend) which has been enabled. The in-memory backend is not suitable for any environment other than development as the policies are lost when Sherpa is stopped. The Consul backend is ideal for non-dev environments and policies will be persisted after Sherpa restarts.

#### Nomad Job Meta Policy Engine

Sherpa can also pull scaling policies from Nomad jobs via the [meta stanza](https://www.nomadproject.io/docs/job-specification/meta.html). Job groups that you wish to be scalable, should be configured with the appropriate keys and values. Sherpa will automatically read these and update its scaling table when changes occur. It useful to note here, meta policies should be strictly configured inside the job group stanza. If the meta keys are configured at the job level, they will not be applied to all groups within the job. 
