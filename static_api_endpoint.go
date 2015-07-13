package lol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"

	"launchpad.net/go-xdg"
)

// StaticAPIEndpoint is used to access Riot REST API for static data.
type StaticAPIEndpoint struct {
	region       *Region
	key          APIKey
	staticRegion *Region
	cachedir     string
	version      string
}

// NewStaticAPIEndpoint is creeating a new Static endpoint. You have
// to pass the Dynamic (i.e. EUW, KR NA) region you are interested in
// fecthing data.
func NewStaticAPIEndpoint(region *Region, key APIKey) (*StaticAPIEndpoint, error) {
	if region.IsDynamic() == false {
		return nil, fmt.Errorf("We need a duynamic region for looking up data")
	}
	res := &StaticAPIEndpoint{
		region: region,
		key:    key,
	}
	res.staticRegion, _ = NewRegion(GLOBAL)

	cacheVersion, err := xdg.Cache.Ensure(path.Join("go-lol", "static-data-cache", "version"))
	if err != nil {
		return nil, err
	}
	//we should get the current version
	versions := make([]string, 0, 10)
	resp, err := http.Get(res.formatURL("/versions", nil))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Could not get current data versions: got error code: %s ", resp.Status)
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&versions)
	if err != nil {
		return nil, fmt.Errorf("Could not parse list of versions: %s", err)
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("Invalid empty list of versions")
	}
	res.version = versions[0]

	// we should create the cache dir
	res.cachedir = path.Join(path.Dir(cacheVersion), res.version)
	err = os.MkdirAll(res.cachedir, 0755)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize cache directory %s: %s", res.cachedir, err)
	}

	return res, nil
}

func (a *StaticAPIEndpoint) formatURL(url string, options map[string]string) string {
	res := fmt.Sprintf("https://%s/api/lol/static-data/%s/v1.2%s?api_key=%s",
		a.staticRegion.url, a.region.code, url, a.key)
	for k, v := range options {
		res = res + "&" + k + "=" + v
	}
	return res
}

func (a *StaticAPIEndpoint) formatCacheFile(url string, options map[string]string) string {
	res := path.Join(a.cachedir, url)
	//beware map are unsorted, we should sort them !
	keys := make([]string, 0, len(options))

	for k := range options {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		res = res + "_" + k + "_" + options[k]
	}

	return res
}

// get data from that endpoint
func (a *StaticAPIEndpoint) cachedGet(url string, options map[string]string, v interface{}) error {
	filepath := a.formatCacheFile(url, options)
	_, err := os.Stat(filepath)
	var reader io.Reader
	var cleanupError error
	if err == nil {
		f, err := os.Open(filepath)
		if err != nil {
			return fmt.Errorf("Could not open cache file %s: %s", filepath, err)
		}
		defer f.Close()
		reader = f
	} else {
		if os.IsNotExist(err) == false {
			return fmt.Errorf("Could not check for cache file existence %s: %s", filepath, err)
		}
		//create cache file
		err = os.MkdirAll(path.Dir(filepath), 0755)
		if err != nil {
			return fmt.Errorf("Could not create cache file %s: %s", filepath, err)
		}
		f, err := os.Create(filepath)
		if err != nil {
			return fmt.Errorf("Could not create cache file %s: %s", filepath, err)
		}
		defer func() {
			f.Close()
			if cleanupError != nil {
				os.RemoveAll(filepath)
			}
		}()

		fullURL := a.formatURL(url, options)
		resp, cleanupError := http.Get(fullURL)
		if err != nil {
			return fmt.Errorf("Could not reach %s: %s", fullURL, err)
		}
		if resp.StatusCode >= 400 {
			resp.Body.Close()
			cleanupError = RESTError{Code: resp.StatusCode}
			return cleanupError
		}
		defer resp.Body.Close()

		var buffer bytes.Buffer

		_, cleanupError = io.Copy(&buffer, io.TeeReader(resp.Body, f))
		if cleanupError != nil {
			return fmt.Errorf("Could not cache data from %s: %s", fullURL, err)
		}
	}

	dec := json.NewDecoder(reader)
	cleanupError = dec.Decode(v)
	return cleanupError
}
