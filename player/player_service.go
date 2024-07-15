package player

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/easyjoker/egame_core"
	cache "github.com/easyjoker/egame_core/cache"

	"github.com/go-redis/redis/v8"
)

const (
	// 目前玩家id最大的key
	PLAYER_MAX_ID_KEY = "player_max_id"

	// 玩家資料的key
	PLAYER_KEY = "player_%d"

	// 玩家帳密的key
	PLAYER_ACCOUNT_KEY = "player_account_%s"
)

// 使用 mysql_connection 產生新玩家
func newPlayer(client *redis.Client, account string, password string, name string, balance float64) (*Player, *egame_core.Error) {
	if nil == client {
		client = cache.GetClient()
	}

	exist, checkErr := client.Exists(context.Background(), fmt.Sprintf(PLAYER_ACCOUNT_KEY, account)).Result()

	if checkErr != nil {
		log.Println("NewPlayer error: ", checkErr)
		return nil, &egame_core.Error{Code: egame_core.RedisError, Error: checkErr}
	}
	if exist == 1 {
		return nil, &egame_core.Error{Code: egame_core.PlayerExisted, Error: fmt.Errorf("player already exists")}
	}

	// 取得新玩家id
	id, err := getPlayerId(client)

	if err != nil {
		return nil, err
	}

	// 保存玩家資料
	// player := &Player{
	// 	Id:       id,
	// 	Account:  account,
	// 	Password: password,
	// 	Name:     name,
	// 	Balance:  balance,
	// }

	mapData := map[string]interface{}{
		"id":       id,
		"account":  account,
		"password": password,
		"name":     name,
		"balance":  balance,
	}

	context := context.Background()
	// var mapData map[string]interface{}
	// parseError := mapstructure.Decode(player, &mapData)
	// if parseError != nil {
	// 	log.Println("NewPlayer error: ", err)
	// 	return nil, &egame_core.Error{Code: egame_core.ParseError, Error: parseError}
	// }

	result := client.HSet(context, fmt.Sprintf(PLAYER_KEY, id), mapData)
	if result.Err() != nil {
		log.Println("NewPlayer error: ", result.Err())
		return nil, &egame_core.Error{Code: egame_core.RedisError, Error: result.Err()}
	}

	client.Set(context, fmt.Sprintf(PLAYER_ACCOUNT_KEY, account), id, 0)

	return getPlayer(client, id)
}

func NewPlayer(account string, password string, name string, balance float64) (*Player, *egame_core.Error) {
	return newPlayer(cache.GetClient(), account, password, name, balance)

}

func LoginPlayer(account string, password string) (*Player, *egame_core.Error) {
	client := cache.GetClient()

	player_id, err := client.Get(context.Background(), fmt.Sprintf(PLAYER_ACCOUNT_KEY, account)).Result()

	if err != nil {
		return nil, &egame_core.Error{Code: egame_core.PlayerNotFound, Error: fmt.Errorf("player not found")}
	}
	if id, parseError := strconv.ParseInt(player_id, 10, 64); parseError != nil {
		log.Println("LoginPlayer error: ", parseError)
		return nil, &egame_core.Error{Code: egame_core.ParseError, Error: parseError}
	} else {
		player, getPlayerError := getPlayer(client, id)
		if getPlayerError != nil {
			return nil, getPlayerError
		}
		if player.Password == password {
			return player, nil
		}

		return nil, &egame_core.Error{Code: egame_core.LoginFailed, Error: fmt.Errorf("login failed")}
	}
}

func GetPlayer(id int64) (player *Player, err *egame_core.Error) {
	return getPlayer(cache.GetClient(), id)
}

// 取得玩家資料
func getPlayer(client *redis.Client, id int64) (player *Player, err *egame_core.Error) {
	if nil == client {
		client = cache.GetClient()
	}

	context := context.Background()
	result, e := client.HGetAll(context, fmt.Sprintf(PLAYER_KEY, id)).Result()
	if e != nil {
		if e == redis.Nil {
			return nil, nil
		}

		log.Println("GetPlayer error: ", e)
		return nil, &egame_core.Error{Code: egame_core.RedisError, Error: e}
	}
	player = &Player{}
	for key, value := range result {
		switch key {
		case "id":
			player.Id, _ = strconv.ParseInt(value, 10, 64)
		case "account":
			player.Account = value
		case "password":
			player.Password = value
		case "name":
			player.Name = value
		case "balance":
			player.Balance, _ = strconv.ParseFloat(value, 64)
		}
	}
	return player, nil
}

// 取得可使用的玩家id
func getPlayerId(client *redis.Client) (int64, *egame_core.Error) {
	if nil == client {
		client = cache.GetClient()
	}
	context := context.Background()
	returned := client.Incr(context, PLAYER_MAX_ID_KEY)
	id, err := returned.Result()
	if err != nil {
		log.Println("getPlayerId error: ", err)
		return 0, &egame_core.Error{Code: egame_core.RedisError, Error: err}
	}

	if id == 0 {
		return getPlayerId(client)
	}

	return id, nil
}
func GetPlayerBalance(id int64) (float64, error) {
	return getPlayerBalance(nil, id)
}

// get player balance from redis
func getPlayerBalance(client *redis.Client, id int64) (float64, error) {
	if nil == client {
		client = cache.GetClient()
	}

	context := context.Background()
	returned := client.HGet(context, fmt.Sprintf(PLAYER_KEY, id), "balance")
	if returned.Err() != nil {
		log.Println("GetPlayerBalance error: ", returned.Err())
		return 0, returned.Err()
	}

	balance, err := returned.Float64()
	if err != nil {
		log.Println("GetPlayerBalance error: ", err)
		return 0, err
	}

	return balance, nil
}

func DecutPlayerBalance(id int64, balance float64) *egame_core.Error {
	return decutPlayerBalance(nil, id, balance)

}

// decut player balance to redis
func decutPlayerBalance(client *redis.Client, id int64, balance float64) *egame_core.Error {
	if nil == client {
		client = cache.GetClient()
	}
	if balance <= 0 {
		return &egame_core.Error{Code: egame_core.InvalidAmount, Error: fmt.Errorf("balance must be greater than 0")}
	}

	last, err := getPlayerBalance(client, id)
	if err != nil {
		return &egame_core.Error{Code: egame_core.RedisError, Error: err}
	}

	if last < balance {
		return &egame_core.Error{Code: egame_core.NotEnoughBalance, Error: fmt.Errorf("balance is not enough")}
	}

	context := context.Background()
	result := client.HSet(context, fmt.Sprintf(PLAYER_KEY, id), "balance", balance-last)
	if result.Err() != nil {
		log.Println("SetPlayerBalance error: ", result.Err())
		return &egame_core.Error{Code: egame_core.RedisError, Error: result.Err()}
	}

	return nil
}

// add player balance to redis
func AddPlayerBalance(id int64, balance float64) *egame_core.Error {
	client := cache.GetClient()
	if balance <= 0 {
		return &egame_core.Error{Code: egame_core.NotEnoughBalance, Error: fmt.Errorf("balance must be greater than 0")}
	}

	last, err := getPlayerBalance(client, id)
	if err != nil {
		return &egame_core.Error{Code: egame_core.RedisError, Error: err}
	}

	context := context.Background()
	result := client.HSet(context, fmt.Sprintf(PLAYER_KEY, id), "balance", last+balance)
	if result.Err() != nil {
		log.Println("AddPlayerBalance error: ", result.Err())
		return &egame_core.Error{Code: egame_core.RedisError, Error: result.Err()}
	}

	return nil
}
