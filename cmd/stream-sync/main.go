package main

import (
 "log"

 "stream-sync/internal/cache"
 "stream-sync/internal/channel"
 "stream-sync/internal/config"
 "stream-sync/internal/playlist"
 "stream-sync/internal/queue"
 "stream-sync/internal/r2"
 "stream-sync/internal/watcher"
 "stream-sync/internal/worker"
)

func main() {

 cfg, err := config.Load("config.yaml")
 if err != nil {
  log.Fatal(err)
 }

 log.Println("Configuration loaded")

 client, err := r2.NewClient(cfg)
 if err != nil {
  log.Fatal(err)
 }

 log.Println("Connected to Cloudflare R2")

 q := queue.New(1000)
 pq := queue.NewPlaylist(1000)

 uploaded := cache.New()

 uploader := playlist.NewUploader(client, cfg)

 pool := worker.New(client, q, uploaded, pq)
 pool.Start(cfg.Workers)

 channels, err := channel.List(cfg.WatchPath)
 if err != nil {
  log.Fatal(err)
 }

 log.Printf("Found %d channels", len(channels))

 playlist.StartWorker(uploader, pq)

 for _, ch := range channels {
  w, err := watcher.New(q, pq)
  if err != nil {
   log.Fatal(err)
  }

  if err := w.WatchChannel(ch); err != nil {
   log.Fatal(err)
  }
 }

 log.Println("Stream Sync started.")

 select {}
}