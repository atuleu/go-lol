package lol

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/nightlyone/lockfile"
	"launchpad.net/go-xdg"
)

type APIKey string

type APIKeyStorer interface {
	Get() APIKey
	Store(k APIKey) error
}

type XdgAPIKeyStorer struct {
	path string
	lock lockfile.Lockfile
	key  string
}

func NewXdgAPIKeyStorer() (*XdgAPIKeyStorer, error) {
	res := &XdgAPIKeyStorer{}
	var err error
	res.path, err = xdg.Config.Ensure("go-lol/api-key")
	if err != nil {
		return nil, fmt.Errorf("Could not find configuration file for API key: %s", err)
	}
	//ensure it is 0700 ! this is user critical data
	info, err := os.Stat(res.path)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModePerm != 700 {
		err := os.Chmod(res.path, 0700)
		if err != nil {
			return nil, err
		}
	}

	//Create a lockfile for reading and setting
	res.lock, err = lockfile.New(path.Join(path.Dir(res.path), "apikey.lock"))
	if err != nil {
		return nil, err
	}

	return res, res.load()
}

func (s *XdgAPIKeyStorer) load() error {
	if err := s.lock.TryLock(); err != nil {
		return fmt.Errorf("Could not lock data: %s", err)
	}
	defer s.lock.Unlock()

	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	var data []byte
	data, err = ioutil.ReadAll(f)

	if err != nil {
		return err
	}

	s.key = strings.TrimSpace(string(data))

	// Check key format validity
	if apiKeyRx.MatchString(s.key) == false {
		return fmt.Errorf("Invalid key syntax `%s'", s.key)
	}

	return nil
}

var apiKeyRx = regexp.MustCompile(`\AA[0-9a-f]{8}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{12}\z`)

func (s *XdgAPIKeyStorer) save() error {
	if err := s.lock.TryLock(); err != nil {
		return fmt.Errorf("Could not lock data: %s", err)
	}
	defer s.lock.Unlock()

	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n", s.key)
	return err
}

func (s *XdgAPIKeyStorer) Get() APIKey {
	return APIKey(s.key)
}

func (s *XdgAPIKeyStorer) Store(k APIKey) error {
	key := string(k)
	if apiKeyRx.MatchString(key) == false {
		return fmt.Errorf("Invalid key syntax `%s'", k)
	}
	oldKey := s.key
	s.key = key
	err := s.save()
	if err != nil {
		s.key = oldKey
	}
	return err
}
