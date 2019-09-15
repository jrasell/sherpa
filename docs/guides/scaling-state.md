# Sherpa Scaling State Guide

Scaling state can be accessed by the CLI, API, and UI, giving operators quick and easy insight into scaling events the Sherpa server has undertaken. Individual scaling events include details describing the changes that were made, the resulting Nomad evaluation ID, and the source of the request, whether it be the internal autoscaler or a request to the API.

## Garbage Collection

The scaling state is periodically garbage collected to ensure backend storage use does not grow indefinitely. When the GC process runs, it will remove all scaling events which were triggered over 24 hours ago.
