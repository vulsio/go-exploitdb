package commands

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/inconshreveable/log15"
	"github.com/mozqnet/go-exploitdb/db"
	"github.com/mozqnet/go-exploitdb/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cveIDRegexp       = regexp.MustCompile(`^CVE-\d{1,}-\d{1,}$`)
	exploitDBIDRegexp = regexp.MustCompile(`^\d{1,}$`)
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the data of exploit",
	Long:  `Search the data of exploit`,
	RunE:  searchExploit,
}

func init() {
	RootCmd.AddCommand(searchCmd)

	searchCmd.PersistentFlags().String("type", "", "All Exploits by CVE: CVE  |  by ID: ID (default: CVE)")
	viper.BindPFlag("type", searchCmd.PersistentFlags().Lookup("type"))
	viper.SetDefault("type", "CVE")

	searchCmd.PersistentFlags().String("param", "", "All Exploits: None  |  by CVE: [CVE-xxxx]  | by ID: [xxxx]  (default: None)")
	viper.BindPFlag("param", searchCmd.PersistentFlags().Lookup("param"))
	viper.SetDefault("param", "")
}

func searchExploit(cmd *cobra.Command, args []string) (err error) {
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

	searchType := viper.GetString("type")
	param := viper.GetString("param")

	var results = []*models.Exploit{}
	switch searchType {
	case "CVE":
		if !cveIDRegexp.Match([]byte(param)) {
			log15.Error("Specify the search type [CVE] parameters like `--param CVE-xxxx-xxxx`")
			return errors.New("Invalid CVE Param")
		}
		results = driver.GetExploitByCveID(param)
	case "ID":
		if !exploitDBIDRegexp.MatchString(param) {
			log15.Error("Specify the search type [ID] parameters like `--param 10000`")
			return errors.New("Invalid ID Param")
		}
		results = driver.GetExploitByID(param)
	default:
		log15.Error("Specify the search type [ CVE / ID].")
		return errors.New("Invalid Type")
	}
	fmt.Println("")
	fmt.Println("Results: ")
	fmt.Println("---------------------------------------")
	if len(results) == 0 {
		fmt.Println("No Record Found")
	}
	for _, r := range results {
		fmt.Println("\n[*]CVE-ExploitID Reference:")
		fmt.Printf("  CVE: %s\n", r.CveID)
		fmt.Printf("  Exploit Type: %s\n", r.ExploitType)
		fmt.Printf("  Exploit Unique ID: %s\n", r.ExploitUniqueID)
		fmt.Printf("  URL: %s\n", r.URL)
		fmt.Printf("  Description: %s\n", r.Description)
		fmt.Printf("\n[*]Exploit Detail Info: ")
		if r.OffensiveSecurity != nil {
			fmt.Printf("\n  [*]OffensiveSecurity: ")
			os := r.OffensiveSecurity
			if os.Document != nil {
				fmt.Println("\n  - Document:")
				fmt.Printf("    Path: %s\n", os.Document.DocumentURL)
				fmt.Printf("    File Type: %s\n", os.Document.Type)
			}
			if os.ShellCode != nil {
				fmt.Println("\n  - Exploit Code or Proof of Concept:")
				fmt.Printf("    %s\n", os.ShellCode.ShellCodeURL)
			}
		}
		fmt.Println("---------------------------------------")
	}
	return nil
}
