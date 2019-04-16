package cluster

import "errors"

//Constants used to specify various conditions and state within the system
const (
	REPLICATION_MODE_FULL_REPLICATE byte = iota
	REPLICATION_MODE_SHARD

	clusterStateStarting byte = iota
	clusterStateStarted
	clusterStateEnded

	clusterModeACTIVE byte = iota
	clusterModePASSIVE
)

//CHAN_SIZE for all defined channels
const CHAN_SIZE = 1024 * 64

//Global errors that can be returned to callers
var (
	ErrNotEnoughReplica = errors.New("not enough replica")
	ErrNotFound         = errors.New("data not found")
	ErrTimedOut         = errors.New("not found as a result of timing out")
	ErrNotStarted       = errors.New("node not started, call Start()")
)