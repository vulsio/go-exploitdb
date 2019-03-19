package db

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/inconshreveable/log15"
	"github.com/mozqnet/go-exploitdb/models"
	pb "gopkg.in/cheggaaa/pb.v1"
)

/**
# Redis Data Structure

- HASH
  ┌───┬───────────────────────┬────────────┬────────────────┬────────────────────────────────┐
  │NO │         HASH          │   FIELD    │     VALUE      │            PURPOSE             │
  └───┴───────────────────────┴────────────┴────────────────┴────────────────────────────────┘
  ┌───┬───────────────────────┬────────────┬────────────────┬────────────────────────────────┐
  │ 1 │EXPLOIT#E#$EXPLOITDBID │  $CVEID    │ $EXPLOIT JSON  │ TO GET EXPLOIT FROM EXPLOITDBID│
  ├───┼───────────────────────┼────────────┼────────────────┼────────────────────────────────┤
  │ 2 │EXPLOIT#C#$CVEID       │$EXPLOITDBID│ $EXPLOIT JSON  │ TO GET EXPLOIT FROM CVEID      │
  └───┴───────────────────────┴──────  ────┴────────────────┴────────────────────────────────┘
**/

const (
	dialectRedis      = "redis"
	exploitDBIDPrefix = "EXPLOIT#E#"
	cveIDPrefix       = "EXPLOIT#C#"
)

// RedisDriver is Driver for Redis
type RedisDriver struct {
	name string
	conn *redis.Client
}

// Name return db name
func (r *RedisDriver) Name() string {
	return r.name
}

// OpenDB opens Database
func (r *RedisDriver) OpenDB(dbType, dbPath string, debugSQL bool) (locked bool, err error) {
	if err = r.connectRedis(dbPath); err != nil {
		err = fmt.Errorf("Failed to open DB. dbtype: %s, dbpath: %s, err: %s", dbType, dbPath, err)
	}
	return
}

func (r *RedisDriver) connectRedis(dbPath string) error {
	var err error
	var option *redis.Options
	if option, err = redis.ParseURL(dbPath); err != nil {
		log15.Error("Failed to parse url.", "err", err)
		return err
	}
	r.conn = redis.NewClient(option)
	err = r.conn.Ping().Err()
	return err
}

// MigrateDB migrates Database
func (r *RedisDriver) MigrateDB() error {
	return nil
}

// GetExploitByCveID :
func (r *RedisDriver) GetExploitByCveID(cveID string) (exploits []*models.Exploit) {
	result := r.conn.HGetAll(cveIDPrefix + cveID)
	if result.Err() != nil {
		log15.Error("Failed to get cve.", "err", result.Err())
		return nil
	}

	for _, j := range result.Val() {
		var exploit models.Exploit
		if err := json.Unmarshal([]byte(j), &exploit); err != nil {
			log15.Error("Failed to Unmarshal json.", "err", err)
			return nil
		}
		exploits = append(exploits, &exploit)
	}
	return exploits
}

// GetExploitByID :
func (r *RedisDriver) GetExploitByID(exploitDBID string) (exploits []*models.Exploit) {
	results := r.conn.HGetAll(exploitDBIDPrefix + exploitDBID)
	if results.Err() != nil {
		log15.Error("Failed to get cve.", "err", results.Err())
		return nil
	}
	for _, j := range results.Val() {
		var exploit models.Exploit
		if err := json.Unmarshal([]byte(j), &exploit); err != nil {
			log15.Error("Failed to Unmarshal json.", "err", err)
			return nil
		}
		exploits = append(exploits, &exploit)
	}
	return exploits
}

// GetExploitAll :
func (r *RedisDriver) GetExploitAll() (exploits []*models.Exploit) {
	log15.Error("redis does not correspond to all")
	return
}

// GetExploitMultiByCveID :
func (r *RedisDriver) GetExploitMultiByCveID(cveIDs []string) (exploitsMap map[string][]*models.Exploit) {
	exploitsMap = map[string][]*models.Exploit{}
	rs := map[string]*redis.StringStringMapCmd{}

	pipe := r.conn.Pipeline()
	for _, cveID := range cveIDs {
		rs[cveID] = pipe.HGetAll(cveIDPrefix + cveID)
	}
	if _, err := pipe.Exec(); err != nil {
		if err != redis.Nil {
			log15.Error("Failed to get multi cve json.", "err", err)
			return nil
		}
	}

	for cveID, results := range rs {
		var exploits []*models.Exploit
		for _, j := range results.Val() {
			var exploit models.Exploit
			if results.Err() != nil {
				log15.Error("Failed to Get Explit", "err", results.Err())
				continue
			}
			if err := json.Unmarshal([]byte(j), &exploit); err != nil {
				log15.Error("Failed to Unmarshal json.", "err", err)
				return nil
			}
			exploits = append(exploits, &exploit)
		}
		exploitsMap[cveID] = exploits
	}
	return exploitsMap
}

//InsertExploit :
func (r *RedisDriver) InsertExploit(exploits []*models.Exploit) (err error) {
	bar := pb.StartNew(len(exploits))

	var noCveIDExploitCount, cveIDExploitCount int
	for _, exploit := range exploits {
		var pipe redis.Pipeliner
		pipe = r.conn.Pipeline()
		bar.Increment()

		j, err := json.Marshal(exploit)
		if err != nil {
			return fmt.Errorf("Failed to marshal json. err: %s", err)
		}

		if 0 < len(exploit.CveID) {
			if result := pipe.HSet(cveIDPrefix+exploit.CveID, exploit.ExploitUniqueID, string(j)); result.Err() != nil {
				return fmt.Errorf("Failed to HSet CVE. err: %s", result.Err())
			}
			cveIDExploitCount++
		} else {
			noCveIDExploitCount++
		}

		// NoCveID -> NONE
		if len(exploit.CveID) == 0 {
			exploit.CveID = "NONE"
		}
		if result := pipe.HSet(exploitDBIDPrefix+exploit.ExploitUniqueID, exploit.CveID, string(j)); result.Err() != nil {
			return fmt.Errorf("Failed to HSet Exploit. err: %s", result.Err())
		}

		if _, err = pipe.Exec(); err != nil {
			return fmt.Errorf("Failed to exec pipeline. err: %s", err)
		}
	}
	bar.Finish()
	log15.Info("No CveID Exploit Count", "count", noCveIDExploitCount)
	log15.Info("CveID Exploit Count", "count", cveIDExploitCount)
	return nil
}
