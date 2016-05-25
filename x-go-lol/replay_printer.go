package xlol

import (
	"fmt"

	"github.com/charlywazzap/go-ansi"
	"github.com/charlywazzap/go-lol"
)

// ReplayPrinter can be used to display Replay info on command line /
// text output
type ReplayPrinter struct {
	region *lol.Region
	api    *lol.StaticAPIEndpoint
}

// NewReplayPrinter creates a new replay printer for a region
func NewReplayPrinter(region *lol.Region, key lol.APIKey) (*ReplayPrinter, error) {
	res := &ReplayPrinter{
		region: region,
	}
	var err error
	res.api, err = lol.NewStaticAPIEndpoint(region, key)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type playerSelection struct {
	name     string
	champion string
}

// Display nicely display on stdout the replay information
func (p *ReplayPrinter) Display(r *Replay) {
	ansi.ResetColor()
	ansi.Printf("  GameID:%s/%d at %s -- ", r.MetaData.GameKey.PlatformID,
		r.MetaData.GameKey.ID, r.MetaData.StartTime)
	if len(r.GameInfo.Participants) == 0 || r.MetaData.GameKey.PlatformID != p.region.PlatformID() {
		ansi.SetForeground(ansi.Red)
		ansi.Printf("  No information available\n")
		ansi.ResetColor()
		return
	}

	ansi.SetForeground(ansi.Yellow)
	ansi.Printf(" map : %d\n", r.GameInfo.Map)
	maxNameLength := 0
	blueTeam := make([]playerSelection, 0, len(r.GameInfo.Participants))
	redTeam := make([]playerSelection, 0, len(r.GameInfo.Participants))

	for _, part := range r.GameInfo.Participants {
		res := playerSelection{
			name: part.Name,
		}

		champ, err := p.api.GetChampion(part.Champion)
		if err != nil {
			res.champion = fmt.Sprintf("Unknown ChampionID:%d", part.Champion)
		} else {
			res.champion = champ.Name
		}

		if part.TeamID == 100 {
			blueTeam = append(blueTeam, res)
		} else if part.TeamID == 200 {
			redTeam = append(redTeam, res)
		}
		if len(res.name) > maxNameLength {
			maxNameLength = len(res.name)
		}
		if len(res.champion) > maxNameLength {
			maxNameLength = len(res.champion)
		}
	}

	colFormat := fmt.Sprintf(" %%%ds |", maxNameLength)
	//prints the blue team first
	ansi.Printf("  ")
	ansi.SetForeAndBackground(ansi.White, ansi.Blue)
	ansi.Printf("|")
	for _, p := range blueTeam {
		ansi.Printf(colFormat, p.name)
	}
	ansi.ResetColor()
	ansi.Printf("\n  ")
	ansi.SetForeAndBackground(ansi.Blue, ansi.Default)
	ansi.Printf("|")
	for _, p := range blueTeam {
		ansi.Printf(colFormat, p.champion)
	}
	ansi.ResetColor()
	ansi.Printf("\n  ")
	ansi.SetForeAndBackground(ansi.White, ansi.Magenta)
	ansi.Printf("|")
	for _, p := range redTeam {
		ansi.Printf(colFormat, p.name)
	}
	ansi.ResetColor()
	ansi.Printf("\n  ")
	ansi.SetForeAndBackground(ansi.Magenta, ansi.Default)
	ansi.Printf("|")
	for _, p := range redTeam {
		ansi.Printf(colFormat, p.champion)
	}
	ansi.ResetColor()
	ansi.Printf("\n")
}
