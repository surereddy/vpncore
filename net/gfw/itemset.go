package gfw

import (
	"bufio"
	"os"
	"sort"
	"sync"
	"sync/atomic"
)

type itemMap map[string]bool

type ItemSet struct {
	filePath string
	moniter  int32

	sync.RWMutex
	itemMap
}

func NewItemSet(filePath string, capacity int) *ItemSet {
	il := &ItemSet{
		filePath:  filePath,
		itemMap: make(itemMap, capacity),
	}

	il.monitorConfigFile()
	return il
}

func (this *ItemSet) monitorConfigFile() {

	if exists, _ := IsFileExists(this.filePath); exists != true {
		return
	}

	if atomic.LoadInt32(&this.moniter) == 1 {
		return
	}

	configFileChanged := make(chan bool)
	go func() {
		for {
			select {
			case <-configFileChanged:
				this.Load()
			}
		}
	}()
	atomic.StoreInt32(&this.moniter, 1)
	MonitorFileChanegs(this.filePath, configFileChanged)

}

func (this *ItemSet) Hit(item string) bool {
	this.RLock()
	_, ok := this.itemMap[item]
	this.RUnlock()
	return ok
}

func (this *ItemSet) Save() bool {
	outFile, err := os.OpenFile(this.filePath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		return false
	}
	defer outFile.Close()

	this.RLock()
	var keys []string
	for k := range this.itemMap {
		keys = append(keys, k)
	}
	this.RUnlock()
	sort.Strings(keys)
	for _, k := range keys {
		outFile.WriteString(k + "\n")
	}

	this.monitorConfigFile()
	return true
}

func (this *ItemSet) Load() bool {
	inFile, _ := os.Open(this.filePath)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	this.Lock()
	defer this.Unlock()

	for scanner.Scan() {
		this.itemMap[scanner.Text()] = true
	}

	return len(this.itemMap) != 0
}

func (this *ItemSet) Clear() {
	this.Lock()
	this.itemMap = make(itemMap)
	this.Unlock()
}

func (this *ItemSet) AddItem(key string) {
	if _, ok := this.itemMap[key]; !ok {
		this.itemMap[key] = true
	}
}

func (this *ItemSet) IsEmpty() bool {
	this.RLock()
	isEmpty := (len(this.itemMap) == 0)
	this.RUnlock()
	return isEmpty
}
