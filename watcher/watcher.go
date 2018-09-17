package main

import (
	"distributed-scoreboard/utils"
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	//server := "127.0.0.1:2181"
	// Parse arguments
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Printf("Expected 2 arguments, got %d\n", len(args))
		os.Exit(1)
	}

	server := args[0]
	listSize, err := strconv.Atoi(args[1])
	utils.ExitIfError(err, "Could not convert size to integer.")
	if listSize > 25 || listSize < 1 {
		fmt.Printf("Display size expected in range [1, 25], got %d", listSize)
		os.Exit(1)
	}

	/* ********************** Register with server ****************** */
	// Connect to server
	conn, _, err := zk.Connect([]string{server}, time.Second)
	utils.ExitIfError(err, "Conn: Could not connect to Zk server")
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

	/* ********************* Initialize data structures ********************** */
	ch := make(chan utils.Update, 50)
	var online map[string]bool
	knownPlayers := make(map[string]bool)
	recentScores := &utils.RecentScores {
		ListSize: 		0,
		MaxListSize: 	listSize,
		Queue:			make(chan utils.Data, listSize),
	}

	/* ********************* Track online players *************************** */
	go watchOnlineStatus(server, ch)


	/* ********************* Display score and status ************************* */
	for true {
		updateMsg := <-ch
		if updateMsg.Type == utils.OnlineStatusUpdate {
			online = make(map[string]bool)
			if len(updateMsg.Players) > 0 {
				for _, player := range strings.Split(updateMsg.Players, ",") {
					online[player] = true
					if _, ok := knownPlayers[player]; !ok {
						knownPlayers[player] = true
						go watchPlayerScores(server, player, ch)
					}
				}
			}
		} else {
			data := utils.Data {
				Score: updateMsg.Score,
				Player: updateMsg.Players,
			}
			recentScores.Push(data)
		}
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		recentScores.Display(online)
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

func watchOnlineStatus(server string, ch chan utils.Update) {
	// Create a new connection for the new routine
	conn, _, err := zk.Connect([]string{server}, time.Second)
	utils.ExitIfError(err, "Conn: Could not connect to Zk server")
	defer conn.Close()

	// Keep watching forever
	for true {
		children, _, ech, err := conn.ChildrenW(utils.OnlineDir)
		utils.ExitIfError(err, "watchOnlineStatus: Could not watch for online players")
		// Get all the children, convert to a csv
		updateMsg := utils.Update{
			Type:      utils.OnlineStatusUpdate,
			Players:   strings.Join(children[:],","),
			Score:     0,
			Timestamp: 0,
		}

		// send into the channel from which display reads
		ch <- updateMsg

		// wait for zookeeper to send something into the channel
		_ = <-ech
	}
}


func watchPlayerScores(server string, player string, ch chan utils.Update) {
	// Create a new connection for the new routine
	conn, _, err := zk.Connect([]string{server}, time.Second)
	utils.ExitIfError(err, "Conn: Could not connect to Zk server")
	defer conn.Close()

	var data map[string]int64
	znodePath := utils.GetZnodePath(utils.ScoreDir, player)
	for true {
		rawData, _, ech, err := conn.GetW(znodePath)
		utils.ExitIfError(err, "watchPlayerScores: Could not GetW score for player " + player)
		json.Unmarshal(rawData, &data)
		utils.ExitIfError(err, "watchPlayerScores: Could not convert byte array to int64")
		// Send the main routine updated score
		updateMsg := utils.Update{
			Type:      utils.ScoreUpdate,
			Players:   player,
			Score:     data["score"],
			Timestamp: data["timestamp"],
		}

		// send into the channel from which display reads
		ch <- updateMsg

		// wait for zookeeper to send something into the channel
		_ = <-ech
	}
}

func display(online map[string]bool) {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
	//fmt.Println(online)
	for player := range online {
		fmt.Printf("%-20s\t*\n", player)
	}
}