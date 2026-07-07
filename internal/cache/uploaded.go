package cache

import "sync"

type channelData struct {
 set   map[string]struct{}
 order []string
}

type Uploaded struct {
 mu   sync.RWMutex
 data map[string]*channelData
}

func New() *Uploaded {
 return &Uploaded{
  data: make(map[string]*channelData),
 }
}

func (u *Uploaded) Add(channel, file string, keep int) []string {

 u.mu.Lock()
 defer u.mu.Unlock()

 if keep <= 0 {
  keep = 6
 }

 if _, ok := u.data[channel]; !ok {
  u.data[channel] = &channelData{
   set:   make(map[string]struct{}),
   order: make([]string, 0),
  }
 }

 ch := u.data[channel]

 if _, ok := ch.set[file]; ok {
  return nil
 }

 ch.set[file] = struct{}{}
 ch.order = append(ch.order, file)

 var deleted []string

 for len(ch.order) > keep {
  old := ch.order[0]
  ch.order = ch.order[1:]
  delete(ch.set, old)
  deleted = append(deleted, old)
 }

 return deleted
}

func (u *Uploaded) Exists(channel, file string) bool {

 u.mu.RLock()
 defer u.mu.RUnlock()

 ch, ok := u.data[channel]
 if !ok {
  return false
 }

 _, ok = ch.set[file]
 return ok
}

func (u *Uploaded) List(channel string) map[string]struct{} {

 u.mu.RLock()
 defer u.mu.RUnlock()

 result := make(map[string]struct{})

 ch, ok := u.data[channel]
 if !ok {
  return result
 }

 for file := range ch.set {
  result[file] = struct{}{}
 }

 return result
}

func (u *Uploaded) Remove(channel, file string) {

 u.mu.Lock()
 defer u.mu.Unlock()

 ch, ok := u.data[channel]
 if !ok {
  return
 }

 delete(ch.set, file)

 newOrder := make([]string, 0, len(ch.order))
 for _, f := range ch.order {
  if f != file {
   newOrder = append(newOrder, f)
  }
 }

 ch.order = newOrder
}

func (u *Uploaded) Clear(channel string) {

 u.mu.Lock()
 defer u.mu.Unlock()

 delete(u.data, channel)
}