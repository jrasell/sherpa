# Sherpa Strict Checking Guide

Strict policy checking is a feature that can be enabled on a Sherpa server to provide tighter scaling limitations. When enabled, any scaling action on a job group must have a corresponding policy. The policy is then checked to ensure the following criteria are met:
* the scaling policy is enabled
* a scaling action will not break the min/max counts
* a scaling action will not break the cooldown

If either is broken, the scaling action will be declined.
