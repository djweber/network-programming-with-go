package ch03

import (
	"context"
	"fmt"
	"io"
	"time"
)

const defaultPingInterval = 30 * time.Second

func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
  var interval time.Duration

  select {
    case <-ctx.Done():
      fmt.Println("done")
    case interval = <-reset:
      fmt.Println("reset")
    default:
  }

  if interval <= 0 {
    interval = defaultPingInterval
  }

  timer := time.NewTimer(interval)

  defer func() {
    if !timer.Stop() {
      <-timer.C
    }
  }()

  for {
    select {
    case <- ctx.Done():
      return
    case newInterval := <-reset:
      if !timer.Stop() {
        <-timer.C
      }
      if newInterval > 0 {
        interval = newInterval
      }
    case <-timer.C:
      if _,err := w.Write([]byte("ping")); err != nil {
        return
      }
    }

    _ = timer.Reset(interval)
  }
}
