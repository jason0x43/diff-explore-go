package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

type commit struct {
	Commit      string
	AuthorName  string
	AuthorEmail string
	Timestamp   int64
	Subject     string
}

func gitLog() []commit {
	out, err := exec.Command(
		"git",
		"log",
		"--date=iso8601-strict",
		"--pretty=format:{%n"+
			"  \"commit\": \"%H\",%n"+
			"  \"authorName\": \"%aN\",%n"+
			"  \"authorEmail\": \"%aE\",%n"+
			"  \"timestamp\": %at,%n"+
			"  \"subject\": \"%f\"%n"+
			"},",
	).Output()
	if err != nil {
		log.Fatal(err)
	}

	commitJson := "[" + strings.TrimRight(string(out), ",\n") + "]"

	var commits []commit
	json.Unmarshal([]byte(commitJson), &commits)

	return commits
}

type stat struct {
	Change rune
	Path string
}

func gitDiffStat(a string, b string) []stat {
	out, err := exec.Command(
		"git",
		"diff",
		"--name-status",
		a + ".." + b,
	).Output()
	if err != nil {
		log.Fatal(err)
	}

	var stats []stat

	for _, line := range strings.Split(string(out), "\n") {
		if len(line) > 0 {
			s := stat{}
			s.Change = rune(line[0])
			s.Path = strings.TrimSpace(line[1:])
			stats = append(stats, s)
		}
	}

	return stats
}
