package mahjong

import (
	"encoding/json"
	"fmt"
)

// BroadcastRemainTile broadcasts remain tile
func (room Room) BroadcastRemainTile(num uint) {
	room.IO.BroadcastTo(room.Name, "remainTile", num)
}

// BroadcastStopWaiting broadcasts stop waiting signal
func (room Room) BroadcastStopWaiting() {
	room.IO.BroadcastTo(room.Name, "stopWaiting")
}

// BroadcastReady broadcasts the player's name who is ready
func (room Room) BroadcastReady(name string) {
	room.IO.BroadcastTo(room.Name, "broadcastReady", name)
}

// BroadcastGameStart broadcasts player list
func (room Room) BroadcastGameStart() {
	room.IO.BroadcastTo(room.Name, "broadcastGameStart", room.GetPlayerList())
}

// BroadcastWindAndRound broadcasts wind and round index
func (room Room) BroadcastWindAndRound(wind int, round int) {
	room.IO.BroadcastTo(room.Name, "broadcastWindAndRound", wind, round)
}

// BroadcastOpenDoor broadcasts poen door index and dice number
func (room Room) BroadcastOpenDoor(idx int, dice int) {
	room.IO.BroadcastTo(room.Name, "broadcastOpenDoor", idx, dice)
}

// BroadcastSetSeat broadcasts east seat idx
func (room Room) BroadcastSetSeat(idx int) {
	room.IO.BroadcastTo(room.Name, "broadcastSetSeat", idx)
}

// BroadcastBanker broadcasts banker ID
func (room Room) BroadcastBanker(id int, keepWin int) {
	room.IO.BroadcastTo(room.Name, "broadcastBanker", id, keepWin)
}

// BroadcastBuHua broadcasts the player's flower
func (room Room) BroadcastBuHua(flowers [][]string) {
	room.IO.BroadcastTo(room.Name, "broadcastBuHua", flowers)
}

// BroadcastHua broadcasts the player's draw flower
func (room Room) BroadcastHua(id int, tile Tile) {
	room.IO.BroadcastTo(room.Name, "broadcastHua", id, tile.ToString())
}

// BroadcastDraw broadcasts the player's id who draw a tile
func (room Room) BroadcastDraw(id int, num uint) {
	room.IO.BroadcastTo(room.Name, "broadcastDraw", id, num)
}

// BroadcastThrow broadcasts the player's id and the tile he threw
func (room Room) BroadcastThrow(id int, tile Tile) {
	room.IO.BroadcastTo(room.Name, "broadcastThrow", id, tile.ToString())
}

// BroadcastCommand broadcasts the player's id and the command he made
func (room Room) BroadcastCommand(from int, to int, command int, tile string, score int) {
	if command == COMMAND["ONGON"] {
		room.IO.BroadcastTo(room.Name, "broadcastCommand", from, to, command, "", score)
	} else {
		room.IO.BroadcastTo(room.Name, "broadcastCommand", from, to, command, tile, score)
	}
}

// BroadcastTing broadcasts the player's id which is ting
func (room Room) BroadcastTing(id int) {
	room.IO.BroadcastTo(room.Name, "broadcastTing", id)
}

// BroadcastEnd broadcasts the game result
func (room Room) BroadcastEnd(data []GameResult) {
	result, _ := json.Marshal(data)
	room.IO.BroadcastTo(room.Name, "end", string(result))
}

// BroadcastRobGon broadcasts rob gon
func (room Room) BroadcastRobGon(id int, tile Tile) {
	room.IO.BroadcastTo(room.Name, "robGon", id, tile.ToString())
}

// BroadcastGameEnd broadcasts when the game end
func (room Room) BroadcastGameEnd() {
	room.IO.BroadcastTo(room.Name, "gameEnd")
}

// BroadcastConversation broadcasts when player have something to say
func (room Room) BroadcastCoversation(id int, sentence string) {
	if id == -1 || sentence == "" {
		return
	}
	if id == 2 && len(sentence) > 6 {
		fmt.Println("player2 speak too long")
	}else {
		room.IO.BroadcastTo(room.Name, "speak", room.Players[id].ID, sentence)
	}
	
	// result := []Sentence{}
	// err := game.GossipDB.C("sentences").Find(bson.M{"Situation": action, "$or": []bson.M{bson.M{"Name": ""}, bson.M{"Name": "123"}}}).All(&result)
	// if err != nil {
	// 	(err)
	// 	return
	// }
	// fmt.Println(result[rand.Intn(len(result))].Sentence)
	// if result != nil {
	// 	room.IO.BroadcastTo(room.Name, "speak", id, result[rand.Intn(len(result))].Sentence)
	// }

}
