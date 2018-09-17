package main

import (
	"distributed-scoreboard/utils"
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"strconv"
	"time"
)

func main() {
	// Parse arguments
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Printf("Expected 2 arguments, got %d\n", len(args))
		os.Exit(1)
	}

	player := args[1]
	server := args[0]

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

	// Player goes online - create an ephemeral node
	onlineZnodePath := utils.GetZnodePath(utils.OnlineDir, player)
	znode, err = conn.Create(onlineZnodePath, []byte(player), utils.FlagEphemeral, acl)
	utils.ExitIfError(err, "Could not go online (ephemeral node not created)")
	fmt.Printf("Znode %s created. %s is online.\n", znode, player)

	// Register with scoreboard
	scoreZnodePath := utils.GetZnodePath(utils.ScoreDir, player)
	znode, err = conn.Create(scoreZnodePath, []byte("0"), utils.FlagRegular, acl)
	if znode != "" {
		fmt.Printf("Znode %s created.\n", znode)
	}
	fmt.Printf("%s can now post scores.\n", player)

	/* ************************* Send data ************************** */
	// Get the initial data from zookeeper
	data, stat, err := conn.Get(scoreZnodePath)
	utils.ExitIfError(err, "Could not get the initial score")
	tmp, _ := strconv.Atoi(string(data))
	score := int64(tmp)

	playerData := map[string]int64{"score": score, "timestamp": 0}

	// Post scores
	for true {
		// Update playerData
		score++
		playerData["score"] = score
		playerData["timestamp"] = time.Now().Unix()
		serializedPlayerData, err := json.Marshal(playerData)
		utils.ExitIfError(err, "Could not serialize playerData")

		// Send data to zk
		stat, err = conn.Set(scoreZnodePath, serializedPlayerData, stat.Version)
		utils.ExitIfError(err, "Could not post score to ZK")
		fmt.Printf("%d\n", playerData["score"])

		// sleep
		time.Sleep(3000 * time.Millisecond)
	}
}