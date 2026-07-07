package playlist

import (
 "errors"
 "os"
 "path"
 "strings"
)

var ErrPlaylistNotReady = errors.New("playlist not ready")

const DelaySegments = 0

func GetMediaSequence(localPlaylist string) string {
 data, err := os.ReadFile(localPlaylist)
 if err != nil {
  return ""
 }

 for _, line := range strings.Split(string(data), "\n") {
  line = strings.TrimSpace(line)
  if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:") {
   return strings.TrimPrefix(line, "#EXT-X-MEDIA-SEQUENCE:")
  }
 }

 return ""
}

func Rewrite(
 localPlaylist string,
 publicBaseURL string,
 channel string,
 isReady func(fileName string) bool,
 addEndList bool,
) ([]byte, error) {

 data, err := os.ReadFile(localPlaylist)
 if err != nil {
  return nil, err
 }

 publicBaseURL = strings.TrimRight(publicBaseURL, "/")

 lines := strings.Split(string(data), "\n")
 out := make([]string, 0, len(lines)+4)

 var pending []string
 var segments [][]string

 hasAllowCache := false

 for _, rawLine := range lines {
  line := strings.TrimSpace(rawLine)

  if line == "" {
   continue
  }

  if strings.HasPrefix(line, "#EXT-X-START") ||
   strings.HasPrefix(line, "#EXT-X-ENDLIST") ||
   strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE") {
   continue
  }

  if strings.HasPrefix(line, "#EXT-X-ALLOW-CACHE") {
   hasAllowCache = true
   continue
  }

  if strings.HasPrefix(line, "#EXTINF") ||
   strings.HasPrefix(line, "#EXT-X-BYTERANGE") {
   pending = append(pending, line)
   continue
  }

  if strings.HasSuffix(line, ".js") {
   fileName := path.Base(line)

   item := append([]string{}, pending...)
   item = append(item, publicBaseURL+"/"+channel+"/"+fileName)

   segments = append(segments, item)
   pending = nil
   continue
  }

  if line == "#EXTM3U" {
   out = append(out, line)

   if !hasAllowCache {
    out = append(out, "#EXT-X-ALLOW-CACHE:NO")
   }

   continue
  }

  out = append(out, line)
 }

 limit := len(segments) - DelaySegments
 if limit < 1 {
  limit = len(segments)
 }

 for i := 0; i < limit; i++ {
  out = append(out, segments[i]...)
 }

 return []byte(strings.Join(out, "\n")+"\n"), nil
}