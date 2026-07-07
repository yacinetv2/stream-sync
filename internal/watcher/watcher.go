package watcher

import (
 "log"
 "path/filepath"
 "sync"
 "time"

 "stream-sync/internal/queue"

 "github.com/fsnotify/fsnotify"
)

type Watcher struct {
 fs            *fsnotify.Watcher
 queue         *queue.Queue
 playlistQueue *queue.PlaylistQueue

 mu     sync.Mutex
 timers map[string]*time.Timer
}

func New(q *queue.Queue, pq *queue.PlaylistQueue) (*Watcher, error) {
 fs, err := fsnotify.NewWatcher()
 if err != nil {
  return nil, err
 }

 return &Watcher{
  fs:            fs,
  queue:         q,
  playlistQueue: pq,
  timers:        make(map[string]*time.Timer),
 }, nil
}

func (w *Watcher) WatchChannel(path string) error {
 if err := w.fs.Add(path); err != nil {
  return err
 }

 log.Println("Watching:", path)

 go func() {
  for {
   select {
   case event := <-w.fs.Events:
    if event.Name == "" {
     continue
    }

    if filepath.Base(event.Name) == "index.m3u8" {
     continue
    }

    if event.Op&fsnotify.Create == 0 {
     continue
    }

    if filepath.Ext(event.Name) != ".js" {
     continue
    }

    if err := waitUntilStable(event.Name); err != nil {
     log.Println(err)
     continue
    }

    channel := filepath.Base(filepath.Dir(event.Name))
    fileName := filepath.Base(event.Name)

    w.queue.Push(queue.Job{
     Channel:      channel,
     LocalPath:    event.Name,
     RemoteKey:    filepath.ToSlash(filepath.Join(channel, fileName)),
     PlaylistPath: filepath.Join(filepath.Dir(event.Name), "index.m3u8"),
    })

    log.Printf("Upload queued: %s", filepath.ToSlash(filepath.Join(channel, fileName)))

   case err := <-w.fs.Errors:
    if err != nil {
     log.Println(err)
    }
   }
  }
 }()

 return nil
}