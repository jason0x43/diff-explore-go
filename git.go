package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var renameThreshold = 50

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
		"-q",
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
	Adds    int
	Dels    int
	Path    string
	OldPath string
}

func gitDiffStat(start, end string) []stat {
	commit := start
	if end != "" {
		commit += ".." + end
	}

	out, err := exec.Command(
		"git",
		"diff",
		"--numstat",
		fmt.Sprintf("--find-renames=%d", renameThreshold),
		commit,
	).Output()
	if err != nil {
		log.Fatal(err)
	}

	var stats []stat

	outStr := strings.TrimSuffix(string(out), "\n")
	for _, line := range strings.Split(outStr, "\n") {
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, "\t")

		s := stat{}
		s.Adds, _ = strconv.Atoi(parts[0])
		s.Dels, _ = strconv.Atoi(parts[1])

		s.Path = strings.TrimSpace(parts[2])
		if strings.Contains(s.Path, " => ") {
			s.OldPath, s.Path, _ = strings.Cut(s.Path, " => ")
		}
		stats = append(stats, s)
	}

	return stats
}

type diffOptions struct {
	ignoreWhitespace bool
}

func gitDiff(start, end, path, oldPath string, options diffOptions) []string {
	commit := start
	command := "diff-index"
	if end != "" {
		command = "diff-tree"
		commit += ".." + end
	}

	args := []string{
		"git",
		command,
		"-M",
		"--patience",
		fmt.Sprintf("--find-renames=%d", renameThreshold),
		"-p",
	}

	if options.ignoreWhitespace {
		args = append(args, "-w")
	}

	args = append(args, commit, "--", path)

	if oldPath != "" {
		args = append(args, oldPath)
	}

	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	outStr := strings.TrimSuffix(string(out), "\n")
	return strings.Split(outStr, "\n")
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