package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type dataWrapper struct {
	Data Data `json:"data"`
}

//Data -
type Data struct {
	Games GamesWrapper `json:"games"`
}

// GamesWrapper - wrapper
type GamesWrapper struct {
	Games []Game `json:"game"`
}

// Game - game
type Game struct {
	HomeCity       string    `json:"home_team_city"`
	HomeTeamName   string    `json:"home_team_name"`
	HomeWins       string    `json:"home_win"`
	HomeLosses     string    `json:"home_loss"`
	AwayCity       string    `json:"away_team_city"`
	AwayTeamName   string    `json:"away_team_name"`
	AwayWins       string    `json:"away_win"`
	AwayLosses     string    `json:"away_loss"`
	Linescore      Linescore `json:"linescore"`
	TimeDateString string    `json:"time_date"`
}

// GetHomeTeam - home team
func (g Game) GetHomeTeam() Team {
	runs, err := strconv.Atoi(g.Linescore.Runs.HomeRuns)
	if err != nil {
		panic(err)
	}
	wins, err := strconv.Atoi(g.HomeWins)
	if err != nil {
		panic(err)
	}
	losses, err := strconv.Atoi(g.HomeLosses)
	if err != nil {
		panic(err)
	}
	t, err := time.Parse("2006/01/02 3:04", g.TimeDateString)
	if err != nil {
		panic(err)
	}
	return Team{
		City:     g.HomeCity,
		TeamName: g.HomeTeamName,
		Wins:     wins,
		Losses:   losses,
		Runs:     runs,
		GameTime: t,
	}
}

// GetAwayTeam - away team
func (g Game) GetAwayTeam() Team {
	runs, err := strconv.Atoi(g.Linescore.Runs.AwayRuns)
	if err != nil {
		panic(err)
	}
	wins, err := strconv.Atoi(g.AwayWins)
	if err != nil {
		panic(err)
	}
	losses, err := strconv.Atoi(g.AwayLosses)
	if err != nil {
		panic(err)
	}
	t, err := time.Parse("2006/01/02 3:04", g.TimeDateString)
	if err != nil {
		panic(err)
	}
	return Team{
		City:     g.AwayCity,
		TeamName: g.AwayTeamName,
		Wins:     wins,
		Losses:   losses,
		Runs:     runs,
		GameTime: t,
	}
}

// Linescore - linescore
type Linescore struct {
	Runs Runs `json:"r"`
}

// Runs - runs
type Runs struct {
	HomeRuns string `json:"home"`
	AwayRuns string `json:"away"`
}

// Team - team
type Team struct {
	City     string
	TeamName string
	Runs     int
	Wins     int
	Losses   int
	GameTime time.Time
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

	teams := []Team{}
	for _, g := range d.Data.Games.Games {
		teams = append(teams, g.GetHomeTeam())
		teams = append(teams, g.GetAwayTeam())
	}
	fmt.Printf("%#+v", teams)

	csvFile, _ := os.Open("runs_pool_master/Sheet1-Table 1.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	scores := getTeamScores(lines)
	fmt.Printf("%#+v", scores)
}

func computeTeamScores(scores map[string]TeamScores, teams []Team) {
	completedTeamScores := []TeamScores{}
	updatedTeamScores := []TeamScores{}
	for _, team := range teams {
		for key, teamScore := range scores {
			if strings.Contains(key, team.TeamName) {
				if teamScore.IsCompleted() {
					completedTeamScores = append(completedTeamScores, teamScore)
					break
				}
				teamScore.Scores[team.Runs] = true
				if teamScore.IsCompleted() {
					teamScore.GamesPlayed = team.Wins + team.Losses
					teamScore.GameTimeCompleted = team.GameTime
					completedTeamScores = append(completedTeamScores, teamScore)
				}
				updatedTeamScores = append(updatedTeamScores, teamScore)
				break
			}
		}
	}
}

func getTeamScores(lines [][]string) map[string]TeamScores {
	m := map[string]TeamScores{}
	for i, line := range lines {
		if i == 0 {
			continue
		}
		if len(line) != 18 {
			panic(fmt.Sprintf("unknown columns"))
		}
		t := TeamScores{Scores: map[int]bool{}}
		t.TeamName = line[2]
		for i := 0; i < 14; i++ {
			fmt.Printf("here - line value: %v, i: %v\n", line[i+3], i)
			if line[i+3] != "" {
				t.Scores[i] = true
			}
		}
		if line[0] != "" {
			t.Player = &Player{
				Name: line[0],
				Paid: line[1] != "",
			}
		}
		m[t.TeamName] = t
	}
	return m
}

// TeamScores - run scores
type TeamScores struct {
	TeamName string
	//0, 1,... 12, >12-> 13
	Scores            map[int]bool
	Completed         bool
	GameTimeCompleted *time.Time
	GamesPlayed       int
	Player            *Player
}

// IsComplete - is it completed?
func (t TeamScores) IsComplete() bool {
	for _, c := range t.Scores {
		if !c {
			return false
		}
	}
	return true
}

//Player - player
type Player struct {
	Name string
	Paid bool
}
