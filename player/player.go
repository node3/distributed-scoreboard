package main

import (
	"bufio"
	"distributed-scoreboard/utils"
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)


func main() {
	/* ********************** Parse arguments ****************** */
	args := os.Args[1:]
	player, server, manual := parseArgs(args)

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
		for {
			// Prompt user for score
			fmt.Print("Enter score: ")
			text, err := reader.ReadString('\n')
			utils.ExitIfError(err, "Could not read score from console")
			text = strings.TrimSpace(text)
			tmp, err = strconv.Atoi(string(text))
			if err != nil {
				fmt.Println("Could not convert input from string to int. Please enter a numeric value.")
				continue
			}

			// Validate input score to be grater than equal to 0
			if tmp < 0 {
				fmt.Println("Please enter a non-negative value")
				continue
			}

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
		// Parse the additional parameters to automate the score generation
		count, err := strconv.Atoi(args[2])
		utils.ExitIfError(err, "Could not convert count string input to int")

		intMeanDelay, err := strconv.Atoi(args[3])
		utils.ExitIfError(err, "Could not convert meanDelay string input to int")
		meanDelay := int64(intMeanDelay)
		if meanDelay < utils.MinMeanDelay || meanDelay > utils.MaxMeanDelay {
			fmt.Printf("Please enter a meanDelay of value between [%d, %d]\n", utils.MinMeanDelay, utils.MaxMeanDelay)
			os.Exit(1)
		}

		intMeanScore, err := strconv.Atoi(args[4])
		utils.ExitIfError(err, "Could not convert meanScore string input to int")
		meanScore := int64(intMeanScore)
		//fmt.Println(count, meanDelay, meanScore)
		if meanScore < utils.MinMeanScore {
			fmt.Printf("Please enter a meanScore of value >= %d to allow a variety in scores\n", utils.MinMeanScore)
			os.Exit(1)
		}

		// define the standard deviations
		var stdDevScore, stdDevDelay, delay int64
		if meanScore - utils.MaxStdDev4Score >= 0 {
			stdDevScore = utils.MaxStdDev4Score
		} else {
			stdDevScore = meanScore
		}
		if meanDelay - utils.MaxStdDev4Delay >= 0 {
			stdDevDelay = utils.MaxStdDev4Delay
		} else {
			stdDevDelay = meanDelay
		}

		// define min and max score and delays
		minScore := meanScore - stdDevScore
		maxScore := meanScore + stdDevScore
		minDelay := meanDelay - stdDevDelay
		maxDelay := meanDelay + stdDevDelay

		for count > 0 {
			// Update playerData
			score = randomize(minScore, maxScore, stdDevScore, meanScore)
			delay = randomize(minDelay, maxDelay, stdDevDelay, meanDelay)
			playerData["score"] = score
			playerData["timestamp"] = time.Now().Unix()
			serializedPlayerData, err := json.Marshal(playerData)
			utils.ExitIfError(err, "Could not serialize playerData")

			// Send data to zk
			stat, err = conn.Set(scoreZnodePath, serializedPlayerData, stat.Version)
			utils.ExitIfError(err, "Could not post score to ZK")
			fmt.Printf("Score %d posted to zookeeper.\n", score)

			// sleep
			count--
			fmt.Printf("Sleeping for %s\n",time.Duration(delay) * time.Second)
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}
}

func parseArgs(args []string) (string, string, bool) {
	if len(args) != 2 && len(args) != 5 {
		fmt.Printf("Expected 2 or 5 arguments, got %d\n", len(args))
		os.Exit(1)
	}

	player := args[1]
	server := args[0]

	// Check for automated run or manual input
	manual := true

	if len(args) == 5 {
		manual = false
	}
	return player, server, manual
}

func randomize(min, max, stddev, mean int64) int64 {
	// Generate a number with normal probabilistic distribution
	num := int64(rand.NormFloat64() * float64(stddev) + float64(mean))
	if num < 0 {
		return 0
	}
	return num
	// Normalize the number to lie between min and max
	//return (max - min) * (num/math.MaxInt64) + min
}