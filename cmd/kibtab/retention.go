// Command kibtab starts the Kibtab instance.
//
// The retention job removes the audit rows that the retention days hold.
// Release v0.7.0 adds the job.
package main

import (
	"context"
	"log"
	"time"
)

// auditRetentionEvery is the pause between two runs of the retention job.
const auditRetentionEvery = 24 * time.Hour

// auditPruner removes the audit rows that the retention days hold. It
// returns the number of removed rows.
type auditPruner func(ctx context.Context) (int64, error)

// startAuditRetention runs the retention job in the background. It
// returns a function that stops the job.
func startAuditRetention(engine dbEngine) func() {
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(auditRetentionEvery)
	go runAuditRetention(ctx, ticker.C, engine.PruneAudit)
	return func() {
		ticker.Stop()
		cancel()
	}
}

// runAuditRetention prunes the audit rows on each tick. It stops when
// ctx ends. A failed run writes one log line and it waits for the next
// tick.
func runAuditRetention(ctx context.Context, ticks <-chan time.Time,
	prune auditPruner) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			removed, err := prune(ctx)
			if err != nil {
				log.Printf("audit retention: %v", err)
				continue
			}
			log.Printf("audit retention removed %d audit rows", removed)
		}
	}
}
