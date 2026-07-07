package events

type UploadEvent struct {
 Channel      string
 LocalPath    string
 RemoteKey    string
 PlaylistPath string
}

type Bus struct {
 Uploads chan UploadEvent
}

func New() *Bus {
 return &Bus{
  Uploads: make(chan UploadEvent, 1000),
 }
}