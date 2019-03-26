package commands

import (
	"github.com/inconshreveable/log15"
	"github.com/mozqnet/go-exploitdb/db"
	"github.com/mozqnet/go-exploitdb/fetcher"
	"github.com/mozqnet/go-exploitdb/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fetchGitHubReposCmd = &cobra.Command{
	Use:   "githubrepos",
	Short: "Fetch the data of github repos",
	Long:  `Fetch the data of github repos`,
	RunE:  fetchGitHubRepos,
}

func init() {
	fetchCmd.AddCommand(fetchGitHubReposCmd)
}

func fetchGitHubRepos(cmd *cobra.Command, args []string) (err error) {
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

	log15.Info("Fetching GitHub Repos Exploit")
	var exploits []*models.Exploit
	if exploits, err = fetcher.FetchGitHubRepos(viper.GetBool("deep")); err != nil {
		log15.Error("Failed to fetch GitHubRepo Exploit", "err", err)
		return err
	}
	log15.Info("GitHub Repos Exploit", "count", len(exploits))

	log15.Info("Insert Exploit into go-exploitdb.", "db", driver.Name())
	if err := driver.InsertExploit(exploits); err != nil {
		log15.Error("Failed to insert.", "dbpath", viper.GetString("dbpath"), "err", err)
		return err
	}
	return nil
}
