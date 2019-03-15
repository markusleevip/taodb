package client

import (
	"flag"
)

type Options struct {
	port   		string
	logto      	string
	loglevel   	string
	ip			string
}

func parseArgs() *Options {
	ip :=flag.String("ip","127.0.0.1"," server ip")
	port := flag.String("port", ":7398", "")
	logto := flag.String("log", "stdout", "Write log messages to this file. 'stdout' and 'none' have special meanings")
	loglevel := flag.String("log-level", "DEBUG", "The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")
	flag.Parse()

	return &Options{
		port:   	*port,
		logto:     	*logto,
		loglevel:  	*loglevel,
		ip:			*ip,
	}
}
