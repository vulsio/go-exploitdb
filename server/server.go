package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/inconshreveable/log15"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mozqnet/go-exploitdb/db"
	"github.com/spf13/viper"
)

// Start :
func Start(logDir string, driver db.DB) error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	logPath := filepath.Join(logDir, "access.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if _, err := os.Create(logPath); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: f,
	}))

	// Routes
	e.GET("/health", health())
	e.GET("/cves/:cve", getExploitByCveID(driver))
	e.GET("/id/:id", getExploitByID(driver))

	bindURL := fmt.Sprintf("%s:%s", viper.GetString("bind"), viper.GetString("port"))
	log15.Info("Listening...", "URL", bindURL)

	e.Logger.Fatal(e.Start(bindURL))
	return nil
}

func health() echo.HandlerFunc {
	return func(context echo.Context) error {
		return context.String(http.StatusOK, "")
	}
}

func getExploitByCveID(driver db.DB) echo.HandlerFunc {
	return func(context echo.Context) (err error) {
		cve := context.Param("cve")
		log15.Debug("Params", "CVE", cve)

		exploits := driver.GetExploitByCveID(cve)
		if err != nil {
			log15.Error("Failed to get Exploit Info by CVE.", "err", err)
		}
		return context.JSON(http.StatusOK, exploits)
	}
}

func getExploitByID(driver db.DB) echo.HandlerFunc {
	return func(context echo.Context) (err error) {
		exploitDBID := context.Param("id")
		log15.Debug("Params", "ExploitDBID", exploitDBID)

		exploit := driver.GetExploitByID(exploitDBID)
		if err != nil {
			log15.Error("Failed to get Exploit Info by EDB-ID.", "err", err)
		}
		return context.JSON(http.StatusOK, exploit)
	}
}
