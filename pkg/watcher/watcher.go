package watcher

// Watcher is the interface to satisfy when creating a watcher to run blocking queries on an API
// and process updates as required.
type Watcher interface {

	// Run is used to run the watcher loop. The input channel should be the channel where objects
	// are sent to be acted upon when they are determined to be new.
	Run(updateChan chan interface{})
}

// IndexHasChange is used to check whether a returned blocking query has an updated index, compared
// to a tracked value.
func IndexHasChange(new, old uint64) bool {
	if new <= old {
		return false
	}
	return true
}

// MaxFound is used to determine which value passed is the greatest. This is used to track the most
// recently found highest index value.
func MaxFound(new, old uint64) uint64 {
	if new <= old {
		return old
	}
	return new
}
