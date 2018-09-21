make:
	mkdir -p src/distributed-scoreboard
	cp -r player utils watcher src/distributed-scoreboard/
	export GOPATH=`pwd`
	sudo apt install golang-go -y
	go get github.com/samuel/go-zookeeper/zk
	go install distributed-scoreboard/player
	go install distributed-scoreboard/watcher
	sudo cp $(GOPATH)/bin/* /usr/bin/

.ONESHELL:
