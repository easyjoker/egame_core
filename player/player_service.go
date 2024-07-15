package player

import (
	"context"
	"egame_core"
	cache "egame_core/cache"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
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
func NewPlayer(client *redis.Client, account string, password string, name string, balance float64) (*Player, *egame_core.Error) {
	if nil == client {
		client = cache.GetClient()
	}

	pwd, err := GetPlayerAccountAndPassword(client, name)
	if err != nil {
		return nil, err
	}

	if pwd != "" {
		return nil, &egame_core.Error{Code: egame_core.PlayerExisted, Error: fmt.Errorf("player already exists")}
	}

	// 取得新玩家id
	id, err := GetPlayerId(client)

	if err != nil {
		return nil, err
	}

	// 保存玩家資料
	player := &Player{
		Id:       uint64(id),
		Account:  account,
		Password: password,
		Name:     name,
		Balance:  balance,
	}

	context := context.Background()
	var mapData map[string]interface{}
	parseError := mapstructure.Decode(player, &mapData)
	if parseError != nil {
		log.Println("NewPlayer error: ", err)
		return nil, &egame_core.Error{Code: egame_core.ParseError, Error: parseError}
	}

	result := client.HSet(context, fmt.Sprintf(PLAYER_ACCOUNT_KEY, name), mapData)
	if result.Err() != nil {
		log.Println("NewPlayer error: ", result.Err())
		return nil, &egame_core.Error{Code: egame_core.RedisError, Error: result.Err()}
	}

	return player, nil
}

// 取得玩家帳密
func GetPlayerAccountAndPassword(client *redis.Client, account string) (string, *egame_core.Error) {
	if nil == client {
		client = cache.GetClient()
	}
	context := context.Background()
	returned := client.Get(context, fmt.Sprintf(PLAYER_ACCOUNT_KEY, account))
	if returned.Err() != nil {
		if returned.Err() == redis.Nil {
			return "", nil
		}

		log.Println("GetPlayerAccount error: ", returned.Err())
		return "", &egame_core.Error{Code: egame_core.RedisError, Error: returned.Err()}
	}

	return returned.Val(), nil
}

// 取得玩家資料
func GetPlayer(client *redis.Client, id uint64) (player *Player, err *egame_core.Error) {
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

	if len(result) == 0 {
		return nil, nil
	}

	decodeErr := mapstructure.Decode(result, &player)
	if decodeErr != nil {
		log.Println("GetPlayer error: ", decodeErr)
		return nil, &egame_core.Error{Code: egame_core.ParseError, Error: decodeErr}
	}

	for key, value := range result {
		switch key {
		case "id":
			player.Id = id
		case "account":
			player.Account = value
		case "password":
			player.Password = value
		case "name":
			player.Name = value
		case "balance":
			float, parseError := strconv.ParseFloat(value, 64)
			if parseError != nil {
				log.Println("GetPlayer error: ", parseError)
				return nil, &egame_core.Error{Code: egame_core.ParseError, Error: parseError}
			}
			player.Balance = float
		}
	}

	return player, nil
}

// 取得可使用的玩家id
func GetPlayerId(client *redis.Client) (int64, *egame_core.Error) {
	if nil == client {
		client = cache.GetClient()
	}
	context := context.Background()
	returned := client.Incr(context, PLAYER_MAX_ID_KEY)
	id, err := returned.Result()
	if err != nil {
		log.Println("GetPlayerId error: ", err)
		return 0, &egame_core.Error{Code: egame_core.RedisError, Error: err}
	}

	if id == 0 {
		return GetPlayerId(client)
	}

	return id, nil
}

// get player balance from redis
func GetPlayerBalance(client *redis.Client, id uint64) (float64, error) {
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

// decut player balance to redis
func DecutPlayerBalance(client *redis.Client, id uint64, balance float64) *egame_core.Error {
	if nil == client {
		client = cache.GetClient()
	}
	if balance <= 0 {
		return &egame_core.Error{Code: egame_core.InvalidAmount, Error: fmt.Errorf("balance must be greater than 0")}
	}

	last, err := GetPlayerBalance(client, id)
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
func AddPlayerBalance(client *redis.Client, id uint64, balance float64) *egame_core.Error {
	if nil == client {
		client = cache.GetClient()
	}
	if balance <= 0 {
		return &egame_core.Error{Code: egame_core.NotEnoughBalance, Error: fmt.Errorf("balance must be greater than 0")}
	}

	last, err := GetPlayerBalance(client, id)
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
