package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type notifyFunc func(string, string)

func watchRepo(path string, notify notifyFunc) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	gitDir := filepath.Join(path, getGitDir())

	go func() {
		err = filepath.WalkDir(
			path,
			func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() {
					if isIgnored(path) || path == gitDir {
						return fs.SkipDir
					}

					return watcher.Add(path)
				}
				return nil
			},
		)

		if err != nil {
			log.Fatal(err)
		}

		notify("ready", "")

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if !isIgnored(event.Name) {
					absPath, err := filepath.Abs(event.Name)
					if err != nil {
						absPath = path
					}
					notify(fmt.Sprintf("%d", event.Op), absPath)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	return watcher
}
