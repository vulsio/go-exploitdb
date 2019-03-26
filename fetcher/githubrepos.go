package fetcher

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/mozqnet/go-exploitdb/extractor"
	"github.com/mozqnet/go-exploitdb/models"
	"github.com/mozqnet/go-exploitdb/util"
)

// FetchGitHubRepos :
func FetchGitHubRepos(deep bool) (exploits []*models.Exploit, err error) {
	for year := 1999; year <= time.Now().Year(); year++ {
		log15.Info("Fetching GitHub Repository", "year", year)
		page := 1
		maxPage := 10
		for {
			if maxPage < page {
				break
			}
			spaceCode := "%20"
			url := fmt.Sprintf("https://api.github.com/search/repositories?q=CVE%s%d+in:name&&page=%d&per_page=100", spaceCode, year, page)
			githubJSON, err := util.FetchURL(url)
			if err != nil {
				return nil, err
			}
			var githubs models.GitHubRepoJSON
			if err = json.Unmarshal(githubJSON, &githubs); err != nil {
				return nil, err
			}
			// for github rate limit
			if 1000 < len(githubs.Items) {
				log15.Warn("More than 1000 data can not be acquired due to rate limit of github")
			}
			for _, github := range githubs.Items {
				cveIDs := extractor.ExtractCveID([]byte(github.FullName))
				if len(cveIDs) == 0 {
					continue
				}
				for _, cveID := range cveIDs {
					githubRepoExploit := &models.Exploit{
						ExploitUniqueID: fmt.Sprintf("%s-%s", models.GitHubRepositoryType, github.URL),
						ExploitType:     models.GitHubRepositoryType,
						URL:             github.URL,
						CveID:           cveID,
						Description:     github.Description,
						GitHubRepository: &models.GitHubRepository{
							ExploitUniqueID: github.URL,
							Star:            github.Star,
							Fork:            github.Fork,
						},
					}
					exploits = append(exploits, githubRepoExploit)
				}
			}
			if page == 1 {
				totalPageSize := (githubs.TotalCount / 100) + 1
				if totalPageSize < maxPage {
					maxPage = totalPageSize
				}
			}
			page++
		}
	}
	return exploits, nil
}
