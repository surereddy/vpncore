package gfw

import (
	"os"
	"github.com/go-fsnotify/fsnotify"
	"errors"
)

func IsFileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.Mode() & os.ModeType == 0 {
			return true, nil
		}
		return false, errors.New(path + " exists but is not regular file")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MonitorFileChanegs(path string, changed chan bool) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Add(path)
	if err != nil {
		return err
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				if (event.Op & fsnotify.Write == fsnotify.Write) ||
					(event.Op & fsnotify.Rename == fsnotify.Rename) {
					changed <- true
				}
			case <-watcher.Errors:
				continue
			}
		}
	}()
	return nil
}

