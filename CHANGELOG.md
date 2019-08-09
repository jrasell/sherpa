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
