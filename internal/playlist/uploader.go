package playlist

import (
 "log"
 "os"
 "path/filepath"

 "stream-sync/internal/config"
 "stream-sync/internal/r2"
)

type Uploader struct {
 Client *r2.Client
 Config *config.Config
}

func NewUploader(client *r2.Client, cfg *config.Config) *Uploader {
 return &Uploader{
  Client: client,
  Config: cfg,
 }
}

func (u *Uploader) UploadPlaylist(channel string) error {

 localPlaylist := filepath.Join(
  u.Config.WatchPath,
  channel,
  "index.m3u8",
 )

 if _, err := os.Stat(localPlaylist); err != nil {
  if os.IsNotExist(err) {
   log.Printf("Playlist skipped, file not found: %s", channel)
   return nil
  }
  return err
 }

 data, err := Rewrite(
  localPlaylist,
  u.Config.R2.PublicURL,
  channel,
  func(fileName string) bool {
   return true
  },
  false,
 )
 if err != nil {
  return err
 }

 objectName := filepath.ToSlash(filepath.Join("live", channel, "index.m3u8"))

 if err := u.Client.UploadBytes(data, objectName); err != nil {
  return err
 }

 log.Printf("Playlist updated: %s", channel)

 return nil
}