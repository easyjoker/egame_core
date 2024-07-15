package egame_core

type ErrorCode int

// Error codes 除了 Success 之外，其他都是錯誤代碼都必須是 4 位數字
const (
	// Success indicates the operation was successful
	Success ErrorCode = 0
	// ParseError indicates the error occurred while parsing the request
	ParseError ErrorCode = 1

	// PlyaerNotFound indicates the player was not found
	PlayerNotFound ErrorCode = 1000
	// PlayerExisted indicates the player name is already existed
	PlayerExisted ErrorCode = 1001
	// NotEnoughBalance indicates the player does not have enough balance
	NotEnoughBalance ErrorCode = 1002
	// LoginFailed indicates the login failed
	LoginFailed ErrorCode = 1003
	// InvalidAmount indicates the amount is invalid
	InvalidAmount ErrorCode = 2000

	// RedisError indicates the Redis error
	RedisError ErrorCode = 4000
)

type Error struct {
	Code  ErrorCode `json:"code"`
	Error error     `json:"message"`
}
