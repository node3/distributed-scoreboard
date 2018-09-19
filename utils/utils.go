package utils

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"strings"
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

	// player automation
	MaxStdDev4Score = 100
	MinMeanScore    = 10
	MaxStdDev4Delay = 5
	MinMeanDelay    = 5
	MaxMeanDelay    = 10
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
	output := make([]string, rs.ListSize)
	fmt.Println("Most recent scores")
	fmt.Println("------------------")
	for i := rs.ListSize-1; i >= 0; i-- {
		d := rs.Pop()
		if _, ok := online[d.Player]; ok {
			onlineStatus = "**"
		} else {
			onlineStatus = ""
		}
		output[i] = fmt.Sprintf("%-20s\t%d\t%s\n", d.Player, d.Score, onlineStatus)
		rs.Push(d)
	}
	fmt.Println(strings.Join(output, ""))
}

type HighestScores struct {
	ListSize    int
	MaxListSize int
	Records     []Data
}

func (hs *HighestScores) Push(d Data) {
	if hs.ListSize == hs.MaxListSize {
		if d.Score < hs.Records[0].Score {
			return
		}

		for i := 1; i < hs.ListSize; i++ {
			if d.Score < hs.Records[i].Score {
				hs.Records[i-1] = d
				return
			} else {
				hs.Records[i-1] = hs.Records[i]
			}
		}
		hs.Records[hs.ListSize-1] = d
	} else {
		hs.ListSize++
		for i := hs.ListSize-2; i >= 0 ; i-- {
			if d.Score < hs.Records[i].Score {
				hs.Records[i+1] = hs.Records[i]
			} else {
				hs.Records[i+1] = d
				return
			}
		}
		hs.Records[0] = d
	}
}

func (hs *HighestScores) Display(online map[string]bool) {
	onlineStatus := ""
	fmt.Println("Highest scores")
	fmt.Println("------------------")
	for i := hs.ListSize-1; i >= 0; i-- {
		if _, ok := online[hs.Records[i].Player]; ok {
			onlineStatus = "**"
		} else {
			onlineStatus = ""
		}
		fmt.Printf("%-20s\t%d\t%s\n", hs.Records[i].Player, hs.Records[i].Score, onlineStatus)
	}
}

func GetZnodePath(dir string, player string) string {
	return dir + "/" + player
}

func ExitIfError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s. Error: %s.\nProgram will now terminate.\n", msg, err)
		//panic(err)
		os.Exit(1)
	}
}