package utils

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"time"
)

const (
	FlagRegular = int32(0)
	FlagSequence = int32(zk.FlagSequence)
	FlagEphemeral = int32(zk.FlagEphemeral)
	OnlineDir = string("/online")
	ScoreDir = string("/score")
)

func exitIfError(err error, msg string, player string) {
	if err != nil {
		fmt.Printf("%s %s. %s\n", player, msg, err)
		//panic(err)
		os.Exit(1)
	}
}

func register(server string, port string, player string) (*zk.Conn, []zk.ACL) {
	// Connect to server
	conn, _, err := zk.Connect([]string{server}, time.Second)
	exitIfError(err, "could not connect to Zk server", player)
	defer conn.Close()

	// Create the parent directory for online nodes
	acl := zk.WorldACL(zk.PermAll)
	znode, err := conn.Create(OnlineDir, []byte("Online node directory"), FlagRegular, acl)
	if znode != "" {
		fmt.Println("Online node directory created")
	}

	// Create the parent directory for scores
	znode, err = conn.Create(ScoreDir, []byte("Scoreboard directory"), FlagRegular, acl)
	if znode != "" {
		fmt.Println("Scoreboard directory created")
	}
}