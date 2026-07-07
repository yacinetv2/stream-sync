package playlist

import (
 "sync"
 "time"
)

type SyncManager struct {
 mu     sync.Mutex
 timers map[string]*time.Timer
}

func NewSyncManager() *SyncManager {
 return &SyncManager{
  timers: make(map[string]*time.Timer),
 }
}