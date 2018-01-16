package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

var timestamp string

func initAndLog(logger *Logger) {
	InitLoggers(logger)

	Stdout.Printf("Stdout.Printf")
	Stderr.Printf("Stderr.Printf")
	Trace.Printf("Trace.Printf")
	Debug.Printf("Debug.Printf")
	Info.Printf("Info.Printf")
	Warning.Printf("Warning.Printf")
	Error.Printf("Error.Printf")
}

func getExpected(ts string, loggers ...string) string {
	out := []string{}
	for _, l := range loggers {
		switch l {
		case "stdout":
			out = append(out, fmt.Sprintf("%s log_test.go:19: Stdout.Printf", ts))
		case "stderr":
			out = append(out, fmt.Sprintf("%s log_test.go:20: Stderr.Printf", ts))
		case "trace":
			out = append(out, fmt.Sprintf("TRACE: %s log_test.go:21: Trace.Printf", ts))
		case "debug":
			out = append(out, fmt.Sprintf("DEBUG: %s log_test.go:22: Debug.Printf", ts))
		case "info":
			out = append(out, fmt.Sprintf("INFO: %s Info.Printf", ts))
		case "warning":
			out = append(out, fmt.Sprintf("WARNING: %s log_test.go:24: Warning.Printf", ts))
		case "error":
			out = append(out, fmt.Sprintf("ERROR: %s log_test.go:25: Error.Printf", ts))
		default:
			panic("unknown")
		}
	}
	out = append(out, "")
	return strings.Join(out, "\n")
}

func TestLog_init_empty(t *testing.T) {
	if os.Getenv("GO_TEST_FORKED_PROCESS") == "1" {
		initAndLog(nil)
		os.Exit(0)
		return
	}
	out, err := helperCommand(t)
	if err != nil {
		t.Errorf("got %q", err)
	}
	expected := getExpected(timestamp, "stdout", "stderr")
	if string(out) != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}

func TestLog_trace(t *testing.T) {
	if os.Getenv("GO_TEST_FORKED_PROCESS") == "1" {
		initAndLog(&Logger{Trace: os.Stdout})
		os.Exit(0)
		return
	}

	out, err := helperCommand(t)
	if err != nil {
		t.Errorf("got %q", err)
	}
	expected := getExpected(timestamp, "stdout", "stderr", "trace")
	if string(out) != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}

func TestLog_all(t *testing.T) {
	if os.Getenv("GO_TEST_FORKED_PROCESS") == "1" {
		initAndLog(&Logger{
			os.Stdout,
			os.Stdout,
			os.Stdout,
			os.Stdout,
			os.Stderr,
		})
		os.Exit(0)
		return
	}

	out, err := helperCommand(t)
	if err != nil {
		t.Errorf("got %q", err)
	}
	expected := getExpected(timestamp, "stdout", "stderr", "trace", "debug", "info", "warning", "error")
	if string(out) != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}

func TestLog(t *testing.T) {
	if os.Getenv("GO_TEST_FORKED_PROCESS") == "1" {
		initAndLog(&Logger{
			os.Stdout,
			os.Stdout,
			os.Stdout,
			os.Stdout,
			os.Stderr,
		})
		os.Exit(0)
		return
	}

	out, err := helperCommand(t)
	if err != nil {
		t.Errorf("got %q", err)
	}
	expected := getExpected(timestamp, "stdout", "stderr")
	if string(out) != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}

func ExampleLog_stdout() {
	// init logger
	InitLoggers(&Logger{
		ioutil.Discard,
		ioutil.Discard,
		os.Stdout,
		os.Stdout,
		os.Stderr,
	})

	Stdout.Printf("Stdout.Printf")
	Stderr.Printf("Stderr.Printf")
}

func helperCommand(t *testing.T) (string, error) {
	cs := []string{fmt.Sprintf("-test.run=%s", helperCaller()), "--"}
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_TEST_FORKED_PROCESS=1"}

	out, err := cmd.CombinedOutput()

	return string(out), err
}

func helperCaller() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	me := runtime.FuncForPC(pc)
	if me == nil {
		return "unnamed"
	}
	full := me.Name()
	p := strings.Split(full, ".")
	return p[len(p)-1]
}

func init() {
	timestamp = time.Now().UTC().Format("2006/01/02 15:04:05")
}
