package runner

import "testing"

func TestFeed(t *testing.T) {
	tests := []struct {
		name         string
		expect, fail string
		line         string
		wantDone     bool
		wantCode     int
		wantReason   string
	}{
		{"sentinel pass", "", "", "##srig-exit:0##", true, 0, "sentinel"},
		{"sentinel custom code", "", "", "##srig-exit:5##", true, 5, "sentinel"},
		{"sentinel with CR", "", "", "##srig-exit:0##\r", true, 0, "sentinel"},
		{"fail match", "PASS", "PANIC", "kernel PANIC now", true, 1, "fail"},
		{"expect match", "All tests passed", "", "All tests passed!", true, 0, "expect"},
		{"no match", "PASS", "FAIL", "boot ok", false, 0, ""},
		{"sentinel beats fail on same line", "", "exit", "##srig-exit:0## exit", true, 0, "sentinel"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEvaluator(tt.expect, tt.fail)
			if err != nil {
				t.Fatalf("NewEvaluator: %v", err)
			}
			o, done := e.Feed(tt.line)
			if done != tt.wantDone {
				t.Fatalf("done = %v, want %v", done, tt.wantDone)
			}
			if done && (o.ExitCode != tt.wantCode || o.Reason != tt.wantReason) {
				t.Errorf("outcome = %+v, want code %d reason %q", o, tt.wantCode, tt.wantReason)
			}
		})
	}
}

func TestTimeout(t *testing.T) {
	withExpect, _ := NewEvaluator("READY", "")
	if o := withExpect.Timeout(); o.ExitCode != 1 {
		t.Errorf("expect+timeout code = %d, want 1", o.ExitCode)
	}
	noExpect, _ := NewEvaluator("", "")
	if o := noExpect.Timeout(); o.ExitCode != 0 {
		t.Errorf("no-expect timeout code = %d, want 0", o.ExitCode)
	}
}

func TestNewEvaluatorBadRegex(t *testing.T) {
	if _, err := NewEvaluator("(", ""); err == nil {
		t.Fatal("expected regex compile error")
	}
}

func TestResultFrom(t *testing.T) {
	if r := ResultFrom(Outcome{ExitCode: 0, Reason: "expect"}, "sess_1"); r.Result != "pass" {
		t.Errorf("result = %q, want pass", r.Result)
	}
	if r := ResultFrom(Outcome{ExitCode: 1, Reason: "timeout"}, "sess_1"); r.Result != "fail" {
		t.Errorf("result = %q, want fail", r.Result)
	}
}
