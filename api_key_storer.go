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

// An APIKey is used to authenticate access to Riot Game API
type APIKey string

// Check returns an error if the APIKey is not valid, nil otherwise
func (k APIKey) Check() error {
	if len(k) == 0 {
		return fmt.Errorf("Key is empty")
	}

	if apiKeyRx.MatchString(string(k)) == false {
		return fmt.Errorf("Invalid API key syntax `%s'", k)
	}

	return nil
}

// An APIKeyStorer is used to safely Store and Get an APIKey for the
// user
type APIKeyStorer interface {
	Get() (APIKey, bool)
	Store(k APIKey) error
}

// XdgAPIKeyStorer Store and Get an APIkey from a file located in the
// $XDG_CONFIG_HOME directory
type XdgAPIKeyStorer struct {
	path string
	lock lockfile.Lockfile
	key  APIKey
}

// NewXdgAPIKeyStorer creates a new XdgAPIKeyStorer
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

	str := strings.TrimSpace(string(data))
	if len(str) == 0 {
		return nil
	}

	s.key = APIKey(str)
	return s.key.Check()
}

var apiKeyRx = regexp.MustCompile(`\A[0-9a-f]{8}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{12}\z`)

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

// Get retrieve the APIKey from user's XDG_CONFIG_HOME
func (s *XdgAPIKeyStorer) Get() (APIKey, bool) {
	return s.key, len(s.key) > 0
}

// Store saves the new APIKey k to user's XDG_CONFIG_HOME
func (s *XdgAPIKeyStorer) Store(k APIKey) error {
	if err := k.Check(); err != nil {
		return err
	}
	oldKey := s.key
	s.key = k
	err := s.save()
	if err != nil {
		s.key = oldKey
	}
	return err
}
