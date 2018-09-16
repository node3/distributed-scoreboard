package utils

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
)

const (
	FlagRegular = int32(0)
	FlagSequence = int32(zk.FlagSequence)
	FlagEphemeral = int32(zk.FlagEphemeral)
	OnlineDir = string("/online")
	ScoreDir = string("/score")
)

type Data struct {
	Score int64 	`json:"score"`
	Timestamp int64	`json:"timestamp"`
}


func GetZnodePath(dir string, player string) string {
	return dir + "/" + player
}

func ExitIfError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s. %s\n", msg, err)
		//panic(err)
		os.Exit(1)
	}
}

//func Register(server string, player string) (*zk.Conn, []zk.ACL) {
//
//}