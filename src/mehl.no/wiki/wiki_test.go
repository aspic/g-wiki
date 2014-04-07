package main

import "testing"

func TestSanitizeLog(t *testing.T) {
    input := `{"Hash": "a926492", "Message":"escape"""here  ";", "Time":"28 hours ago"}`
    expected := `{"Hash": "a926492", "Message":" escape   here   ; ", "Time":"28 hours ago"}`
    sanitized := string(parseLog([]byte(input)))

    if expected != sanitized {
        t.Errorf("Sanitized: %s did not match %s", sanitized, expected)
    }
}
