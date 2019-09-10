package commands

import (
	"github.com/inconshreveable/log15"
	"github.com/mozqnet/go-exploitdb/db"
	server "github.com/mozqnet/go-exploitdb/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start go-exploitdb HTTP server",
	Long:  `Start go-exploitdb HTTP server`,
	RunE:  executeServer,
}

func init() {
	RootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().String("bind", "", "HTTP server bind to IP address (default: loop back interface")
	viper.BindPFlag("bind", serverCmd.PersistentFlags().Lookup("bind"))
	viper.SetDefault("bind", "127.0.0.1")

	serverCmd.PersistentFlags().String("port", "", "HTTP server port number (default: 1326")
	viper.BindPFlag("port", serverCmd.PersistentFlags().Lookup("port"))
	viper.SetDefault("port", "1326")

}

func executeServer(cmd *cobra.Command, args []string) (err error) {
	logDir := viper.GetString("log-dir")
	driver, locked, err := db.NewDB(viper.GetString("dbtype"), viper.GetString("dbpath"), viper.GetBool("debug-sql"))
	if err != nil {
		if locked {
			log15.Error("Failed to initialize DB. Close DB connection before fetching", "err", err)
		}
		return err
	}

	log15.Info("Starting HTTP Server...")
	if err = server.Start(logDir, driver); err != nil {
		log15.Error("Failed to start server.", "err", err)
		return err
	}

	return nil
}
