package player

import (
	"encoding/json"
	"reflect"
	"testing"
)

// 測試玩家的行為
func TestNewPlayer(t *testing.T) {
	player := Player{Id: 1, Name: "test", Balance: 100}

	marshed, err := json.Marshal(player)

	if err != nil {
		t.Errorf("Player marshed has error")
	}

	var restored Player
	err = json.Unmarshal(marshed, &restored)
	if err != nil {
		t.Errorf("Player unmarshed is not correct")
	}

	if !reflect.DeepEqual(player, restored) {
		t.Errorf("Player unmarshed is not the same as original player")
	}
}
