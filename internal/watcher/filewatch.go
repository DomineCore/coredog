package watcher

// Listen some dir and call CoreFileService.Pub to publish the corefile infomation to CorefileService.

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FileWatcher struct {
	watch    *fsnotify.Watcher
	receiver chan string
}

func NewFileWatcher(recevier chan string) *FileWatcher {
	w := new(FileWatcher)
	w.watch, _ = fsnotify.NewWatcher()
	w.receiver = recevier
	return w
}

func pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {

		return false
	}
	return s.IsDir()

}

func (fw *FileWatcher) NFSWatch(dir string) error {
	for {
		go func() {
			ok, err := pathExist(dir)
			if err != nil {

			}
			if !ok {

			}
			if isDir(dir) {
				filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
					var creationTime time.Time
					fileInfo, _ := os.Stat(path)
					// 获取文件的创建日期
					switch runtime.GOOS {
					case "linux":
						creationTimeSpec := fileInfo.Sys().(*syscall.Stat_t).Ctim
						creationTime = time.Unix(int64(creationTimeSpec.Sec), int64(creationTimeSpec.Nsec))
					case "windows":
						creationTime = fileInfo.ModTime()
					}
					if time.Now().After(creationTime.Add(-1*time.Minute)) && time.Now().Before(creationTime) {
						fw.watch.Events <- fileInfo.Name()
					}
					return nil
				})
			}
			go fw.watchEvents()
		}()
	}
}

func (fw *FileWatcher) Watch(dir string) error {
	ok, err := pathExist(dir)
	if err != nil {
		return errors.Wrapf(err, "unexcepted error")
	}
	if !ok {
		return fmt.Errorf("dir is not exist:%s", dir)
	}
	if !isDir(dir) {
		return fmt.Errorf("input path is not a valid dir:%s", dir)
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = fw.watch.Add(path)
			if err != nil {
				return err
			}
			logrus.Infof("started watch corefile in:%s", path)
		}
		return nil
	})
	go fw.watchEvents()
	return nil
}

func (fw *FileWatcher) watchEvents() {
	for {
		select {
		case ev := <-fw.watch.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					file, err := os.Stat(ev.Name)
					if err == nil && file.IsDir() {
						fw.watch.Add(ev.Name)
						logrus.Infof("new subdir created,start to watch it:%s", ev.Name)
					}
				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					_, err := isFileWriteComplete(ev.Name)
					if err != nil {
						logrus.Errorf("file write")
					}
					logrus.Infof("capture a file:%s", ev.Name)
					// send file to receiver channel
					fw.receiver <- ev.Name
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					// if subdir removed
					// then remove the subdir's watch
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						fw.watch.Remove(ev.Name)
						logrus.Infof("subdir is removed, no more to watch:%s", ev.Name)
					}
				}

				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					// if subdir renamed
					// then the subdir's will be remove watch
					fw.watch.Remove(ev.Name)
					logrus.Infof("subdir is renamed, no more to watch:%s", ev.Name)
				}
			}
		case err := <-fw.watch.Errors:
			{
				logrus.Errorf("unexcepted watch error:%s", err)
				return
			}
		}
	}
}

func isFileWriteComplete(filePath string) (ok bool, err error) {
	for {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			err = errors.Wrapf(err, "failed to get file info")
			ok = false
			return ok, err
		}

		// get file size
		initialSize := fileInfo.Size()

		time.Sleep(1 * time.Second)

		// get file size again
		fileInfo, err = os.Stat(filePath)
		if err != nil {
			err = errors.Wrapf(err, "failed to get file info again")
			ok = false
			return ok, err
		}

		// check filesize changed
		if fileInfo.Size() == initialSize {
			ok = true
			return ok, nil
		} else {
			continue
		}
	}
}
