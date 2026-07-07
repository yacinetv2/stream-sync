package watcher

import (
 "os"
 "time"
)

func waitUntilStable(path string) error {
 var lastSize int64 = -1
 stableCount := 0

 for {
  info, err := os.Stat(path)
  if err != nil {
   return err
  }

  size := info.Size()

  if size == lastSize && size > 0 {
   stableCount++
  } else {
   stableCount = 0
  }

  if stableCount >= 8 {
   return nil
  }

  lastSize = size
  time.Sleep(300 * time.Millisecond)
 }
}