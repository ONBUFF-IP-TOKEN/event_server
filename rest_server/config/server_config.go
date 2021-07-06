package config

import (
	"sync"

	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
)

var once sync.Once
var currentConfig *ServerConfig

type TokenInfo struct {
	MainnetHost      string   `yaml:"mainnet_host"`
	ServerWalletAddr string   `yaml:"server_wallet_address"`
	ServerPrivateKey string   `yaml:"server_private_key"`
	TokenAddrs       []string `yaml:"token_address"`
}

type ApiAuth struct {
	AuthEnable        bool   `yaml:"auth_enable"`
	JwtSecretKey      string `yaml:"jwt_secret_key"`
	TokenExpiryPeriod int64  `yaml:"token_expiry_period"`
}

type ServerConfig struct {
	baseconf.Config `yaml:",inline"`

	MysqlDBAuth baseconf.DBAuth `yaml:"mysql_db_auth"`
	Token       TokenInfo       `yaml:"token_info"`
	Auth        ApiAuth         `yaml:"api_auth"`
}

func GetInstance(filepath ...string) *ServerConfig {
	once.Do(func() {
		if len(filepath) <= 0 {
			panic(baseconf.ErrInitConfigFailed)
		}
		currentConfig = &ServerConfig{}
		if err := baseconf.Load(filepath[0], currentConfig); err != nil {
			currentConfig = nil
		}
	})

	return currentConfig
}
