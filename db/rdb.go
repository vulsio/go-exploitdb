package db

import (
	"fmt"

	"github.com/cheggaaa/pb"
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	sqlite3 "github.com/mattn/go-sqlite3"

	// Required MySQL.  See http://jinzhu.me/gorm/database.html#connecting-to-a-database
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	// Required SQLite3.
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mozqnet/go-exploitdb/models"
	"github.com/mozqnet/go-exploitdb/util"
)

const (
	dialectSqlite3    = "sqlite3"
	dialectMysql      = "mysql"
	dialectPostgreSQL = "postgres"
)

// RDBDriver :
type RDBDriver struct {
	name string
	conn *gorm.DB
}

// Name return db name
func (r *RDBDriver) Name() string {
	return r.name
}

// OpenDB opens Database
func (r *RDBDriver) OpenDB(dbType, dbPath string, debugSQL bool) (locked bool, err error) {
	r.conn, err = gorm.Open(dbType, dbPath)
	if err != nil {
		msg := fmt.Sprintf("Failed to open DB. dbtype: %s, dbpath: %s, err: %s", dbType, dbPath, err)
		if r.name == dialectSqlite3 {
			switch err.(sqlite3.Error).Code {
			case sqlite3.ErrLocked, sqlite3.ErrBusy:
				return true, fmt.Errorf(msg)
			}
		}
		return false, fmt.Errorf(msg)
	}
	r.conn.LogMode(debugSQL)
	if r.name == dialectSqlite3 {
		r.conn.Exec("PRAGMA foreign_keys = ON")
	}
	return false, nil
}

// MigrateDB migrates Database
func (r *RDBDriver) MigrateDB() error {
	if err := r.conn.AutoMigrate(
		&models.Exploit{},
		&models.OffensiveSecurity{},
		&models.Document{},
		&models.ShellCode{},
		&models.GitHubRepository{},
	).Error; err != nil {
		return fmt.Errorf("Failed to migrate. err: %s", err)
	}

	var errs gorm.Errors
	errs = errs.Add(r.conn.Model(&models.Exploit{}).AddIndex("idx_exploit_exploit_cve_id", "cve_id").Error)
	errs = errs.Add(r.conn.Model(&models.OffensiveSecurity{}).AddIndex("idx_offensive_secyrity_exploit_unique_id", "exploit_unique_id").Error)
	errs = errs.Add(r.conn.Model(&models.Document{}).AddIndex("idx_exploit_document_exploit_unique_id", "exploit_unique_id").Error)
	errs = errs.Add(r.conn.Model(&models.ShellCode{}).AddIndex("idx_exploit_shell_code_exploit_unique_id", "exploit_unique_id").Error)
	errs = errs.Add(r.conn.Model(&models.GitHubRepository{}).AddIndex("idx_exploit_github_repository_exploit_unique_id", "exploit_unique_id").Error)

	for _, e := range errs {
		if e != nil {
			return fmt.Errorf("Failed to create index. err: %s", e)
		}
	}
	return nil
}

// InsertExploit :
func (r *RDBDriver) InsertExploit(exploits []*models.Exploit) (err error) {
	log15.Info(fmt.Sprintf("Inserting %d Exploits", len(exploits)))
	return r.deleteAndInsertExploit(r.conn, exploits)
}

func (r *RDBDriver) deleteAndInsertExploit(conn *gorm.DB, exploits []*models.Exploit) (err error) {
	bar := pb.StartNew(len(exploits))
	tx := conn.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	old := models.Exploit{}
	result := tx.Where(&models.Exploit{}).First(&old)
	if !result.RecordNotFound() {
		// Delete all old records
		var errs gorm.Errors
		errs = errs.Add(tx.Delete(models.Document{}).Error)
		errs = errs.Add(tx.Delete(models.ShellCode{}).Error)
		errs = errs.Add(tx.Delete(models.OffensiveSecurity{}).Error)
		errs = errs.Add(tx.Delete(models.GitHubRepository{}).Error)
		errs = errs.Add(tx.Delete(models.Exploit{}).Error)
		errs = util.DeleteNil(errs)
		if len(errs.GetErrors()) > 0 {
			return fmt.Errorf("Failed to delete old records. err: %s", errs.Error())
		}
	}

	var noCveIDExploitCount, cveIDExploitCount int
	for _, exploit := range exploits {
		if err = tx.Create(exploit).Error; err != nil {
			return fmt.Errorf("Failed to insert. exploitTypeID: %s, err: %s", exploit.ExploitUniqueID, err)
		}
		if 0 < len(exploit.CveID) {
			cveIDExploitCount++
		} else {
			noCveIDExploitCount++
		}
		bar.Increment()
	}
	bar.Finish()
	log15.Info("No CveID Exploit Count", "count", noCveIDExploitCount)
	log15.Info("CveID Exploit Count", "count", cveIDExploitCount)
	return nil
}

// GetExploitByID :
func (r *RDBDriver) GetExploitByID(exploitUniqueID string) []*models.Exploit {
	es := []*models.Exploit{}
	var errs gorm.Errors
	errs = errs.Add(r.conn.Where(&models.Exploit{ExploitUniqueID: exploitUniqueID}).Find(&es).Error)
	for _, e := range es {
		switch e.ExploitType {
		case models.OffensiveSecurityType:
			os := &models.OffensiveSecurity{}
			errs = errs.Add(r.conn.Preload("Document").Preload("ShellCode").Where(&models.OffensiveSecurity{ExploitUniqueID: e.ExploitUniqueID}).First(&os).Error)
			e.OffensiveSecurity = os

		case models.GitHubRepositoryType:
			gh := &models.GitHubRepository{}
			errs = errs.Add(r.conn.Where(&models.GitHubRepository{ExploitUniqueID: e.ExploitUniqueID}).First(&gh).Error)
			e.GitHubRepository = gh
		}
	}
	for _, e := range errs.GetErrors() {
		if !gorm.IsRecordNotFoundError(e) {
			log15.Error("Failed to get exploit by ExploitDB-ID", "err", e)
		}
	}
	return es
}

// GetExploitAll :
func (r *RDBDriver) GetExploitAll() []*models.Exploit {
	es := []*models.Exploit{}
	docs := []*models.Document{}
	shells := []*models.ShellCode{}
	offensiveSecurities := []*models.OffensiveSecurity{}
	var errs gorm.Errors

	errs = errs.Add(r.conn.Find(&es).Error)
	errs = errs.Add(r.conn.Find(&offensiveSecurities).Error)
	errs = errs.Add(r.conn.Find(&docs).Error)
	errs = errs.Add(r.conn.Find(&shells).Error)
	if len(errs.GetErrors()) > 0 {
		log15.Error("Failed to delete old records", "err", errs.Error())
	}

	for _, e := range es {
		for _, o := range offensiveSecurities {
			for _, d := range docs {
				if o.ID == d.OffensiveSecurityID {
					o.Document = d
				}
			}
			for _, s := range shells {
				if o.ID == s.OffensiveSecurityID {
					o.ShellCode = s
				}
			}
			if e.ID == o.ExploitID {
				e.OffensiveSecurity = o
			}
		}
	}
	return es
}

// GetExploitMultiByID :
func (r *RDBDriver) GetExploitMultiByID(exploitUniqueIDs []string) map[string][]*models.Exploit {
	exploits := map[string][]*models.Exploit{}
	for _, exploitUniqueID := range exploitUniqueIDs {
		exploits[exploitUniqueID] = r.GetExploitByID(exploitUniqueID)
	}
	return exploits
}

// GetExploitByCveID :
func (r *RDBDriver) GetExploitByCveID(cveID string) []*models.Exploit {
	es := []*models.Exploit{}
	var errs gorm.Errors
	errs = errs.Add(r.conn.Where(&models.Exploit{CveID: cveID}).Find(&es).Error)
	for _, e := range es {
		switch e.ExploitType {
		case models.OffensiveSecurityType:
			os := &models.OffensiveSecurity{}
			errs = errs.Add(r.conn.Preload("Document").Preload("ShellCode").Where(&models.OffensiveSecurity{ExploitUniqueID: e.ExploitUniqueID}).First(&os).Error)
			e.OffensiveSecurity = os

		case models.GitHubRepositoryType:
			gh := &models.GitHubRepository{}
			errs = errs.Add(r.conn.Where(&models.GitHubRepository{ExploitUniqueID: e.ExploitUniqueID}).First(&gh).Error)
			e.GitHubRepository = gh
		}
	}
	for _, e := range errs.GetErrors() {
		if !gorm.IsRecordNotFoundError(e) {
			log15.Error("Failed to get exploit by CveID", "err", e)
		}
	}
	return es
}

// GetExploitMultiByCveID :
func (r *RDBDriver) GetExploitMultiByCveID(cveIDs []string) (exploits map[string][]*models.Exploit) {
	exploits = map[string][]*models.Exploit{}
	for _, cveID := range cveIDs {
		exploits[cveID] = r.GetExploitByCveID(cveID)
	}
	return exploits
}
