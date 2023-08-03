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
	TotalKills   uint           `json:"total_kills"`
	Players      []string       `json:"players"`
	Kills        map[string]int `json:"kills"`
	Ranking      []string       `json:"ranking"`
	KillsByMeans map[string]int `json:"kills_by_means"`
}

type LineType uint8

const (
	Other LineType = iota
	InitGame
	ShutdownGame
	Kill
	ClientUserinfoChanged
)

var Causes = [...]string{
	"MOD_UNKNOWN",
	"MOD_SHOTGUN",
	"MOD_GAUNTLET",
	"MOD_MACHINEGUN",
	"MOD_GRENADE",
	"MOD_GRENADE_SPLASH",
	"MOD_ROCKET",
	"MOD_ROCKET_SPLASH",
	"MOD_PLASMA",
	"MOD_PLASMA_SPLASH",
	"MOD_RAILGUN",
	"MOD_LIGHTNING",
	"MOD_BFG",
	"MOD_BFG_SPLASH",
	"MOD_WATER",
	"MOD_SLIME",
	"MOD_LAVA",
	"MOD_CRUSH",
	"MOD_TELEFRAG",
	"MOD_FALLING",
	"MOD_SUICIDE",
	"MOD_TARGET_LASER",
	"MOD_TRIGGER_HURT",
	"MOD_NAIL",
	"MOD_CHAINGUN",
	"MOD_PROXIMITY_MINE",
	"MOD_KAMIKAZE",
	"MOD_JUICED",
	"MOD_GRAPPLE",
}

type Death struct {
	killer string
	victim string
	cause  string
}

func parseLine(line string) LineType {
	action_line := regexp.MustCompile(`\s*\d+:\d+\s(\w+):`)

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
	action_line := regexp.MustCompile(`\s*\d+:\d+\sKill:\s+\d+\s+\d+\s+\d+:\s([a-zA-Z\s<>]+)\skilled\s([a-zA-Z\s]+)\sby\s(MOD_\w*)`)

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
	action_line := regexp.MustCompile(`\s*\d+:\d+\sClientUserinfoChanged:\s+\d+\s+n\\([a-zA-Z\s]+)`)

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

func createKillsByMeans() map[string]int {
	kills_by_means := make(map[string]int)
	for _, cause := range Causes {
		kills_by_means[cause] = 0
	}
	return kills_by_means
}

func Parse(path string) ([]Match, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var matches []Match
	var match *Match

	is_world := regexp.MustCompile(`<world>`)

	// reading line by line, due to performance reasons
	for scanner.Scan() {
		line := scanner.Text()

		line_type := parseLine(line)

		// skip lines that are not relevant
		if line_type == Other {
			continue
		}

		// decide what to do according to the kind of the log
		switch line_type {
		case InitGame:
			matches = append(matches, Match{TotalKills: 0, Kills: make(map[string]int), KillsByMeans: make(map[string]int)})
			match = &matches[len(matches)-1]
			match.KillsByMeans = createKillsByMeans()
			continue
		case Kill:
			death, error := parseKill(line)
			if error != nil {
				log.Fatal(error)
			}

			match.KillsByMeans[death.cause]++
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

	return matches, nil
}
