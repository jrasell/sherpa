package state

import "github.com/gofrs/uuid"

// ClusterInfo is our high level cluster information which holds unique identifiers for each
// cluster.
type ClusterInfo struct {

	// ID is the uuid assigned to this cluster.
	ID uuid.UUID

	// Name is the human friendly name of the cluster. It is designed to be used for human
	// identification of clusters as it is much easier than a UUID.
	Name string
}

// ClusterMember represents an individual members of a cluster and a member is considered eligible
// to perform leader actions.
type ClusterMember struct {

	// ID is the unique identifier of the Sherpa server instance.
	ID uuid.UUID

	// Addr is the Sherpa server address where the API is listening and present. This address
	// should include the protocol, host, and port.
	Addr string

	// AdvertiseAddr is the Sherpa server advertise address which can be used for NAT traversal
	// when redirecting requests to the cluster leader.
	AdvertiseAddr string
}
