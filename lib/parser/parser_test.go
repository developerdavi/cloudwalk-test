package parser

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	result, err := Parse("test/input.log")

	if err != nil {
		t.Error("failed to parse input file")
	}

	got, err := json.Marshal(result)

	if err != nil {
		t.Error("failed to parse results to json")
	}

	expected, err := os.ReadFile("test/output.json")

	if err != nil {
		t.Error("failed to read expected output")
	}

	if string(got) != string(expected) {
		t.Error("results are different from expected")
	}

	_, err = Parse("test/invalid.log")

	if err == nil {
		t.Error("failed to detect invalid input file")
	}
}

type LineTypeTest struct {
	line     string
	expected LineType
}

func TestParseLine(t *testing.T) {
	cases := []LineTypeTest{
		{
			`  0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0`,
			InitGame,
		},
		{
			` 15:00 Exit: Timelimit hit.`,
			Other,
		},
		{
			` 20:34 ClientConnect: 2`,
			Other,
		},
		{
			` 20:34 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\xian/default\hmodel\xian/default\g_redteam\\g_blueteam\\c1\4\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
			ClientUserinfoChanged,
		},
		{
			` 20:37 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
			ClientUserinfoChanged,
		},
		{
			` 20:37 ClientBegin: 2`,
			Other,
		},
		{
			` 20:37 ShutdownGame:`,
			ShutdownGame,
		},
		{
			` 20:37 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\bot_minplayers\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0`,
			InitGame,
		},
		{
			` 20:38 ClientConnect: 2`,
			Other,
		},
		{
			` 20:54 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`,
			Kill,
		},
		{
			` 22:18 Kill: 2 2 7: Is ga la mi do killed Isgalamido by MOD_ROCKET_SPLASH`,
			Kill,
		},
	}

	for _, c := range cases {
		got := parseLine(c.line)

		if got != c.expected {
			t.Errorf("got %v, expected %v", got, c.expected)
		}
	}
}

func TestParseKill(t *testing.T) {
	cases := []struct {
		line     string
		expected Death
	}{
		{
			` 20:54 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`,
			Death{
				killer: "<world>",
				victim: "Isgalamido",
				cause:  "MOD_TRIGGER_HURT",
			},
		},
		{
			` 22:18 Kill: 2 2 7: Is ga la mi do killed Isgalamido by MOD_ROCKET_SPLASH`,
			Death{
				killer: "Is ga la mi do",
				victim: "Isgalamido",
				cause:  "MOD_ROCKET_SPLASH",
			},
		},
		{
			`  2:38 Kill: 7 5 7: Assasinu Credi killed Dono da Bola by MOD_ROCKET_SPLASH`,
			Death{
				killer: "Assasinu Credi",
				victim: "Dono da Bola",
				cause:  "MOD_ROCKET_SPLASH",
			},
		},
	}

	for _, c := range cases {
		got, _ := parseKill(c.line)

		if got != c.expected {
			t.Errorf("got %v, expected %v", got, c.expected)
		}
	}
}
