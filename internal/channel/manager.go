package channel

import (
 "os"
 "path/filepath"
)

func List(root string) ([]string, error) {

 entries, err := os.ReadDir(root)
 if err != nil {
  return nil, err
 }

 var channels []string

 for _, entry := range entries {

  if !entry.IsDir() {
   continue
  }

  channels = append(channels, filepath.Join(root, entry.Name()))
 }

 return channels, nil
}