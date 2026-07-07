package worker

import (
 "fmt"
 "log"
 "os/exec"
 "path/filepath"

 "stream-sync/internal/cache"
 "stream-sync/internal/queue"
 "stream-sync/internal/r2"
)

const RcloneRemote = "r2:yacine"

type Pool struct {
 Client        *r2.Client
 Queue         *queue.Queue
 Uploaded      *cache.Uploaded
 PlaylistQueue *queue.PlaylistQueue
}

func New(client *r2.Client, q *queue.Queue, uploaded *cache.Uploaded, pq *queue.PlaylistQueue) *Pool {
 return &Pool{
  Client:        client,
  Queue:         q,
  Uploaded:      uploaded,
  PlaylistQueue: pq,
 }
}

func uploadWithRclone(localPath, remoteKey string) error {
 dst := RcloneRemote + "/" + remoteKey

 cmd := exec.Command(
  "rclone",
  "copyto",
  localPath,
  dst,
  "--retries", "3",
  "--low-level-retries", "3",
  "--ignore-checksum",
 )

 out, err := cmd.CombinedOutput()
 if len(out) > 0 {
  log.Printf("rclone output: %s", string(out))
 }

 if err != nil {
  return fmt.Errorf("%w: %s", err, string(out))
 }

 return nil
}

func (p *Pool) Start(workers int) {
 for i := 0; i < workers; i++ {
  go func(id int) {

   log.Printf("Worker %d started", id)

   for {
    job := p.Queue.Pop()

    log.Printf("[Worker %d] rclone uploading %s", id, job.RemoteKey)

    if err := uploadWithRclone(job.LocalPath, job.RemoteKey); err != nil {
     log.Printf("[Worker %d] rclone upload failed: %s: %v", id, job.RemoteKey, err)
     continue
    }

    fileName := filepath.Base(job.RemoteKey)

    p.Uploaded.Add(job.Channel, fileName, 0)

    log.Printf("[Worker %d] rclone uploaded: %s", id, job.RemoteKey)

    if p.PlaylistQueue != nil {
     p.PlaylistQueue.Push(queue.PlaylistJob{
      Channel:      job.Channel,
      PlaylistPath: job.PlaylistPath,
     })
     log.Printf("[Worker %d] Playlist queued: %s", id, job.Channel)
    }
   }

  }(i + 1)
 }
}