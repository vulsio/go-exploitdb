package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	"github.com/k0kubun/pp"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
)

// GenWorkers :
func GenWorkers(num int) chan<- func() {
	tasks := make(chan func())
	for i := 0; i < num; i++ {
		go func() {
			for f := range tasks {
				f()
			}
		}()
	}
	return tasks
}

// GetDefaultLogDir :
func GetDefaultLogDir() string {
	defaultLogDir := "/var/log/go-exploitdb"
	if runtime.GOOS == "windows" {
		defaultLogDir = filepath.Join(os.Getenv("APPDATA"), "go-exploitdb")
	}
	return defaultLogDir
}

// SetLogger :
func SetLogger(logDir string, quiet, debug, logJSON bool) {
	stderrHundler := log15.StderrHandler
	logFormat := log15.LogfmtFormat()
	if logJSON {
		logFormat = log15.JsonFormatEx(false, true)
		stderrHundler = log15.StreamHandler(os.Stderr, logFormat)
	}

	lvlHundler := log15.LvlFilterHandler(log15.LvlInfo, stderrHundler)
	if debug {
		lvlHundler = log15.LvlFilterHandler(log15.LvlDebug, stderrHundler)
	}
	if quiet {
		lvlHundler = log15.LvlFilterHandler(log15.LvlDebug, log15.DiscardHandler())
		pp.SetDefaultOutput(ioutil.Discard)
	}

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.Mkdir(logDir, 0700); err != nil {
			log15.Error("Failed to create log directory", "err", err)
		}
	}
	var hundler log15.Handler
	if _, err := os.Stat(logDir); err == nil {
		logPath := filepath.Join(logDir, "go-exploitdb.log")
		hundler = log15.MultiHandler(
			log15.Must.FileHandler(logPath, logFormat),
			lvlHundler,
		)
	} else {
		hundler = lvlHundler
	}
	log15.Root().SetHandler(hundler)
}

// FetchURL returns HTTP response body
func FetchURL(url string, apiKey ...string) ([]byte, error) {
	var errs []error
	httpProxy := viper.GetString("http-proxy")

	resp, body, err := gorequest.New().Proxy(httpProxy).Get(url).Type("text").EndBytes()
	if len(errs) > 0 || resp == nil {
		return nil, fmt.Errorf("HTTP error. errs: %v, url: %s", err, url)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP error. errs: %v, status code: %d, url: %s", err, resp.StatusCode, url)
	}
	// for github rate limit
	if strings.Contains(url, "github") {
		var rateErr error
		if rateErr = WaitForRateLimit(resp.Header); rateErr != nil {
			return nil, rateErr
		}
	}
	return body, nil
}

// WaitForRateLimit :
func WaitForRateLimit(header http.Header) (err error) {
	if header == nil {
		return nil
	}
	if remaindings, ok := header["X-Ratelimit-Remaining"]; ok {
		for _, remainding := range remaindings {
			var r int
			if r, err = strconv.Atoi(remainding); err != nil {
				log15.Error("Wrong http header", "X-Ratelimit-Remaining", remaindings, "err", err)
				return err
			}
			if 1 < r {
				return nil
			}
		}
	}
	if resets, ok := header["X-Ratelimit-Reset"]; ok {
		for _, reset := range resets {
			var r int64
			if r, err = strconv.ParseInt(reset, 10, 64); err != nil {
				log15.Error("Wrong http header", "X-Ratelimit-Reset", reset, "err", err)
				return err
			}
			// add 1s
			duration := time.Until(time.Unix(r, 0))
			log15.Info("Sleep for GitHub rate limit", "duration", duration)
			time.Sleep(duration)
		}
	}
	return nil
}

// DeleteRecordNotFound deletes gorm.ErrRecordNotFound in errs
func DeleteRecordNotFound(errs []error) (new []error) {
	for _, err := range errs {
		if err != nil && err != gorm.ErrRecordNotFound {
			new = append(new, err)
		}
	}
	return new
}

// DeleteNil deletes nil in errs
func DeleteNil(errs []error) (new []error) {
	for _, err := range errs {
		if err != nil {
			new = append(new, err)
		}
	}
	return new
}
