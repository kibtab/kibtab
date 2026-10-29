package main

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

// retentionRun starts runAuditRetention in the background. It returns a
// tick channel and a function that stops the job and waits for it.
func retentionRun(prune auditPruner) (chan time.Time, func()) {
	ticks := make(chan time.Time, 1)
	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		runAuditRetention(ctx, ticks, prune)
		close(done)
	}()
	return ticks, func() {
		cancel()
		<-done
	}
}

func TestRunAuditRetentionPrunesEachTick(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	called := make(chan int64, 1)
	ticks, stop := retentionRun(func(context.Context) (int64, error) {
		called <- 4
		return 4, nil
	})

	ticks <- time.Now()
	if removed := <-called; removed != 4 {
		t.Fatalf("prune: got %d, want 4", removed)
	}
	stop()

	if !strings.Contains(buf.String(), "removed 4 audit rows") {
		t.Fatalf("log: got %q, want a removal line", buf.String())
	}
}

func TestRunAuditRetentionLogsAFailedPrune(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	called := make(chan int64, 1)
	ticks, stop := retentionRun(func(context.Context) (int64, error) {
		called <- 0
		return 0, errors.New("the database refused the delete")
	})

	ticks <- time.Now()
	<-called
	stop()

	failure := "audit retention: the database refused the delete"
	if !strings.Contains(buf.String(), failure) {
		t.Fatalf("log: got %q, want the failure line", buf.String())
	}
}
