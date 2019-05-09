# AutoScaler

Sherpa can run using an internal autoscaler which uses Nomad job resource utilization stats to determine whether a group is eligible for scaling. The autoscaler will iterate stored policies, performing calculations to figure out each job group CPU and memory consummation. If consumption is above 80% or below 20%, Sherpa will request scaling of the job group as long as the scaling policies limits are not violated.
