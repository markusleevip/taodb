package cluster

const (
	nodeStateConnecting = iota
	nodeStateConnected
	nodeStateDisconnected
	nodeStateHandshake
)

type remoteNodeState uint8