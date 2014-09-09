package main

import "testing"

func TestParseLogLine(t *testing.T) {
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

func TestListDirectoriesShouldReturnDirectories(t *testing.T) {
	path := "/test/test2"
	dirs := listDirectories(path)
	expectedLength := 3

	if len(dirs) != expectedLength {
		t.Errorf("Directories size should be %d, was %d", expectedLength, len(dirs))
	}
	if dirs[0].Path != "" {
		t.Errorf("Wrong root path, was %s", dirs[0].Path)
	}
}
