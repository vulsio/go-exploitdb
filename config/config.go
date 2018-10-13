package config

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/inconshreveable/log15"
)

// CommonConfig :
type CommonConfig struct {
	Debug     bool
	DebugSQL  bool
	Quiet     bool
	Deep      bool
	DBPath    string
	DBType    string
	HTTPProxy string
}

// SearchConfig :
type SearchConfig struct {
	SearchType  string
	SearchParam string
}

// ServerConfig :
type ServerConfig struct {
	Bind string
	Port string
}

// CommonConf :
var CommonConf CommonConfig

// SearchConf :
var SearchConf SearchConfig

// ServerConf :
var ServerConf ServerConfig

// Validate :
func (p *CommonConfig) Validate() bool {
	if p.DBType == "sqlite3" {
		if ok, _ := valid.IsFilePath(p.DBPath); !ok {
			log15.Error("SQLite3 DB path must be a *Absolute* file path.", "dbpath", p.DBPath)
			return false
		}
	}

	_, err := valid.ValidateStruct(p)
	if err != nil {
		log15.Error("Invalid Struct", "err", err)
		return false
	}
	return true
}
