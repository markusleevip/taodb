module github.com/markusleevip/taodb

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190320223903-b7391e95e576
	golang.org/x/net => github.com/golang/net v0.0.0-20190320064053-1272bf9dcd53
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190322080309-f49334f85ddc
	golang.org/x/text => github.com/golang/text v0.3.0
)

require (
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/syndtr/goleveldb v1.0.0
)
