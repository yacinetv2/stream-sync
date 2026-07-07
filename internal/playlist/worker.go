package playlist

import (
 "log"
 "sync"
 "time"

 "stream-sync/internal/queue"
)

const PlaylistDelay = 0 * time.Second

func StartWorker(uploader *Uploader, pq *queue.PlaylistQueue) {

 var mu sync.Mutex
 timers := make(map[string]*time.Timer)

 go func() {

  log.Println("Playlist worker started")

  for {
   job := pq.Pop()

   mu.Lock()

   if t, ok := timers[job.Channel]; ok {
    t.Stop()
   }

   channel := job.Channel

   timers[channel] = time.AfterFunc(PlaylistDelay, func() {

    if err := uploader.UploadPlaylist(channel); err != nil {
     log.Printf("Playlist upload failed: %s: %v", channel, err)
    } else {
     log.Printf("Playlist uploaded after delay: %s", channel)
    }

    mu.Lock()
    delete(timers, channel)
    mu.Unlock()
   })

   mu.Unlock()

   log.Printf("Playlist delayed: %s", job.Channel)
  }
 }()
}