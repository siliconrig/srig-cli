// Package runner decides pass/fail for `srig run` from streamed serial output.
package runner

import (
	"regexp"
	"strconv"
	"strings"
)

// Outcome is the terminal result of a run.
type Outcome struct {
	ExitCode int
	Reason   string // "sentinel" | "fail" | "expect" | "timeout" | "ran"
	Matched  string
}

var sentinelRe = regexp.MustCompile(`^##srig-exit:(\d+)##`)

// Evaluator matches serial lines against the sentinel and optional patterns.
type Evaluator struct {
	expect *regexp.Regexp
	fail   *regexp.Regexp
}

// NewEvaluator compiles the optional expect/fail regexes (empty string = unset).
func NewEvaluator(expect, fail string) (*Evaluator, error) {
	e := &Evaluator{}
	var err error
	if expect != "" {
		if e.expect, err = regexp.Compile(expect); err != nil {
			return nil, err
		}
	}
	if fail != "" {
		if e.fail, err = regexp.Compile(fail); err != nil {
			return nil, err
		}
	}
	return e, nil
}

// Feed evaluates one line of serial output. It returns (outcome, true) when a
// terminal condition is reached; precedence is sentinel > fail > expect.
func (e *Evaluator) Feed(line string) (Outcome, bool) {
	line = strings.TrimRight(line, "\r\n")
	if m := sentinelRe.FindStringSubmatch(line); m != nil {
		code, err := strconv.Atoi(m[1])
		if err != nil {
			code = 1 // malformed sentinel value → treat as failure, never a false pass
		}
		return Outcome{ExitCode: code, Reason: "sentinel", Matched: line}, true
	}
	if e.fail != nil && e.fail.MatchString(line) {
		return Outcome{ExitCode: 1, Reason: "fail", Matched: line}, true
	}
	if e.expect != nil && e.expect.MatchString(line) {
		return Outcome{ExitCode: 0, Reason: "expect", Matched: line}, true
	}
	return Outcome{}, false
}

// Timeout is the outcome when the watch window elapses with no terminal line:
// a fail if --expect was set (never matched), otherwise a clean pass.
func (e *Evaluator) Timeout() Outcome {
	if e.expect != nil {
		return Outcome{ExitCode: 1, Reason: "timeout"}
	}
	return Outcome{ExitCode: 0, Reason: "ran"}
}

// Result is the machine-readable form of an Outcome for --json.
type Result struct {
	Result    string `json:"result"` // "pass" | "fail"
	ExitCode  int    `json:"exit_code"`
	Reason    string `json:"reason"`
	Matched   string `json:"matched,omitempty"`
	SessionID string `json:"session_id"`
}

// ResultFrom maps an Outcome (exit 0 = pass) to a Result.
func ResultFrom(o Outcome, sessionID string) Result {
	res := "fail"
	if o.ExitCode == 0 {
		res = "pass"
	}
	return Result{Result: res, ExitCode: o.ExitCode, Reason: o.Reason, Matched: o.Matched, SessionID: sessionID}
}
