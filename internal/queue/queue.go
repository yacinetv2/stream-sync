package queue

type Job struct {
 Channel      string
 LocalPath    string
 RemoteKey    string
 PlaylistPath string
}

type PlaylistJob struct {
 Channel      string
 PlaylistPath string
}

type Queue struct {
 Jobs chan Job
}

type PlaylistQueue struct {
 Jobs chan PlaylistJob
}

func New(size int) *Queue {
 return &Queue{
  Jobs: make(chan Job, size),
 }
}

func (q *Queue) Push(job Job) {
 q.Jobs <- job
}

func (q *Queue) Pop() Job {
 return <-q.Jobs
}

func NewPlaylist(size int) *PlaylistQueue {
 return &PlaylistQueue{
  Jobs: make(chan PlaylistJob, size),
 }
}

func (q *PlaylistQueue) Push(job PlaylistJob) {
 q.Jobs <- job
}

func (q *PlaylistQueue) Pop() PlaylistJob {
 return <-q.Jobs
}