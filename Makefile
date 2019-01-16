.PHONY: bin/server
bin/server:
	GOPATH=$(CURDIR) go install master client server clientf clientread

.PHONY: run
run: bin/server
	-pkill -f bin/server
	-pkill -f bin/master
	-pkill -f bin/client
	bin/master & bin/server -port 7070 -n & bin/server -port 7071 -n & bin/server -port 7072 -n

master: bin/server
	-pkill -f bin/master
	bin/master