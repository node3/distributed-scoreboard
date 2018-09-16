package main

import (
	"distributed-scoreboard/utils"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"strconv"
	"time"
)

func main() {
	// Initialize data structures
	var isOnline map[string]bool

	//server := "127.0.0.1:2181"
	// Parse arguments
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Printf("Expected 2 arguments, got %d\n", len(args))
		os.Exit(1)
	}

	server := args[0]
	displaySize, err := strconv.Atoi(args[1])
	utils.ExitIfError(err, "Could not convert size to integer.")
	if displaySize > 25 || displaySize < 1 {
		fmt.Printf("Display size expected in range [1, 25], got %d", displaySize)
		os.Exit(1)
	}

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

	/* ********************* Track online players *************************** */
	for true {
		children, _, ech, err := conn.ChildrenW(utils.OnlineDir)
		utils.ExitIfError(err, "Could not watch for online players")
		isOnline := make(map[string]bool)
		for _, name := range children {
			isOnline[name] = true
		}
		_ = <-ech
	}

	/* ********************* Keep Watching for scores *************************** */
	//for true {
	//	children, _, err := conn.Children(utils.ScoreDir)
	//	utils.ExitIfError(err, "Could not get children")
	//
	//	for _, name := range children {
	//		znodePath := utils.GetZnodePath(utils.ScoreDir, name)
	//		fmt.Println(znodePath)
	//		data, _, err := conn.Get(znodePath)
	//		utils.ExitIfError(err, "Could not data from "+znodePath)
	//		fmt.Printf("%s stored at %s \n", data, znodePath)
	//		//err = conn.Delete("/dir/" + name, 0)
	//	}
	//	time.Sleep(500 * time.Millisecond)
	//}
	// Create a data structure to store score and online presence data
	// Get initial tree structure, set watches
	// Create a display function
	// Display the data and update on watch
}


