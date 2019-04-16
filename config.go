package taodb

import "time"

type Config struct{
	Timeout time.Duration
	IdleTimeout time.Duration

	TCPKeepAlive time.Duration
}