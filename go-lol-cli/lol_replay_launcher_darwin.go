package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"time"

	"github.com/charlywazzap/go-lol"
)

type darwinLauncher struct {
	launcherPath string
	clientPath   string
}

const (
	launcherReleasesBasepath = "Contents/LoL/RADS/solutions/lol_game_client_sln/releases"
	launcherPath             = "deploy/LeagueOfLegends.app/Contents/MacOS/LeagueofLegends"
	clientReleasesBasepath   = "Contents/LoL/RADS/projects/lol_air_client/releases"
	clientPath               = "deploy/bin/LolClient"
)

func (l *darwinLauncher) lookForExecutableReleases(basepath, binPath string) (string, error) {
	filesInfo, err := ioutil.ReadDir(basepath)
	if err != nil {
		return "", fmt.Errorf("Could not check for releases: %s", err)
	}

	versions := make([]version, 0, len(filesInfo))
	for _, inf := range filesInfo {
		if inf.IsDir() == false {
			continue
		}
		v, err := newVersion(inf.Name())
		if err != nil {
			return "", err
		}
		versions = append(versions, v)
	}

	sort.Sort(sort.Reverse(versionList(versions)))
	if len(versions) == 0 {
		return "", fmt.Errorf("Could not find a release in %s", basepath)
	}

	res := path.Join(basepath, versions[0].String(), binPath)

	info, err := os.Stat(res)
	if err != nil {
		return "", fmt.Errorf("Could not find binary %s: %s", res, err)
	}
	if info.Mode().Perm()&111 == 111 {
		return "", fmt.Errorf("Wrong persion %v on binary %s", info.Mode().Perm(), res)
	}

	return res, nil
}

func NewLolReplayLauncher(basepath string) (LolReplayLauncher, error) {
	if len(basepath) == 0 {
		basepath = path.Join("/Applications", "League Of Legends.app")
	}

	info, err := os.Stat(basepath)
	if err != nil {
		return nil, fmt.Errorf("Could not check for League Of Legend application in %s:  %s", basepath, err)
	}

	if info.IsDir() == false {
		return nil, fmt.Errorf("%s is not a directory/App Bundle, please check your LoL installation", basepath)
	}

	res := &darwinLauncher{}
	res.launcherPath, err = res.lookForExecutableReleases(path.Join(basepath, launcherReleasesBasepath), launcherPath)
	if err != nil {
		return nil, fmt.Errorf("Could not find LoL launcher: %s", err)
	}

	res.clientPath, err = res.lookForExecutableReleases(path.Join(basepath, clientReleasesBasepath), clientPath)
	if err != nil {
		return nil, fmt.Errorf("Could not find LoL Client: %s", err)
	}

	return res, nil
}

func (l *darwinLauncher) Launch(address string, region *lol.Region, id lol.GameID, encryptionKey string) error {
	cmd := exec.Command(l.launcherPath, MaestroParam1, MaestroParam2, l.clientPath,
		fmt.Sprintf(`spectator %s %s %d %s`, address, encryptionKey, id, region.PlatformID()))

	cmd.Dir = path.Dir(l.launcherPath)
	cmd.Env = append(os.Environ(), "riot_launched=true")

	now := time.Now()
	debugErr, err := ioutil.TempFile("", fmt.Sprintf("go-lol-client-debug-stderr-%d.txt", now.Unix()))
	if err != nil {
		return err
	}
	defer debugErr.Close()

	debugOut, err := ioutil.TempFile("", fmt.Sprintf("go-lol-client-debug-stdout-%d.txt", now.Unix()))
	if err != nil {
		return err
	}
	defer debugOut.Close()

	log.Printf("Connecting client output to %s and %s", debugErr.Name(), debugOut.Name())
	cmd.Stdout = debugOut
	cmd.Stderr = debugErr

	return cmd.Run()
}
