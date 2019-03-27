#!/bin/bash
go fmt ./...
go build  ./main/server
go build  ./main/client