package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type dataWrapper struct {
	Data Data `json:"data"`
}

type Data struct {
	Games GamesWrapper `json:"games"`
}
type GamesWrapper struct {
	Games []Game `json:"game"`
}

type Game struct {
	//TEAMCITY = data: games: game: "home_team_city" / "away_team_city"
	HomeCity string `json:"home_team_city"`
	//TEAMNAME = data: games: game: "home_team_name" / "away_team_name"
	HomeTeamName string `json:"home_team_name"`
	//RUNSCORE = data: games: game: linescore: r: "home" / "away"
	HomeWins   string `json:"home_win"`
	HomeLosses string `json:"home_loss"`
	//TEAMCITY = data: games: game: "home_team_city" / "away_team_city"
	AwayCity string `json:"away_team_city"`
	//TEAMNAME = data: games: game: "home_team_name" / "away_team_name"
	AwayTeamName string `json:"away_team_name"`
	//RUNSCORE = data: games: game: linescore: r: "home" / "away"
	AwayWins   string    `json:"away_win"`
	AwayLosses string    `json:"away_loss"`
	Linescore  Linescore `json:"linescore"`
}

type Linescore struct {
	Runs Runs `json:"r"`
}

type Runs struct {
	HomeRuns string `json:"home"`
	AwayRuns string `json:"away"`
}

type HomeTeam struct {
	//TEAMCITY = data: games: game: "home_team_city" / "away_team_city"
	City string `json:"home_team_city"`
	//TEAMNAME = data: games: game: "home_team_name" / "away_team_name"
	TeamName string `json:"home_team_name"`
	//RUNSCORE = data: games: game: linescore: r: "home" / "away"
	RunScore string `json:"linescore.r.home"`
	Wins     string `json:"home_win"`
	Losses   string `json:"home_loss"`
}

type AwayTeam struct {
	//TEAMCITY = data: games: game: "home_team_city" / "away_team_city"
	City string `json:"away_team_city"`
	//TEAMNAME = data: games: game: "home_team_name" / "away_team_name"
	TeamName string `json:"away_team_name"`
	//RUNSCORE = data: games: game: linescore: r: "home" / "away"
	RunScore string `json:"linescore.r.away"`
	Wins     string `json:"away_win"`
	Losses   string `json:"away_loss"`
}

func main() {
	t := time.Now().AddDate(0, 0, -1)
	day := t.Day()
	month := t.Format("01")
	year := t.Year()
	url := fmt.Sprintf("https://gd2.mlb.com/components/game/mlb/year_%v/month_%v/day_%v/master_scoreboard.json", year, month, day)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("error getting url status: %v", resp.StatusCode))
	}
	defer resp.Body.Close()
	d := dataWrapper{}
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		panic(fmt.Sprintf("error decoding json: %v", err))
	}

	fmt.Printf("%+v", d)
}
