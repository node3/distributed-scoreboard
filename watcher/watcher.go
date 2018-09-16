package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"time"
)

const (
	onlineDir = "/online"
	scoreDir = "/score"
	FlagRegular = int32(0)
)

func exitIfError(err error, msg string, player string) {
	if err != nil {
		fmt.Printf("%s %s. %s\n", player, msg, err)
		//panic(err)
		os.Exit(1)
	}
}

func getZnodePath(dir string, player string) string {
	return dir + "/" + player
}

func main() {
	player := "Watcher"

	// Connect to server
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	exitIfError(err, "could not connect to Zk server", player)
	defer conn.Close()

	// Create the parent directory for online node
	acl := zk.WorldACL(zk.PermAll)
	znode, err := conn.Create(onlineDir, []byte("Online node directory"), FlagRegular, acl)
	if znode != "" {
		fmt.Println("Online node directory created")
	}

	// Create the parent directory for scores
	znode, err = conn.Create(scoreDir, []byte("Scoreboard directory"), FlagRegular, acl)
	if znode != "" {
		fmt.Println("Scoreboard directory created")
	}

	children, _, err := conn.Children(scoreDir)
	exitIfError(err, "Could not get children", player)
	for _, name := range children {
		znodePath := getZnodePath(scoreDir, name)
		fmt.Println(znodePath)
		data, _, err := conn.Get(znodePath)
		exitIfError(err, "could not data from " + znodePath, player)
		fmt.Printf("/dir/%s: %s\n", name, string(data))
		//err = conn.Delete("/dir/" + name, 0)
	}
	// Create a data structure to store score and online presence data
	// Get initial tree structure, set watches
	// Create a display function
	// Display the data and update on watch
}