package fetcher

import (
	"fmt"

	"github.com/mozqnet/go-exploitdb/extractor"
	"github.com/mozqnet/go-exploitdb/models"
	"github.com/mozqnet/go-exploitdb/util"
	"github.com/russross/blackfriday/v2"
)

// FetchAwesomePoc :
func FetchAwesomePoc(deep bool) (exploits []*models.Exploit, err error) {
	url := "https://raw.githubusercontent.com/qazbnm456/awesome-cve-poc/master/README.md"
	readme, err := util.FetchURL(url)
	if err != nil {
		return nil, err
	}
	r := &extractor.AwesomePocReader{}
	blackfriday.Run(readme, blackfriday.WithRenderer(r))

	for poc := range r.AwesomePoc {
		exploit := &models.Exploit{
			ExploitType:     models.AwesomePocType,
			ExploitUniqueID: fmt.Sprintf("%s-%s", models.AwesomePocType, poc.URL),
			URL:             poc.URL,
			Description:     poc.Description,
			CveID:           poc.CveID,
		}
		exploits = append(exploits, exploit)
	}
	return exploits, nil
}
