package audit

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"wtc/internal/store"
)

type Event struct {
	ID        string `json:"id"`
	Seq       int64  `json:"seq"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
	Action    string `json:"action"`
	Message   string `json:"message"`
}

type Recorder struct {
	st         *store.Store
	mu         sync.Mutex
	seq        int64
	bufferSize int
}

func New(st *store.Store, bufferSize int) *Recorder {
	recorder := &Recorder{st: st, bufferSize: bufferSize}
	var seq int64
	if err := st.Get(store.KeyAuditSeq, &seq); err == nil {
		recorder.seq = seq
	}
	return recorder
}

func (r *Recorder) Append(source, action, message string) Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	event := Event{
		ID:        uuid.NewString(),
		Seq:       r.seq,
		Timestamp: time.Now().UnixMilli(),
		Source:    source,
		Action:    action,
		Message:   message,
	}
	_ = r.st.Put(fmt.Sprintf("audit/%d", r.seq), event)
	_ = r.st.Put(store.KeyAuditSeq, r.seq)
	return event
}

func (r *Recorder) Seq() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.seq
}

func (r *Recorder) BufferSize() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.bufferSize
}
