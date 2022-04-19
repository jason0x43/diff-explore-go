package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type notifyFunc func(string, string)

// func debounce(interval time.Duration, input chan string, cb func(arg string)) {
// 	var value string
// 	timer := time.NewTimer(interval)

// 	for {
// 		select {
// 		case value = <-input:
// 			timer.Reset(interval)
// 		case <-timer.C:
// 			if value != "" {
// 				cb(value)
// 			}
// 		}
// 	}
// }

func watchRepo(path string, notify notifyFunc) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	gitDir := filepath.Join(path, getGitDir())

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

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if !isIgnored(event.Name) {
					notify(fmt.Sprintf("%d", event.Op), event.Name)
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
