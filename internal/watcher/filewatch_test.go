package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestPathExist(t *testing.T) {
	os.Mkdir("./test", 0000)
	defer os.Remove("./test")
	exit, err := pathExist("./test")
	if err != nil {
		t.Fatal(err)
	}
	if !exit {
		t.Fail()
	}
}

func TestIsDir(t *testing.T) {
	os.Mkdir("./test", 0000)
	defer os.Remove("./test")
	if !isDir("./test") {
		t.Fail()
	}
}

func TestIsFileWriteComplete(t *testing.T) {
	f, err := os.OpenFile("./testfilewrite", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer os.Remove(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	execute := func(done bool) {
		ok, err := isFileWriteComplete("./testfilewrite")
		if err != nil {
			log.Fatal(err)
		}
		if ok != done {
			log.Fatal(fmt.Errorf("not excepted file state"))
		}
	}
	go execute(false)
	f.WriteString("hello world")
	execute(true)
}

// probelem
func TestWatch(t *testing.T) {
	re := make(chan string)
	fw := NewFileWatcher(re)

	os.Mkdir("./test", os.ModeAppend|os.ModePerm)
	defer os.Remove("./test")
	err := fw.Watch("./test")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		file := <-re
		_, f := filepath.Split(file)
		log.Println("file:" + f)
		if f != "testfile" {
			log.Fatal(fmt.Errorf("error"))
		}
	}()
	go os.OpenFile("./test/testfile", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer os.Remove("./test/testfile")
}
