# Sherpa High Availability

Sherpa supports a multi-server mode for high availability. This mode protects against outages by running multiple Sherpa servers. High availability mode is automatically enabled when using a data store that supports it.

You can tell if a data store supports high availability mode by starting the server and seeing the `HAEnabled` return value from the `system/leader` endpoint. If it is, then Sherpa will automatically use HA mode.

To be highly available, one of the Sherpa server nodes grabs a lock within the data store. The successful server node then becomes the active node; all other nodes become standby nodes. At this point, if the standby nodes receive a request, they will redirect the client depending on the current configuration and state of the cluster. Due to this architecture, HA does not enable increased scalability. In general, the bottleneck of Sherpa is the data store itself, not Sherpa core.

## Client Redirection

The standby nodes will redirect the client using a 307 status code to the active node's redirect address.

What the `cluster-advertise-addr` value should be set to depends on how Sherpa is set up. There are two common scenarios: Sherpa servers accessed directly by clients, and Sherpa servers accessed via a load balancer.

In both cases, the `cluster-advertise-addr` should be a full URL including scheme (http/https), not simply an IP address and port.

### Direct Access

When clients are able to access Sherpa directly, the `cluster-advertise-addr` for each node should be that node's address. For instance, if there are two Sherpa nodes A (accessed via https://a.sherpa.mycompany.com:8000) and B (accessed via https://b.sherpa.mycompany.com:8000), node A would set its `cluster-advertise-addr` to https://a.sherpa.mycompany.com:8000 and node B would set its `cluster-advertise-addr` to https://b.sherpa.mycompany.com:8000.

This way, when A is the active node, any requests received by node B will cause it to redirect the client to node A's `cluster-advertise-addr` at https://a.sherpa.mycompany.com, and vice-versa.

### Behind Load Balancers
Sometimes clients use load balancers as an initial method to access one of the Sherpa servers, but actually have direct access to each Sherpa node. In this case, the Sherpa servers should actually be set up as described in the above section, since for redirection purposes the clients have direct access.

If the only access to the Sherpa servers is via the load balancer, the `cluster-advertise-addr` on each node should be the same: the address of the load balancer. Clients that reach a standby node will be redirected back to the load balancer; at that point hopefully the load balancer's configuration will have been updated to know the address of the current leader. This can cause a redirect loop and as such is not a recommended setup when it can be avoided.
