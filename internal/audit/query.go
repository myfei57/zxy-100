package audit

import (
	"fmt"
)

func (r *Recorder) Recent(limit int) []Event {
	r.mu.Lock()
	seq := r.seq
	r.mu.Unlock()
	if limit < 1 {
		limit = 1
	}
	from := seq - int64(limit) + 1
	if from < 1 {
		from = 1
	}
	events := make([]Event, 0, limit)
	for index := from; index <= seq; index++ {
		var event Event
		if err := r.st.Get(fmt.Sprintf("audit/%d", index), &event); err == nil {
			events = append(events, event)
		}
	}
	return events
}

func (r *Recorder) Count() int64 {
	return r.Seq()
}

func (r *Recorder) Trim(keep int64) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	if keep < 1 {
		return 0
	}
	from := r.seq - keep + 1
	if from < 1 {
		from = 1
	}
	removed := 0
	for index := int64(1); index < from; index++ {
		if r.st.Exists(fmt.Sprintf("audit/%d", index)) {
			_ = r.st.Delete(fmt.Sprintf("audit/%d", index))
			removed++
		}
	}
	return removed
}

func (r *Recorder) BySource(source string, limit int) []Event {
	return filterEvents(r.Recent(r.BufferSize()), func(event Event) bool {
		return event.Source == source
	}, limit)
}

func (r *Recorder) ByAction(action string, limit int) []Event {
	return filterEvents(r.Recent(r.BufferSize()), func(event Event) bool {
		return event.Action == action
	}, limit)
}

func (r *Recorder) Summary() map[string]int {
	summary := make(map[string]int)
	for _, event := range r.Recent(r.BufferSize()) {
		summary[event.Action]++
	}
	return summary
}

func (r *Recorder) Since(seq int64, limit int) []Event {
	r.mu.Lock()
	current := r.seq
	r.mu.Unlock()
	if seq < 1 {
		seq = 1
	}
	if limit < 1 {
		limit = 1
	}
	events := make([]Event, 0, limit)
	for index := seq; index <= current; index++ {
		var event Event
		if err := r.st.Get(fmt.Sprintf("audit/%d", index), &event); err == nil {
			events = append(events, event)
			if len(events) >= limit {
				break
			}
		}
	}
	return events
}

func filterEvents(events []Event, keep func(Event) bool, limit int) []Event {
	if limit < 1 {
		limit = 1
	}
	out := make([]Event, 0, limit)
	for _, event := range events {
		if keep(event) {
			out = append(out, event)
			if len(out) >= limit {
				break
			}
		}
	}
	return out
}
