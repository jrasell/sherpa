# Quick Start Demo
This demo aims to have you running Sherpa and scaling a job in under 5 minutes on your local machine. This will allow you to get an understanding of how the application works, and also hopefully a sense of the power Sherpa can provide to your Nomad cluster.

In order to run the demo, you will need [Consul](https://www.consul.io/downloads.html), [Nomad](https://www.nomadproject.io/downloads.html) and [Sherpa](https://github.com/jrasell/sherpa/releases) downloaded and available to run locally. Once you have these available, you should start Consul and Nomad in local dev mode:
```
$ consul agent -dev
$ nomad agent -dev
```

Once the services startup, you can start the Sherpa process. Ideally in a path to production environment you would run Sherpa as a Nomad service, however, for simplicity we will run the raw binary. The startup flags tell Sherpa to use and perform the following actions:
 * use Consul as its storage backend
 * run the internal autoscaler which will perform resource utilisation checks and trigger scaling if required
 * allow policies to be configured via the API
 * run the Sherpa web UI

Once started, the Sherpa API will be available at http://127.0.0.1:8000, the default bind configuration.
```
$ sherpa server --storage-consul-enabled --autoscaler-enabled --policy-engine-api-enabled --log-level=debug --ui
```

With Sherpa successfully started, we need to run a job on Nomad which we can scale. Nomad provides an [init](https://www.nomadproject.io/docs/commands/job/init.html) utility to write out an example job which is ideal for this situation. 
```
$ nomad init
```

The job is configured to only run 1 instance of the redis container. In order to change this count to be higher and allow us to demonstrate scaling in, we can use this handy sed command:
```
$ sed "s/count = 1/count = 4/g" example.nomad > example-new.nomad
```

We can then register the job on the Nomad cluster:
```
$ nomad run example-new.nomad
```

Once the Nomad job is running, we need to configure a scaling policy for the job and group so the autoscaler has something to evaluate. The policy mostly uses default values apart from the `MinCount` which we set to a low value. You can write the policy using the Sherpa CLI or API, depending on your preference.
```
$ curl -X POST --data '{"Enabled":true,"MinCount":1}' http://127.0.0.1:8000/v1/policy/example/cache
```

From this point on the autoscaler will run every 60s (default period) and assess the resource consumption of the job example. If it believes the job groups are over or under utilised it will suggest a scaling action. If the scaling action does not break any configured thresholds, the updated job specification will be submitted to Nomad. Over the next 3 minutes you should see 3 scaling in events. You can track this via the Sherpa logs, or via the Sherpa UI which is available at http://localhost:8000/ui.

### Important Log Entries

Understanding the Sherpa logs also helps understand the process and feature set. Below are a number of excerpts from the Sherpa logs you have available along with a small description explaining them.

Sherpa performs leadership elections to ensure only one Sherpa instance performs critical tasks such as running the autoscaler. Shortly after startup you should see that the Sherpa instance has obtained cluster leadership and started the protected sub-process.
```
10:38AM INF HTTP server successfully listening addr=127.0.0.1:8000
10:38AM DBG server received leader update message leader-msg="obtained leadership"
10:38AM INF started scaling state garbage collector handler
10:38AM INF starting Sherpa internal auto-scaling engine
```

The autoscaler will log information about the calculations it made to help understand the internals. In this log we see the usage percentages are low and that the autoscaler suggests we scale the job group in. This is successfully triggered as the scaling request does not break any configured thresholds.
```
4:17PM DBG resource utilisation calculation cpu-usage-percentage=8 group=cache job=example mem-usage-percentage=0
4:17PM DBG added group scaling request job=example scaling-req={"count":1,"direction":"in","group":"cache"}
4:17PM INF successfully triggered autoscaling of job evaluation-id=2680e261-d651-3687-4404-fe6f674a50dd id=893e075d-a64d-4415-bd23-c6f73aa4f98f job=example
```

In this set of logs lines, we can see that again the autoscaler suggests that the job group be scaled in. It is found, however, that scaling the job group in would break the minimum threshold. Therefore the request will not be actioned.
```
10:40AM DBG resource utilisation calculation cpu-usage-percentage=11 group=cache job=example mem-usage-percentage=0
10:40AM DBG added group scaling request job=example scaling-req={"count":1,"direction":"in","group":"cache"}
10:40AM DBG scaling action will break job group minimum threshold group=cache job=example
```

When shutting down Sherpa, the server will perform a number of safety tasks. This includes waiting for any in flight autoscaling process to finish, and shutting down the leadership process allowing another instance to take over quickly.
```
4:17PM DBG autoscaler still has in-flight workers, will continue to check
4:17PM DBG exiting autoscaling thread as a result of shutdown request
4:17PM INF successfully drained autoscaler worker pool
4:17PM INF shutting down leadership handler cluster-member-id=ab9e3278-5a18-4965-9e83-ce97e9423e8f cluster-name=sherpa-2e651291-161d-4758-a12d-72294088214c
4:17PM DBG shutting down periodic leader refresh cluster-member-id=ab9e3278-5a18-4965-9e83-ce97e9423e8f cluster-name=sherpa-2e651291-161d-4758-a12d-72294088214c
4:17PM DBG shutting down leader elections cluster-member-id=ab9e3278-5a18-4965-9e83-ce97e9423e8f cluster-name=sherpa-2e651291-161d-4758-a12d-72294088214c
4:17PM INF successfully shutdown server and sub-processes
4:17PM INF HTTP server has been shutdown: http: Server closed
```
