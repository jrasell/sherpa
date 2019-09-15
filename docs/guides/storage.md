# Sherpa Storage Guide

In order to persist policies and scaling state, Sherpa can be configured to use different types of storage backends.

### In-Memory

The in-memory backend is the default backend and ideal for development, but its use is not suggested for any path-to-production environments. Any Sherpa server failures or restarts will result in all data loss.

### Consul

Consul KV provides a scalable and robust backend store for Sherpa. All CRUD operations will be sanitized and then passed through for action within Consul using the official SDK. All data will be stored under the root KV as configured when running the Sherpa server, and can be browsed either using the Sherpa CLI, API or directly via Consul.

The Consul backend is preferable to in-memory as Sherpa server restarts or failures will not result in data loss. Instead the data relies on Consul distributed KV persistence which is proven as the highest scale.
