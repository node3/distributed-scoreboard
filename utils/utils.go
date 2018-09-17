package utils

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
)

const (
	// zookeeper constants
	FlagRegular = int32(0)
	FlagSequence = int32(zk.FlagSequence)
	FlagEphemeral = int32(zk.FlagEphemeral)
	OnlineDir = string("/online")
	ScoreDir = string("/score")

	// type of updates
	OnlineStatusUpdate = 0
	ScoreUpdate = 1

)

type Update struct {
	Type      int    `json:"type"`
	Players   string `json:"onlineNodes"`
	Score     int64  `json:"score"`
	Timestamp int64  `json:"timestamp"`
}

type Data struct {
	Player 	string		`json:"player"`
	Score	int64		`json:"score"`
}

type RecentScores struct {
	ListSize    int
	MaxListSize int
	Queue       chan Data
}

func (rs *RecentScores) Push(d Data) {
	if rs.ListSize == rs.MaxListSize {
		_ = <-rs.Queue
		rs.Queue <- d
	} else if rs.ListSize < rs.MaxListSize {
		rs.Queue <- d
		rs.ListSize++
	} else {
		panic("Push operation performed on channel beyond buffer size")
	}
}

func (rs *RecentScores) Pop() Data {
	if rs.ListSize <= 0 {
		panic("Pop operation performed on channel with size zero")
	} else {
		rs.ListSize--
		return <- rs.Queue
	}
}

func (rs *RecentScores) Display(online map[string]bool) {
	onlineStatus := ""
	listSize := rs.ListSize
	fmt.Println("Most recent scores")
	fmt.Println("------------------")
	for i := 0; i < listSize; i++ {
		d := rs.Pop()
		if _, ok := online[d.Player]; ok {
			onlineStatus = "**"
		} else {
			onlineStatus = ""
		}
		fmt.Printf("%-20s\t%d\t%s\n", d.Player, d.Score, onlineStatus)
		rs.Push(d)
	}
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