package parser

import (
	"bufio"
	"errors"
	"log"
	"os"
	"regexp"
	"sort"
)

type Match struct {
	TotalKills uint           `json:"total_kills"`
	Players    []string       `json:"players"`
	Kills      map[string]int `json:"kills"`
	Ranking    []string       `json:"ranking"`
}

type LineType uint8

const (
	Other LineType = iota
	InitGame
	ShutdownGame
	Kill
	ClientUserinfoChanged
)

type Death struct {
	killer string
	victim string
	cause  string
}

func parseLine(line string) LineType {
	action_line := regexp.MustCompile(`\s+\d+:\d+\s(\w+):`)

	if !action_line.MatchString(line) {
		return Other
	}

	action := action_line.FindStringSubmatch(line)[1]

	switch action {
	case "InitGame":
		return InitGame
	case "ShutdownGame":
		return ShutdownGame
	case "Kill":
		return Kill
	case "ClientUserinfoChanged":
		return ClientUserinfoChanged
	default:
		return Other
	}
}

func parseKill(line string) (Death, error) {
	action_line := regexp.MustCompile(`\s+\d+:\d+\sKill:\s+\d+\s+\d+\s+\d+:\s([a-zA-Z\s<>]+)\skilled\s([a-zA-Z\s]+)\sby\s(MOD_\w*)`)

	if !action_line.MatchString(line) {
		return Death{}, errors.New("invalid kill line: " + line)
	}

	matches := action_line.FindStringSubmatch(line)

	return Death{
		killer: matches[1],
		victim: matches[2],
		cause:  matches[3],
	}, nil
}

func parseClientUserinfoChanged(line string) (string, error) {
	action_line := regexp.MustCompile(`\s+\d+:\d+\sClientUserinfoChanged:\s+\d+\s+n\\([a-zA-Z\s]+)`)

	if !action_line.MatchString(line) {
		return "", errors.New("invalid userinfo line: " + line)
	}

	player := action_line.FindStringSubmatch(line)[1]

	return player, nil
}

func sortKills(match *Match) {
	keys := make([]string, 0, len(match.Kills))
	for k := range match.Kills {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return match.Kills[keys[i]] > match.Kills[keys[j]]
	})
	match.Ranking = keys
}

func Parse() []Match {
	file, err := os.Open("input/qgames.log")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var matches []Match
	var match *Match

	is_world := regexp.MustCompile(`<world>`)

	for scanner.Scan() {
		line := scanner.Text()

		line_type := parseLine(line)

		if line_type == Other {
			continue
		}

		switch line_type {
		case InitGame:
			matches = append(matches, Match{TotalKills: 0, Kills: make(map[string]int)})
			match = &matches[len(matches)-1]
			continue
		case Kill:
			death, error := parseKill(line)
			if error != nil {
				log.Fatal(error)
			}

			match.TotalKills++

			if !is_world.MatchString(death.killer) {
				match.Kills[death.killer]++
			} else {
				match.Kills[death.victim]--
			}

			sortKills(match)
			continue
		case ClientUserinfoChanged:
			player, error := parseClientUserinfoChanged(line)
			if error != nil {
				log.Fatal(error)
			}
			found := false
			for i := 0; i < len(match.Players); i++ {
				if match.Players[i] == player {
					found = true
					break
				}
			}
			if !found {
				match.Players = append(match.Players, player)
			}
		}
	}

	return matches
}
