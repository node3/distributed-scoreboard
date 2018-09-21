make:
	mkdir -p src/distributed-scoreboard
	cp -r player utils watcher src/distributed-scoreboard/
	export GOPATH=`pwd`
	sudo apt install golang-go -y
	go get github.com/samuel/go-zookeeper/zk
	go install distributed-scoreboard/player
	go install distributed-scoreboard/watcher
	sudo cp bin/watcher /usr/bin/watcher
	sudo cp bin/player /usr/bin/player

.ONESHELL:
