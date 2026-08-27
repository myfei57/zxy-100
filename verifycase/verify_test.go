package verifycase

import (
	"sync"
	"testing"

	"wtc/internal/yaw"
)

func TestWtcConcurrentYawCommand(t *testing.T) {
	drive := yaw.NewDrive(0)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for worker := 0; worker < 4; worker++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-start
			for step := 0; step < 300; step++ {
				drive.Request(float64(step%360), "manual")
			}
		}(worker)
	}
	close(start)
	wg.Wait()
	target, seq, _ := drive.Current()
	if seq != 1200 {
		t.Fatalf("expected 1200 acknowledged yaw requests, got %d", seq)
	}
	if !drive.Ack(seq, target) {
		t.Fatal("yaw state inconsistent at rest")
	}
}
