package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mfajri11/family-catering-micro-service/auth-service/pkg/log"
	"gopkg.in/yaml.v2"
)

var (
	tmpCfg   config
	Cfg      config = config{}
	mockCfg  config = config{}
	callPath        = ""
)

func InitMock() {
	tmpCfg = Cfg
	Cfg = mockCfg
}

func DestroyMockConfig() {
	Cfg = tmpCfg
	tmpCfg = config{}
}

func init() {
	mustLoad()
}

type pg struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Username     string        `yaml:"username"`
	DatabaseName string        `yaml:"database-name"`
	MaxPoolSize  int32         `yaml:"max-pool-size"`
	MaxLifeTime  time.Duration `yaml:"max-lifetime"`
	MaxIdleTime  time.Duration `yaml:"max-idle-time"`
	Secret       string
}

func (pg pg) ConnString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", pg.Username, pg.Secret, pg.Host, pg.Port, pg.DatabaseName)
}

type app struct {
	Name                string    `yaml:"name"`
	Host                string    `yaml:"host"`
	Port                string    `yaml:"port"`
	BaseURLOwner        string    `yaml:"base-url-owner"`
	AccessTokenDuration int64     `yaml:"access-token-duration"`
	Security            *security `yaml:"security"`
	address             string
}

func (a app) Address() string {
	return a.address
}

type security struct {
	jwtKey  string `yaml:"jwt-key"`
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
}

func (s security) JWTKey() string {
	return s.jwtKey
}

func (s security) PubKey() *rsa.PublicKey {
	return s.pubKey
}

func (s security) PrivKey() *rsa.PrivateKey {
	return s.privKey
}

func (s *security) init() {
	certPrivFile, err := os.ReadFile(filepath.Join(callPath, "etc/cert/id_rsa.pem"))
	if err != nil {
		panic(err)
	}
	certPubFile, err := os.ReadFile(filepath.Join(callPath, "etc/cert/id_rsa.pub"))
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(certPrivFile)
	s.privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	if s.privKey == nil {
		panic("private key is nil")
	}

	block, _ = pem.Decode(certPubFile)
	s.pubKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	if s.pubKey == nil {
		panic("public key is nil")
	}

}

type redis struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DatabaseName int    `yaml:"database-name"`
	address      string
}

func (r redis) Address() string {
	return r.address
}

type config struct {
	App      app   `yaml:"app"`
	Postgres pg    `yaml:"postgres"`
	Redis    redis `yaml:"redis"`
}

func load() (config, error) {
	env := getenv("AUTH_ENV", "development")

	// we are now at where `MustLoad` is called, but we want to read config file at config/config.<env>.yaml based on environment
	configFile := "config." + env + ".yaml"
	// second return value of runtime.Caller(0) will output the path of caller function (which is in this file path .)
	_, fname, _, _ := runtime.Caller(0)
	callerPath := filepath.Join(fname, "..")
	callPath = filepath.Join(callerPath, "..")
	ConfigPath := filepath.Join(callerPath, configFile)

	// open the configuration file
	f, err := os.OpenFile(ConfigPath, os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		// TODO: wrap error
		return config{}, err
	}

	// TODO: make an option to read data from database in order to make a config change without deployment

	// do not forget for panic if close return an error
	defer func() {
		err := f.Close()
		if err != nil {
			// TODO: wrap error
			log.Error(err, "error close config file")
			panic(err)
		}
	}()

	Cfg = config{}
	Cfg.App.Security = &security{}
	err = yaml.NewDecoder(f).Decode(&Cfg)
	if err != nil {
		// TODO: wrap error
		return config{}, err
	}

	Cfg.App.Security.init()

	Cfg.App.address = fmt.Sprintf("%s:%s", Cfg.App.Host, Cfg.App.Port)
	Cfg.Redis.address = fmt.Sprintf("%s:%d", Cfg.Redis.Host, Cfg.Redis.Port)
	// err = godotenv.Load()
	// if err != nil {
	// 	// TODO: wrap error
	// 	return config{}, err
	// }

	// manually set .env value to config struct because it is the easiest for now
	Cfg.Postgres.Secret = getenv("POSTGRES_SECRET", "secret")
	Cfg.App.Security.jwtKey = getenv("JWT_SECRET_KEY", "secret")

	return Cfg, nil
}

func mustLoad() config {
	var err error
	Cfg, err = load()
	if err != nil {
		panic(err)
	}

	return Cfg
}

// func mustLoad() config {
// 	env := getenv("AUTH_ENV", "development")

// 	// we are now at where `MustLoad` is called, but we want to read config file at config/config.<env>.yaml based on environment
// 	configFile := "config." + env + ".yaml"
// 	// second return value of runtime.Caller(0) will output the path of caller function (which is in this file path .)
// 	_, fname, _, _ := runtime.Caller(0)
// 	callerPath := filepath.Join(fname, "..")
// 	fullConfigPath := filepath.Join(callerPath, configFile)

// 	// open the configuration file
// 	f, err := os.OpenFile(fullConfigPath, os.O_RDONLY|os.O_SYNC, 0)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// do not forget for panic if close return an error
// 	defer func() {
// 		err := f.Close()
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()

// 	Cfg = config{}
// 	err = yaml.NewDecoder(f).Decode(&Cfg)
// 	if err != nil {
// 		panic(err)
// 	}

// 	Cfg.App.address = fmt.Sprintf("%s:%s", Cfg.App.Host, Cfg.App.Port)
// 	Cfg.Redis.address = fmt.Sprintf("%s:%d", Cfg.Redis.Host, Cfg.Redis.Port)
// 	err = godotenv.Load()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// manually set .env value to config struct because it the most easier right now
// 	Cfg.Postgres.Secret = getenv("POSTGRES_SECRET", "secret")
// 	Cfg.App.jwtSecretKey = getenv("JWT_SECRET_KEY", "secret")

// 	return Cfg
// }

func getenv(key, fallback string) string {
	s, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return s
}

func Reload() error {
	var err error
	Cfg, err = load()
	if err != nil {
		// TODO: wrap error
		return err
	}

	return nil
}
