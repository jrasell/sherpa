## 0.4.2 (6 February, 2020)

IMPROVEMENTS:
 * Add pprof Debug API routes for Sherpa server profiling [[GH-124]](https://github.com/jrasell/sherpa/pull/124)

BUG FIXES:
 * Fix panic when performing autoscaling evaluation using external metrics [[GH-128]](https://github.com/jrasell/sherpa/pull/128)

## 0.4.1 (3 January, 2020)

BUG FIXES:
 * Fix incorrect searching of allocs causing missed allocations when performing autoscaling evaluations in jobs with a mix of scalable, and non-scable groups [[GH-121]](https://github.com/jrasell/sherpa/pull/121)

## 0.4.0 (20 December, 2019)

IMPROVEMENTS:
 * Scaling policies can now detail any number of external metrics and thresholds to check. The internal autoscaler will perform lookups on each iteration of the external source, to check whether scaling is needed. Currently `prometheus` is supported as an external source [[GH-104]](https://github.com/jrasell/sherpa/pull/104)
 * Metrics to measure policy backend request latencies [[GH-93]](https://github.com/jrasell/sherpa/pull/93)
 * Metrics to measure scaling state backend request latencies[[GH-94]](https://github.com/jrasell/sherpa/pull/94)
 * Metrics to measure autoscaler performance and health [[GH-95]](https://github.com/jrasell/sherpa/pull/95)
 * Add prometheus format option to metrics endpoint [[GH-101]](https://github.com/jrasell/sherpa/pull/101)
 * Scaling events can be listed by the latest per job group [[GH-107]](https://github.com/jrasell/sherpa/pull/107)
 * Policy list command output sorts alphabetically [[GH-108]](https://github.com/jrasell/sherpa/pull/108)
 * Order scale status list command output by time recent->oldest [[GH-110]](https://github.com/jrasell/sherpa/pull/110)

BUG FIXES:
 * Scale out/in endpoints should use POST rather than PUT method [[GH-103]](https://github.com/jrasell/sherpa/pull/103)
 * Reuse the Consul API client rather than create new for each use [[GH-106]](https://github.com/jrasell/sherpa/pull/106)
 * Fix system metric CLI listing samples incorrectly as counter [[GH-111]](https://github.com/jrasell/sherpa/pull/111)
 * Change incorrect description of `--storage-consul-enabled` CLI flag [[GH-116]](https://github.com/jrasell/sherpa/pull/116)

REFACTOR:
 * Move the system API endpoint into the `server/endpoints/v1` package [[GH-99]](https://github.com/jrasell/sherpa/pull/99)

## 0.3.0 (4 November, 2019)

IMPROVEMENTS:
 * The web UI now includes the count and direction of scaling events [[GH-75]](https://github.com/jrasell/sherpa/pull/75)
 * Custom metadata can now be added to scaling events as string key/value pairs [[GH-75]](https://github.com/jrasell/sherpa/pull/74)
 * Sherpa now supports specifying cooldowns within scaling policies. When attempting to perform a scaling activity, if the job group is in cooldown, the acitivty will be skipped [[GH-77]](https://github.com/jrasell/sherpa/pull/77)
 * The autoscaler provides more granular details regarding scaling decisions [[GH-86]](https://github.com/jrasell/sherpa/pull/86)

BUG FIXES:
 * Fix a big where cluster state was being overwritten on server restart [[GH-83]](https://github.com/jrasell/sherpa/pull/83)
 * Reorder autoscaling direction checks to ensure out actions take priority over in actions [[GH-84]](https://github.com/jrasell/sherpa/issues/84)

REFACTOR:
 * Fix incorrect variable name definition within policy meta [[GH-79]](https://github.com/jrasell/sherpa/pull/79)

## 0.2.1 (11 October, 2019)

IMPROVEMENTS:
 * Update vendored version of panjf2000/ants to latest v2 release [[GH-72]](https://github.com/jrasell/sherpa/pull/72)

BUG FIXES:
 * Fix malformed server log messages within the scale endpoint [[GH-69]](https://github.com/jrasell/sherpa/pull/69)
 * Fix issue when autoscaling a job with a mix of groups with and without scaling policies [[GH-71]](https://github.com/jrasell/sherpa/pull/71)

## 0.2.0 (10 October, 2019)

IMPROVEMENTS:
 * Sherpa server now performs clustering and leadership locking. The internal autoscaler, policy garbage collector, and API endpoints with the exception of `ui` and `system` are leader protected. The sub-process will not run unless the server is considered the cluster leader. Protected API endpoints called on a non-leader server will respond with a HTTP redirect [[GH-45]](https://github.com/jrasell/sherpa/pull/45)
 * Do not attempt to scale a job group if it is currently under deployment [[GH-56]](https://github.com/jrasell/sherpa/pull/56)
 * The Nomad meta policy engine now implements the policy backend interface and can run alongside the Consul storage backend [[GH-58]](https://github.com/jrasell/sherpa/pull/58)

BUG FIXES:
 * Fix a bug where the Consul backend only returned the last group policy for jobs with multiple groups [[GH-51]](https://github.com/jrasell/sherpa/pull/51)
 * The API for writing a policy should accept JobGroupPolicy type rather than a byte array [[GH-52]](https://github.com/jrasell/sherpa/pull/52)
 * Fix panic when attempting to read a policy by job and group name which doesn't exist in the Consul backend [[GH-53]](https://github.com/jrasell/sherpa/pull/53)
 * Fix a problem in the Nomad meta policy watcher which meant the process would perform a blocking query on new cluster [[GH-57]](https://github.com/jrasell/sherpa/pull/57)
 * Fix incorrect response code and msg when scaling breaks thresholds [[GH-62]](https://github.com/jrasell/sherpa/pull/62)

## 0.1.0 (17 September, 2019)

__BACKWARDS INCOMPATIBILITIES:__
 * The server `--storage-consul-path` CLI flag now defaults to `sherpa/` to accomodate policies and state. Setups using the previous default will not be impacted as the policies path appends `policies/`.

IMPROVEMENTS:
 * Use a goroutine pool within autoscaler to limit the number of concurrently running autoscaling threads [[GH-24]](https://github.com/jrasell/sherpa/pull/24)
 * Sherpa stores scaling event details within internal state and is viewable via the API and CLI [[GH-28]](https://github.com/jrasell/sherpa/pull/28)
 * Sherpa can now optionally run with a UI enabled, providing a visuale overview of scaling events [[GH-33]](https://github.com/jrasell/sherpa/pull/33)
 * Update system info to reference more generic storage backend [[GH-41]](https://github.com/jrasell/sherpa/pull/41)

BUG FIXES:
 * Use mutex read lock when reading out all policies from memory backend to remove possible race [[GH-30]](https://github.com/jrasell/sherpa/pull/30)
 * Autoscaler log to debug when it suggests a group should be scaled [[GH-36]](https://github.com/jrasell/sherpa/pull/36)
 * Fix style issue where scale cmd help ended with fullstops [[GH-39]](https://github.com/jrasell/sherpa/pull/39)

## 0.0.2 (9 August, 2019)

IMPROVEMENTS:
 * Scaling policies now support fields to manage internal scaling memory and CPU usage thresholds [[GH-9]](https://github.com/jrasell/sherpa/pull/9)
 
BUG FIXES:
 * Fix incorrect error message return within Consul backend [[GH-11]](https://github.com/jrasell/sherpa/pull/11)
 * Fix internal autoscaler calculations to miltiply first [[GH17]](https://github.com/jrasell/sherpa/pull/17)
 * Filter out the allocations from ResourceTracker that aren't actively running or pending [[GH-18]](https://github.com/jrasell/sherpa/pull/18)
 * Fix bug where meta policy engine would not process job updates [[GH-16]](https://github.com/jrasell/sherpa/pull/16)
 * Check error returned when calling Nomad Job Allocations in scaler rather than ignore it [[GH-21]](https://github.com/jrasell/sherpa/pull/21)

## 0.0.1 (18 May, 2019)

* Initial release.
