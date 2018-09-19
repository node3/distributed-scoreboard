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
	// Parse arguments
	args := os.Args[1:]
	server, listSize := parseArgs(args)

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

	highestScores := &utils.HighestScores {
		ListSize: 		0,
		MaxListSize: 	listSize,
		Records:		make([]utils.Data, listSize),
	}

	/* ********************* Track online players *************************** */
	go watchOnlineStatus(server, ch)

	/* ********************* Display score and status ************************* */
	for {
		updateMsg := <-ch
		if updateMsg.Type == utils.OnlineStatusUpdate {
			// Online status change arrives at channel
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
			// Score post update arrives at channel
			data := utils.Data {
				Score: updateMsg.Score,
				Player: updateMsg.Players,
			}
			recentScores.Push(data)
			highestScores.Push(data)
		}
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		recentScores.Display(online)
		fmt.Println("")
		highestScores.Display(online)
	}
}

func watchOnlineStatus(server string, ch chan utils.Update) {
	// Create a new connection for the new routine
	conn, _, err := zk.Connect([]string{server}, time.Second)
	utils.ExitIfError(err, "watchOnlineStatus: Could not connect to Zk server")
	defer conn.Close()

	// Keep watching forever
	for {
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
	utils.ExitIfError(err, "watchPlayerScores: Could not connect to Zk server")
	defer conn.Close()

	var data map[string]int64
	znodePath := utils.GetZnodePath(utils.ScoreDir, player)
	for {
		rawData, _, ech, err := conn.GetW(znodePath)
		utils.ExitIfError(err, "watchPlayerScores: Could not GetW score for player " + player)
		err = json.Unmarshal(rawData, &data)
		utils.ExitIfError(err, "watchPlayerScores: Could not unmarshal byte array to map")
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

func parseArgs(args []string) (string, int) {
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

	return server, listSize
}