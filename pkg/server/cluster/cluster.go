package cluster

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jrasell/sherpa/pkg/state"
)

func (m *Member) setupCluster() error {

	c, err := m.clusterStorage.GetClusterInfo()
	if err != nil {
		return err
	}

	if c == nil {
		c = &state.ClusterInfo{}
	}

	// If existing state is found, check whether the operator passed cluster name matches the one
	// found. If the passed one does not match the state, we will fail with an error.
	if c.Name != "" && c.ID != uuid.Nil {
		m.logger.Debug().Msg("found existing Sherpa cluster state, verifying data")
		return m.verifyClusterName(c.Name)
	}

	if err := m.generateClusterName(); err != nil {
		return err
	}
	c.Name = m.clusterName

	id, err := m.generateClusterID()
	if err != nil {
		return err
	}
	c.ID = id
	m.logger.Debug().Str("id", c.ID.String()).Msg("successfully generated new cluster ID")

	return m.clusterStorage.PutClusterInfo(c)
}

func (m *Member) verifyClusterName(name string) error {
	if m.clusterName != "" {
		if m.clusterName != name {
			return errors.New("operator configured cluster name does not match discover state cluster name")
		}
	}

	m.logger.Info().Msg("successfully verified state of existing cluster to join")
	m.clusterName = name
	return nil
}

func (m *Member) generateClusterName() error {
	if m.clusterName == "" {
		m.logger.Debug().Msg("generating new cluster name")
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		m.clusterName = fmt.Sprintf("sherpa-%s", id.String())
		m.logger.Debug().Str("name", m.clusterName).Msg("successfully generated new cluster name")
	}
	return nil
}

func (m *Member) generateClusterID() (uuid.UUID, error) {
	m.logger.Debug().Msg("generating new Sherpa cluster ID")
	return uuid.NewV4()
}
