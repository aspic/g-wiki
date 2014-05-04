package main

import "testing"

func ParseLogLine(t *testing.T) {
    input := `a926492 28 hours ago "asdfasdf asdf"asdfasdf asdf test test!!`
    log := parseLog([]byte(input))
    hash := "a926492"
    msg := `"asdfasdf asdf"asdfasdf asdf test test!!`
    time := "28 hours ago"

    if log.Hash != hash {
        t.Errorf("Hash mismatch. Expected: %s got %s", hash, log.Hash)
    }
    if log.Message != msg {
        t.Errorf("Message mismatch. Expected: %s got %s", msg, log.Message)
    }
    if log.Time != time {
        t.Errorf("Time mismatch. Expected: %s got %s", time, log.Time)
    }
}
