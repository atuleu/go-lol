package lol

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path"

	"launchpad.net/go-xdg"

	// imports ability to decode GIG
	_ "image/gif"
	// imports ability to decode JPEG
	_ "image/jpeg"
	// imports ability to decode PNG
	_ "image/png"
)

// An ImageDto describe how to get an image from a Realm
type ImageDto struct {
	Full   string `json:"full"`
	Group  string `json:"group"`
	Sprite string `json:"sprite"`
	H      int    `json:"h"`
	W      int    `json:"w"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

// A Realm represents dataset of data for a game.
type Realm struct {
	// Base URL to get data for the realm
	Cdn        string `json:"cdn"`
	CSSVersion string `json:"css"`
	// Current versio for Data Dragon
	DataDragonVersion string `json:"dd"`
	// Default language for thsi realm
	Locale string `json:"l"`

	LegacyIE6       string            `json:"lg"`
	DataSetsVersion map[string]string `json:"n"`
	Version         string            `json:"v"`
	ProfileIconMax  int               `json:"profileiconmax"`
	Store           string            `json:"store"`

	cachedir string
}

// GetRealm fetches the current Realm used for the selected region of
// the StaticAPIEndpoint. It also initializes the cache structure for
// the Realm data (images)
func (a *StaticAPIEndpoint) GetRealm() (Realm, error) {
	res := Realm{}
	URL := a.formatURL("/realm", nil)
	err := a.get(URL, nil, res)
	if err != nil {
		return Realm{}, err
	}
	res.cachedir, err = xdg.Cache.Ensure(path.Join("go-lol", "realm-images", "version"))
	if err != nil {
		return Realm{}, err
	}
	return res, nil
}

type nullCloser struct {
	buffer bytes.Buffer
}

func (r nullCloser) Read(p []byte) (int, error) {
	return r.buffer.Read(p)
}

func (r nullCloser) Close() error {
	return nil
}

func (r *Realm) openImage(group, name string) (io.ReadCloser, error) {
	version, ok := r.DataSetsVersion[group]
	if ok == false {
		version = r.Version
	}
	cachePath := path.Join(r.cachedir, version, group, name)
	_, err := os.Stat(cachePath)
	if err == nil {
		// image is in path, we open it from here
		return os.Open(cachePath)
	}
	if os.IsNotExist(err) == false {
		return nil, fmt.Errorf("Could not check existance of image in cache (path:%s): %s", cachePath, err)
	}

	//we should get it from http.
	URL := fmt.Sprintf("%s/%s/img/%s/%s", r.Cdn, version, group, name)
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, RESTError{Code: resp.StatusCode}
	}

	// we create the cache file
	err = os.MkdirAll(path.Dir(cachePath), 0755)
	if err != nil {
		return nil, fmt.Errorf("Could not create cache file %s: %s", cachePath, err)
	}

	f, err := os.Create(cachePath)
	if err != nil {
		return nil, fmt.Errorf("Could not create cache file %s: %s", cachePath, err)
	}
	defer f.Close()

	res := nullCloser{}

	_, err = io.Copy(&(res.buffer), io.TeeReader(resp.Body, f))
	if err != nil {
		// we should remove the cache, as it is likely to be
		// incomplete
		f.Close()
		os.RemoveAll(cachePath)
		return nil, err
	}

	return res, nil
}

// GetImage gets the image from the name/full group
func (r *Realm) GetImage(d ImageDto) (image.Image, error) {
	if len(d.Group) == 0 && len(d.Full) == 0 {
		return nil, fmt.Errorf("Invalid ImageDto %+v: missing Full or Group entry", d)
	}

	reader, err := r.openImage(d.Group, d.Full)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	res, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

// GetImageFromSprite gets the image from cutting a sprite
func (r *Realm) GetImageFromSprite(d ImageDto) (image.Image, error) {
	if len(d.Sprite) == 0 || d.H == 0 || d.W == 0 {
		return nil, fmt.Errorf("Invalid ImageDto %+v: missing Sprite, Height or Width", d)
	}

	reader, err := r.openImage("sprite", d.Sprite)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	res, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	casted, ok := res.(subImager)
	if ok == false {
		return nil, fmt.Errorf("Image is decoded, but we are not able to cut it to %d:%d:%d:%d", d.X, d.Y, d.W, d.H)
	}

	return casted.SubImage(image.Rect(d.X, d.Y, d.X+d.W, d.Y+d.H)), nil
}
