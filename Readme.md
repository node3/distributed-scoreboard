## Environment setup
* Please follow these instructions to setup the go environment and the repository.
* You need not clone this repository manually. The script below would appropriately setup everything as discussed below.
* Copy the following code and save it as a shell script (name it install.sh). 
```
mkdir src
git clone https://github.ncsu.edu/atambol/distributed-scoreboard.git src/distributed-scoreboard
sudo apt install golang-go -y
export GOPATH=`pwd`
go get github.com/samuel/go-zookeeper/zk
echo export GOPATH=$GOPATH
```
* Run this command to run the script `chmod 755 install.sh && ./install.sh`.
* You will be prompted for github user-name and password. Please enter your credentials.
* The script has downloaded the github repo and installed golang plus its libraries. 
* The last line of the output gives the command to export GOPATH environment variable that must be set on every terminal.

## Run player
* Ensure GOPATH is set first by running the export command.
```   
cd $GOPATH
go run src/distributed-scoreboard/player/player.go 127.0.0.1 Thor 100 10 500
```

## Run watcher
* Ensure GOPATH is set first by running the export command.
```
cd $GOPATH
go run src/distributed-scoreboard/watcher/watcher.go 127.0.0.1 25
```
