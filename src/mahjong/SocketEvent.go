package mahjong

import (
	"encoding/json"
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

// SocketError is callback of socket error event
func SocketError(so socketio.Socket, err error) {
	log.Println("error:", err)
}

// SocketConnect is callback of socket connect event
func SocketConnect(so socketio.Socket) {
	log.Println("on connection")

	so.Emit("auth")

	so.On("join", func(name string) (string, bool) {
		fmt.Println(name)
		if name == "" {
			return "", true
		}
		_uuid, _err := Login(name, &so)
		return _uuid, _err
	})

	so.On("auth", func(uuid string, room string) int {
		if uuid == "" {
			return -1
		}

		index := FindPlayerByUUID(uuid)
		if index == -1 {
			return -1
		}

		player := PlayerList[index]
		if (player.State&(MATCHED|READY|PLAYING)) != 0 && room != "" && player.Room == room {
			so.Join(room)
		}
		player.Socket = &so
		return player.State
	})

	so.On("gossip", func() string {
		game.GossipDealer = so
		// var record []ActionInfo
		// for i := 0; i < 5; i++ {
		// 	record = append(record, ActionInfo{10, "throw", 0, "f6"})
		// }
		// gossipInfo := &GossipInfo{"Throw", 10, SuitSet{10, 10, 10}, 10, 10,
		// 	true, false, []string{"o7", "o8"}, record, SuitSet{10, 10, 10}}
		// info, err := json.Marshal(gossipInfo)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return "json fail"
		// }
		// fmt.Printf("%T\n", info)
		// game.GossipDealer.Emit("gossipInfo", string(info), func(pid int, situation string) {
		// 	fmt.Println(pid, situation)
		// })

		// logInfo := &GossipLog{"123", record, "testtest"}
		// info, err = json.Marshal(logInfo)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return "json fail"
		// }
		// game.GossipDealer.Emit("logInfo", string(info), func(pid int, situation string) {
		// 	fmt.Println(pid, situation)
		// })
		return "connect"
	})

	so.On("ready", socketReady)
	so.On("getRoomInfo", getRoomInfo)
	so.On("getID", getID)
	so.On("getReadyPlayer", getReadyPlayer)
	so.On("getWindAndRound", getWindAndRound)
	so.On("getHand", getHand)
	so.On("getPlayerList", getPlayerList)
	so.On("getHandCount", getHandCount)
	so.On("getRemainCount", getRemainCount)
	so.On("getDoor", getDoor)
	so.On("getSea", getSea)
	so.On("getFlower", getFlower)
	so.On("getCurrentIdx", getCurrentIdx)
	so.On("getScore", getScore)
	so.On("getOpenIdx", getOpenIdx)
	so.On("getBanker", getBanker)
	so.On("getEastIdx", getEastIdx)
	so.On("getgetTing", getTing)
	so.On("manualInputMessage", manualInputMessage)

	so.On("disconnection", func() {
		log.Println("on disconnect")
		index := Disconnect(so)
		Logout(index)
	})
}

func socketReady(uuid string, room string) int {
	if !Auth(room, uuid) {
		return -1
	}

	c := make(chan int, 1)
	fn := func(id int) {
		c <- id
	}
	go game.Rooms[room].Accept(uuid, fn)
	return <-c
}

func getRoomInfo(uuid string) (string, []string, bool) {
	if uuid == "" {
		return "", []string{}, true
	}
	index := FindPlayerByUUID(uuid)
	if index == -1 {
		return "", []string{}, true
	}
	player := PlayerList[index]
	room := player.Room
	return room, GetNameList(FindPlayerListInRoom(room, -1)), false
}

func getReadyPlayer(room string) []string {
	if game.Rooms[room] == nil {
		return []string{}
	}
	return game.Rooms[room].GetReadyPlayers()
}

func getWindAndRound(room string) (int, int) {
	if game.Rooms[room] == nil {
		return -1, -1
	}
	return game.Rooms[room].GetWindAndRound()
}

func getHand(uuid string, room string) []string {
	if !Auth(room, uuid) || game.Rooms[room].State < DealTile {
		return []string{}
	}
	index := FindPlayerByUUID(uuid)
	id := PlayerList[index].Index
	return game.Rooms[room].Players[id].Hand.ToStringArray()
}

func getID(uuid string, room string) int {
	if !Auth(room, uuid) {
		return -1
	}
	index := FindPlayerByUUID(uuid)
	if PlayerList[index].State != READY {
		return -1
	}
	return PlayerList[index].Index
}

func getPlayerList(room string) []string {
	if game.Rooms[room] == nil {
		return []string{}
	}
	return game.Rooms[room].GetPlayerList()
}

func getHandCount(room string) []int {
	if game.Rooms[room] == nil {
		return []int{}
	}
	return game.Rooms[room].GetHandCount()
}

func getRemainCount(room string) int {
	if game.Rooms[room] == nil {
		return 56
	}
	return game.Rooms[room].GetRemainCount()
}

func getDoor(uuid string, room string) ([][][]string, []int, bool) {
	if !Auth(room, uuid) {
		return [][][]string{}, []int{}, true
	}
	index := FindPlayerByUUID(uuid)
	id := PlayerList[index].Index
	return game.Rooms[room].GetDoor(id)
}

func getSea(room string) ([][]string, bool) {
	if game.Rooms[room] == nil {
		return [][]string{}, true
	}
	return game.Rooms[room].GetSea()
}

func getFlower(room string) ([][]string, bool) {
	if game.Rooms[room] == nil {
		return [][]string{}, true
	}
	return game.Rooms[room].GetFlower()
}

func getCurrentIdx(room string) int {
	if game.Rooms[room] == nil {
		return -1
	}
	return game.Rooms[room].GetCurrentIdx()
}

func getScore(room string) []int {
	if game.Rooms[room] == nil {
		return []int{}
	}
	return game.Rooms[room].GetScore()
}

func getOpenIdx(room string) int {
	if game.Rooms[room] == nil {
		return -1
	}
	return game.Rooms[room].OpenIdx
}

func getBanker(room string) (int, int) {
	if game.Rooms[room] == nil {
		return -1, -1
	}
	return game.Rooms[room].Banker, game.Rooms[room].NumKeepWin
}

func getEastIdx(room string) int {
	if game.Rooms[room] == nil {
		return -1
	}
	return game.Rooms[room].EastIdx
}

func getTing(room string) []bool {
	if game.Rooms[room] == nil {
		return []bool{}
	}
	return game.Rooms[room].GetTing()
}

func manualInputMessage(room string, sentence string, id int, uuid string) {
	fmt.Println(room, sentence, id, uuid)
	var name string
	var actions []ActionInfo
	for _, player := range game.Rooms[room].Players {
		if id == player.ID {
			name = player.UUID
			actions = player.Actions
		}
	}
	if len(actions) > 16 {
		actions = actions[len(actions)-16:]
	}
	logInfo := &GossipLog{name, actions, sentence}
	info, err := json.Marshal(logInfo)
	if err != nil {
		fmt.Println(err)
	}
	game.GossipDealer.Emit("logInfo", string(info))
	game.Rooms[room].IO.BroadcastTo(game.Rooms[room].Name, "speak", id, sentence)
}
