## Environment setup
* Copy the following code and save it as a shell script. Run the script.
```
sudo apt install golang-go -y
mkdir src
git clone https://github.ncsu.edu/atambol/distributed-scoreboard.git src/distributed-scoreboard
export GOPATH=`pwd`
go get github.com/samuel/go-zookeeper/zk
echo GOPATH=$GOPATH
```
* You will be prompted for github user-name and password. Please enter your credentials.
* This has download the github repo and installed golang plus its libraries. 
* The last line of the output gives the GOPATH environment variable that must be set before running any go commands.
* You can set the environment variable using export command as shown below. It must be set on every new terminal opened.

```export GOPATH=<Absolute path of parent directory of src directory>```

## Run player
```   
cd $GOPATH
go run src/distributed-scoreboard/player/player.go 127.0.0.1 Thor 
```

## Run watcher
```
cd $GOPATH
go run src/distributed-scoreboard/watcher/watcher.go 127.0.0.1 25
```
