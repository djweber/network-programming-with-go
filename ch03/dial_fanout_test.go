package ch03

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContextCancelFanOut(t *testing.T) {

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))

	sync := make(chan struct{})

	go func() {
		defer func() { sync <- struct{}{} }()

		var d net.Dialer
		d.Control = func(_, _ string, _ syscall.RawConn) error {
			time.Sleep(time.Second)
			return nil
		}

		conn, err := d.DialContext(ctx, "tcp", "10.0.0.1:80")

		if err != nil {
			t.Log(err)
			return
		}

		conn.Close()

		t.Error("connection did not time out")
	}()

	cancel()

	// wait for result from goroutine's defer
	<-sync

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %q", ctx.Err())
	}
}
