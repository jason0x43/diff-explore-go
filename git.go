package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type commit struct {
	Commit      string
	Decoration  string
	AuthorName  string
	AuthorEmail string
	Timestamp   int64
	Subject     string
}

func isIgnored(path string) bool {
	err := exec.Command(
		"git",
		"check-ignore",
		path,
	).Run()

	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if exiterr.ProcessState.ExitCode() == 1 {
				return false
			}
			log.Fatal(err)
		} else {
			log.Fatal(err)
		}
	}

	return true
}

func getGitDir() string {
	out, err := exec.Command(
		"git",
		"rev-parse",
		"--git-dir").Output()
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(string(out))
}

func gitLog() []commit {
	out, err := exec.Command(
		"git",
		"log",
		"--date=iso8601-strict",
		"--decorate",
		"--pretty=format:{%n"+
			"  \"commit\": \"%H\",%n"+
			"  \"decoration\": \"%d\",%n"+
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
	Path   string
}

func gitDiffStat(start, end string) []stat {
	commit := start
	if end != "" {
		commit += ".." + end
	}

	out, err := exec.Command(
		"git",
		"diff",
		"--name-status",
		commit,
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

func gitDiff(start, end, path string) []string {
	commit := start
	command := "diff-index"
	if end != "" {
		command = "diff-tree"
		commit += ".." + end
	}

	out, err := exec.Command(
		"git",
		command,
		"-p",
		commit,
		path,
	).Output()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(out), "\n")
}

type decor struct {
	branches []string
	tags     []string
	refs     []string
}

func parseDecoration(decoration string) (d decor) {
	decoration = strings.Trim(decoration, " () ")
	decors := strings.Split(decoration, ", ")

	headBranchRe := regexp.MustCompile(`^\w+ ->`)
	tagRe := regexp.MustCompile(`^tag: `)
	refsRe := regexp.MustCompile(`^\w+/`)

	for _, decor := range decors {
		if headBranchRe.MatchString(decor) {
			b := strings.Split(decor, " -> ")[1]
			d.branches = append(d.branches, b)
		} else if tagRe.MatchString(decor) {
			t := strings.SplitN(decor, ": ", 2)[1]
			d.tags = append(d.tags, t)
		} else if refsRe.MatchString(decor) {
			d.refs = append(d.refs, decor)
		} else {
			d.branches = append(d.branches, decor)
		}
	}

	return
}
