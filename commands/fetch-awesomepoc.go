package commands

import (
	"github.com/inconshreveable/log15"
	"github.com/mozqnet/go-exploitdb/db"
	"github.com/mozqnet/go-exploitdb/fetcher"
	"github.com/mozqnet/go-exploitdb/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fetchAwesomePocCmd = &cobra.Command{
	Use:   "awesomepoc",
	Short: "Fetch the data of Awesome Poc",
	Long:  `Fetch the data of Awesome Poc`,
	RunE:  fetchAwesomePoc,
}

func init() {
	fetchCmd.AddCommand(fetchAwesomePocCmd)
}

func fetchAwesomePoc(cmd *cobra.Command, args []string) (err error) {
	driver, locked, err := db.NewDB(
		viper.GetString("dbtype"),
		viper.GetString("dbpath"),
		viper.GetBool("debug-sql"),
	)
	if err != nil {
		if locked {
			log15.Error("Failed to initialize DB. Close DB connection before fetching", "err", err)
		}
		return err
	}

	log15.Info("Fetching Awesome Poc Exploit")
	var exploits []*models.Exploit
	if exploits, err = fetcher.FetchAwesomePoc(viper.GetBool("deep")); err != nil {
		log15.Error("Failed to fetch AwesomePoc Exploit", "err", err)
		return err
	}
	log15.Info("Awesome Poc Exploit", "count", len(exploits))

	log15.Info("Insert Exploit into go-exploitdb.", "db", driver.Name())
	if err := driver.InsertExploit(exploits); err != nil {
		log15.Error("Failed to insert.", "dbpath", viper.GetString("dbpath"), "err", err)
		return err
	}
	return nil
}
