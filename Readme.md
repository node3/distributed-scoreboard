## Overview
This is distributed scoreboard application that keeps track of scores posted by players as well as their online status. The application uses a zookeeper service. The online status is maintained using the ephemeral Znodes and scores are maintained in regular nodes. Player program can be used to insert/post scores to zookeeper. Watchers display the scoreboards by implementing watches on player Znodes. One of the scoreboards displays most recent scores that is implemented suing channels. The other scoreboard is highscore list that is implemented using sorted array. Each player is individually watched by the watcher by using go routines.

The program expects a zookeeper service running on some IP and port that will be passed to the binaries as parameters in the form IP:PORT or just IP if the defualt 2181 port is used.

## Run player
* Manual Mode - input scores manually
```   
go run player.go <IP:PORT> <PlayerName> 
```

* Automated Mode - post scores automatically using Meanscore with some random standard deviations at an random interval with mean of MeanDelay. This will repeat for ScoreCount times.
```   
go run player.go <IP:PORT> <PlayerName> <ScoreCount> <MeanDelay> <MeanScore>
```

## Run watcher

```
go run watcher.go <IP:PORT> <ScoreBoardSize>
```

# Scoreboard output
The scoreboard shows two lists. One is the the recent scores and other is the highest score. Double asterisks signify that the player is online.
```
Most recent scores
------------------
Hulk                	248	**
IronMan             	134	**
Hulk                	155	**
IronMan             	79	**
Hulk                	107	**
IronMan             	104	**
Hulk                	61	**
Hulk                	180	**
IronMan             	207	**
Hulk                	76	**


Highest scores
------------------
THOR                	326	
THOR                	299	
THOR                	297	
THOR                	270	
Hulk                	251	**
Hulk                	248	**
THOR                	246	
IronMan             	225	**
THOR                	209	
IronMan             	207	**
```
