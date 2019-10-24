package scale

// JobGroupIsInCooldown satisfies the JobGroupIsInCooldown func within the Scale interface.
func (s *Scaler) JobGroupIsInCooldown(job, group string, cooldown int, time int64) (bool, error) {

	// Pull the latest scaling event for the job group out of the state.
	last, err := s.state.GetLatestScalingEvent(job, group)
	if err != nil {
		return true, err
	}

	// It is possible to return nil for the last event. This means that we were able to call the
	// backend successfully, but there is no latest event for the job group.
	if last == nil {
		return false, nil
	}

	if (time - int64(cooldown*1000000000)) < last.Time {
		return true, nil
	}
	return false, nil
}
