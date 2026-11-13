package leak

import (
	"context"
	"os"
	"testing"

	"go.uber.org/goleak"
)

// 1. A correct pool leaves no goroutines behind.
func TestPoolNoLeak(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan int)
	wg := StartWorkers(ctx, 3, jobs)

	for i := 0; i < 5; i++ {
		jobs <- i
	}
	cancel()
	wg.Wait()
}

// 2. goleak really does detect a leaked goroutine (test passes when a leak IS found).
func TestLeakIsDetected(t *testing.T) {
	stop := make(chan struct{})
	Leaky(stop)

	if err := goleak.Find(); err == nil {
		t.Fatal("expected goleak to find the leaked goroutine")
	} else {
		t.Log("goleak found the leak, as expected:\n", err)
	}

	close(stop) // clean up so later tests are unaffected
	if err := goleak.Find(); err != nil {
		t.Fatalf("goroutine should have exited: %v", err)
	}
}

// 3. Manual demo: this test FAILS on purpose. Enable it with LEAK_DEMO=1.
func TestDeliberateLeak(t *testing.T) {
	if os.Getenv("LEAK_DEMO") != "1" {
		t.Skip("set LEAK_DEMO=1 to see goleak fail")
	}
	defer goleak.VerifyNone(t)
	Leaky(nil) // nil channel: blocks forever
}
