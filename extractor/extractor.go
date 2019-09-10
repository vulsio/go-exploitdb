package extractor

import (
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"

	"github.com/russross/blackfriday/v2"
)

var (
	// ex:
	// CVE-2018-1111
	// CVE: 2018-1111
	// CVE: CVE-2018–14064
	// CVE: cve-2018-8002
	// CVE: CVE 2018-7719
	// CVE : CVE- 2018-7198
	// CVE : [CVE-2018- 11034]
	// CVE: 2017-13056
	// CVE; 2016-0953
	// CVE : 2012-3755
	// CVE 2012-0550
	cveIDRegexp1 = regexp.MustCompile(`(?i)CVE\s?[-–:;]?\s?(\d{4})[-–]\s?(\d{4,})`)
	// ex: CVE Number:   2011-4189
	cveIDRegexp2 = regexp.MustCompile(`CVE Number\s*[:;]\s*(\d{4})[-–](\d{4,})`)
	// ex:
	// [ 'CVE', '2009-0184' ]
	// ['CVE',    '2005-2265']
	// [ 'CVE', '2003-0352'  ]
	// [ 'CVE', '2011-2404 ']
	cveIDRegexp3 = regexp.MustCompile(`\[\s*'CVE'\s*,\s*'(\d{4})[-–](\d{4,})\s*'\s*\]`)
	// ex: ['CVE'     => '2008-6825']
	cveIDRegexp4 = regexp.MustCompile(`\[\s*'CVE'\s*=>\s*'(\d{4})[-–](\d{4,})\s*'\s*\]`)
	// ex : CVE : [2014-3443]
	cveIDRegexp5 = regexp.MustCompile(`CVE\s*:\s*\[(\d{4})[-–](\d{4,})\]`)
	// ex: "cve20113872",
	// ex: "cve_2011_3556",
	// ex: "CVE20120053",
	// ex: "CVE_2012_4681",
	cveIDRegexp6 = regexp.MustCompile(`(?i)CVE[-_]?(\d{4})[-_]?(\d{4,})`)
	// ex: CVE-2018-1002105
	cveIDRegexpSimple = regexp.MustCompile(`^CVE-\d{4}-\d{4,}$`)
)

// ExtractCveID :
func ExtractCveID(file []byte) (cveIDs []string) {
	uniqCveID := map[string]struct{}{}
	regxps := []*regexp.Regexp{
		cveIDRegexp1,
		cveIDRegexp2,
		cveIDRegexp3,
		cveIDRegexp4,
		cveIDRegexp5,
		cveIDRegexp6,
	}

	for _, re := range regxps {
		results := re.FindAllSubmatch(file, -1)
		for _, matches := range results {
			if 2 < len(matches) {
				cveID := fmt.Sprintf("CVE-%s-%s", matches[1], matches[2])
				uniqCveID[cveID] = struct{}{}
			}
		}
	}

	for cveID := range uniqCveID {
		cveIDs = append(cveIDs, cveID)
	}
	return cveIDs
}

// AwesomePoc :
type AwesomePoc struct {
	CveID       string
	Description string
	URL         string
}

// AwesomePocReader :
type AwesomePocReader struct {
	Start bool
	// unique
	AwesomePoc        map[AwesomePoc]struct{}
	FillingAwesomePoc AwesomePoc
}

// RenderNode :
func (r *AwesomePocReader) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	// start extracting from the text "Resource"
	if string(node.Literal) == "Resource" {
		r.Start = true
		r.AwesomePoc = map[AwesomePoc]struct{}{}
		return 0
	}
	if !r.Start {
		return 0
	}

	emptyAwesomePoc := AwesomePoc{}
	switch node.Type {
	case blackfriday.Text:
		if cveIDRegexpSimple.Match(node.Literal) {
			// change awesome poc
			if r.FillingAwesomePoc != emptyAwesomePoc {
				r.AwesomePoc[r.FillingAwesomePoc] = struct{}{}
			}
			r.FillingAwesomePoc = AwesomePoc{
				CveID: string(node.Literal),
			}
		} else if 0 < len(r.FillingAwesomePoc.CveID) {
			if 0 < len(node.Literal) {
				r.FillingAwesomePoc.Description = string(node.Literal)
			}
		}
	case blackfriday.Link:
		if 0 < len(r.FillingAwesomePoc.CveID) && len(r.FillingAwesomePoc.URL) == 0 {
			if !strings.Contains(string(node.LinkData.Destination), "http") {
				// passed relative path of awesome-cve-poc
				r.FillingAwesomePoc.URL = path.Join("https://github.com/qazbnm456/awesome-cve-poc/blob/master", string(node.LinkData.Destination))
			} else {
				r.FillingAwesomePoc.URL = string(node.LinkData.Destination)
			}
		}
	}
	return 0
}

// RenderHeader :
func (r *AwesomePocReader) RenderHeader(w io.Writer, ast *blackfriday.Node) {
}

// RenderFooter :
func (r *AwesomePocReader) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	// 最後のFillingを挿入
	emptyAwesomePoc := AwesomePoc{}
	if r.FillingAwesomePoc != emptyAwesomePoc {
		r.AwesomePoc[r.FillingAwesomePoc] = struct{}{}
	}
}
