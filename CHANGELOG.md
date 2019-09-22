## 0.2.0 (Unreleased)

IMPROVEMENTS:
 * Sherpa server now performs clustering and leadership locking. The internal autoscaler, policy garbage collector, and API endpoints with the exception of `ui` and `system` are leader protected. The sub-process will not run unless the server is considered the cluster leader. Protected API endpoints called on a non-leader server will respond with a HTTP redirect. [[GH-45]](https://github.com/jrasell/sherpa/pull/45)

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
