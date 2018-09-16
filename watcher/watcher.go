package main

import (
	"distributed-scoreboard/utils"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"time"
)

func exitIfError(err error, msg string, player string) {
	if err != nil {
		fmt.Printf("%s %s. %s\n", player, msg, err)
		//panic(err)
		os.Exit(1)
	}
}

func main() {
	player := "Watcher"
	server := "127.0.0.1:2181"

	/* ********************** Register with server ****************** */
	// Connect to server
	conn, _, err := zk.Connect([]string{server}, time.Second)
	utils.ExitIfError(err, "Could not connect to Zk server")
	defer conn.Close()

	// Create the parent directory for online nodes
	acl := zk.WorldACL(zk.PermAll)
	znode, err := conn.Create(utils.OnlineDir, []byte("Online node directory"), utils.FlagRegular, acl)
	if znode != "" {
		fmt.Println("Online node directory created")
	}

	// Create the parent directory for scores
	znode, err = conn.Create(utils.ScoreDir, []byte("Scoreboard directory"), utils.FlagRegular, acl)
	if znode != "" {
		fmt.Println("Scoreboard directory created")
	}

	/* ********************* Keep Watching *************************** */
	children, _, err := conn.Children(utils.ScoreDir)
	exitIfError(err, "Could not get children", player)
	fmt.Printf("here 3")
	for _, name := range children {
		znodePath := utils.GetZnodePath(utils.ScoreDir, name)
		fmt.Println(znodePath)
		data, _, err := conn.Get(znodePath)
		exitIfError(err, "could not data from " + znodePath, player)
		fmt.Printf("%s stored at %s \n", data, znodePath)
		//err = conn.Delete("/dir/" + name, 0)
	}
	// Create a data structure to store score and online presence data
	// Get initial tree structure, set watches
	// Create a display function
	// Display the data and update on watch
}
