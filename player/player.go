package main

import (
	"bufio"
	"distributed-scoreboard/utils"
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	/* ********************** Parse arguments ****************** */
	args := os.Args[1:]
	if len(args) != 2 && len(args) != 5 {
		fmt.Printf("Expected 2 or 5 arguments, got %d\n", len(args))
		os.Exit(1)
	}

	player := args[1]
	server := args[0]

	// Check for automated run or manual input
	manual := true
	var err error
	if len(args) == 5 {
		manual = false
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

	// Player goes online - creates an ephemeral node
	onlineZnodePath := utils.GetZnodePath(utils.OnlineDir, player)
	znode, err = conn.Create(onlineZnodePath, []byte(player), utils.FlagEphemeral, acl)
	utils.ExitIfError(err, "Could not go online (ephemeral node not created)")
	fmt.Printf("Znode %s created. %s is online.\n", znode, player)

	// Register with scoreboard
	scoreZnodePath := utils.GetZnodePath(utils.ScoreDir, player)
	znode, err = conn.Create(scoreZnodePath, []byte("0"), utils.FlagRegular, acl)
	fmt.Printf("Znode %s created. %s can now post scores.\n", scoreZnodePath, player)

	/* ************************* Send data ************************** */
	// Get the initial data from zookeeper
	data, stat, err := conn.Get(scoreZnodePath)
	utils.ExitIfError(err, "Could not get the initial score")
	tmp, _ := strconv.Atoi(string(data))
	score := int64(tmp)

	playerData := map[string]int64{"score": score, "timestamp": 0}

	// Post scores
	if manual == true {
		reader := bufio.NewReader(os.Stdin)
		for true{
			// Prompt user for score
			fmt.Print("Enter score: ")
			text, err := reader.ReadString('\n')
			utils.ExitIfError(err, "Could not read score from console")
			text = strings.TrimSpace(text)
			tmp, err = strconv.Atoi(string(text))
			utils.ExitIfError(err, "Could not convert score from string to int")

			// Update playerData data structure
			playerData["score"] = int64(tmp)
			playerData["timestamp"] = time.Now().Unix()
			serializedPlayerData, err := json.Marshal(playerData)
			utils.ExitIfError(err, "Could not serialize playerData")

			// Send data to zk
			stat, err = conn.Set(scoreZnodePath, serializedPlayerData, stat.Version)
			utils.ExitIfError(err, "Could not post score to ZK")
			fmt.Printf("Score %d posted to zookeeper.\n", playerData["score"])
		}
	} else {
		count, err := strconv.Atoi(args[2])
		utils.ExitIfError(err, "Could not convert count string input to int")
		u_delay, err := strconv.Atoi(args[3])
		utils.ExitIfError(err, "Could not convert u_delay string input to int")
		u_score, err := strconv.Atoi(args[4])
		utils.ExitIfError(err, "Could not convert u_score string input to int")
		fmt.Println(count, u_delay, u_score)
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
}