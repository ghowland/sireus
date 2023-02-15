package util

import (
	"sync"
)

var (
	AllLocks = make(map[string]*sync.RWMutex)
)

func LockAcquire(lockKey string) {
	_, ok := AllLocks[lockKey]
	if !ok {
		AllLocks[lockKey] = &sync.RWMutex{}
	}

	//log.Printf("Lock: %s", lockKey)
	AllLocks[lockKey].Lock()
}

func LockRelease(lockKey string) {
	//log.Printf("Unlock: %s", lockKey)
	AllLocks[lockKey].Unlock()
}
