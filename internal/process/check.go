package process

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spoonboy-io/switch/internal/extract"

	"github.com/spoonboy-io/koan"

	"github.com/spoonboy-io/switch/internal"
)

var documentTTL map[string]time.Time
var mtx sync.Mutex

// CheckAndRefresh will check what needs to be updated, request the data, process update the cache, plus save
func CheckAndRefresh(ctx context.Context, config internal.Sources, logger *koan.Logger, debug bool) {
	// make our cache
	if documentTTL == nil {
		documentTTL = make(map[string]time.Time)
	}

	// queue sources, not in cache or expired
	var queue internal.Sources

	// iterate over config and check cache time
	for _, cfg := range config {
		// check in map, if it is check expiry
		hash := fmt.Sprintf("%s", md5.Sum([]byte(cfg.Description)))
		if ttl, ok := documentTTL[hash]; ok {
			if time.Now().After(ttl) {
				// cached but expired load to queue
				queue = append(queue, cfg)
			}
		} else {
			// not in cache, get
			queue = append(queue, cfg)
		}
	}

	// cancel if not done in 10secs
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var wg sync.WaitGroup

	for i, _ := range queue {

		wg.Add(1)
		go func(q internal.Source) {
			// make request
			logger.Info(fmt.Sprintf("Requesting `%s` (%s)", q.Description, q.Url))

			req, err := http.NewRequest("GET", q.Url, nil)
			if err != nil {
				logger.Error("bad request", err)
			}

			req = req.WithContext(ctx)
			req.Header.Add("Content-Type", "application/json")

			// form the authorization header if exists
			if q.Token != "" {
				req.Header.Add("Authorization", q.Token)
			}

			// allow insecure certs
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}

			res, err := client.Do(req)
			if err != nil {
				logger.Error("bad request", err)
				wg.Done()
			} else {
				defer res.Body.Close()

				if res.StatusCode != http.StatusOK {
					logger.Error("bad response", fmt.Errorf("(%d): Source: %s, URL: %s", res.StatusCode, q.Description, q.Url))

				}

				// update the docTTL
				hash := fmt.Sprintf("%s", md5.Sum([]byte(q.Description)))
				ttl := time.Now().Add(time.Duration(q.Ttl) * time.Minute)

				mtx.Lock()
				documentTTL[hash] = ttl
				mtx.Unlock()

				// process response
				input, err := io.ReadAll(res.Body)
				if err != nil {
					logger.Error("can't read response", err)
				}

				if debug {
					dout := fmt.Sprintf("Debug: received from `%s`\n\n%s", q.Url, input)
					logger.Info(dout)
				}

				output, err := extract.ParseJSONForKeyValue(q.Extract.Name, q.Extract.Value, input, q.Extract.Root)
				if err != nil {
					logger.Error("problem parsing json", err)
				}

				if debug {
					dout := fmt.Sprintf("Debug: saving parsed JSON for `%s`\n\n%s", q.Description, output)
					logger.Info(dout)
				}

				// save data
				target := fmt.Sprintf("%s/%s", q.Save.Folder, q.Save.Filename)
				if err := os.WriteFile(target, output, 0644); err != nil {
					logger.FatalError("problem writing file", err)
				}

				wg.Done()
			}

		}(queue[i].Source)

	}

	wg.Wait()
}
