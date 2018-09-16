package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"strconv"
	"time"
	"distributed-scoreboard/utils"
)

func getZnodePath(dir string, player string) string {
	return dir + "/" + player
}

func main() {
	player := "Thor"
	score := 0

	// Player goes online - create an ephemeral node
	onlineZnodePath := getZnodePath(utils.OnlineDir, player)
	znode, err = conn.Create(onlineZnodePath, []byte(player), utils.FlagEphemeral, acl)
	utils.exitIfError(err, "could not go online (ephemeral node not created)", player)
	fmt.Printf("Znode %s created. %s is online.\n", znode, player)

	// Register with scoreboard
	scoreZnodePath := getZnodePath(scoreDir, player)
	znode, err = conn.Create(scoreZnodePath, []byte(strconv.Itoa(score)), FlagRegular, acl)
	if znode != "" {
		fmt.Printf("Znode %s created.\n", znode)
	}
	fmt.Printf("%s can now post scores.\n", player)

	// Get the initial data
	data, stat, err := conn.Get(scoreZnodePath)
	exitIfError(err, "could not get the initial score", player)
	score, _ = strconv.Atoi(string(data))

	// Post scores
	for true {
		score++
		stat, err = conn.Set(scoreZnodePath, []byte(strconv.Itoa(score)), stat.Version)
		exitIfError(err, "could not post score", player)
		fmt.Println(strconv.Itoa(score))
	}
}