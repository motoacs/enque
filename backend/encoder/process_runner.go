package encoder

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

// RunResult holds the outcome of a single encoder process execution.
type RunResult struct {
	ExitCode      int
	ErrorMessage  string
	TimedOut      bool
	TimeoutReason string
	UsedJobObject bool
}

// ProgressCallback is called when a progress line is parsed.
type ProgressCallback func(progress Progress)

// LogCallback is called for each raw stderr line.
type LogCallback func(line string)

// ProcessRunner executes an encoder process with stderr monitoring.
type ProcessRunner struct {
	encoderPath string
	adapter     Adapter
	noOutputSec int
	noProgressSec int
}

// NewProcessRunner creates a process runner.
func NewProcessRunner(encoderPath string, adapter Adapter, noOutputSec, noProgressSec int) *ProcessRunner {
	return &ProcessRunner{
		encoderPath:   encoderPath,
		adapter:       adapter,
		noOutputSec:   noOutputSec,
		noProgressSec: noProgressSec,
	}
}

// Run executes the encoder process with the given args.
// ctx can be cancelled to kill the process.
// stderrWriter receives raw stderr bytes for logging.
// progressCb is called on parsed progress updates (throttled to ~500ms).
// logCb is called for each raw stderr line.
func (r *ProcessRunner) Run(ctx context.Context, args []string, stderrWriter io.Writer, progressCb ProgressCallback, logCb LogCallback) RunResult {
	cmd := exec.CommandContext(ctx, r.encoderPath, args...)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return RunResult{ExitCode: -1, ErrorMessage: fmt.Sprintf("stderr pipe: %v", err)}
	}

	// Also capture stdout (NVEncC writes very little to stdout)
	cmd.Stdout = nil

	if err := cmd.Start(); err != nil {
		return RunResult{ExitCode: -1, ErrorMessage: fmt.Sprintf("start: %v", err)}
	}

	// Set up Job Object for process management (Windows: ensures child cleanup)
	jo, joErr := SetupJobObject(cmd)
	usedJobObject := joErr == nil && jo != nil
	if jo != nil {
		defer jo.Close()
	}

	// Set up timeout guard
	var tg *TimeoutGuard
	cancelCh := make(chan struct{})
	if r.noOutputSec > 0 || r.noProgressSec > 0 {
		tg = NewTimeoutGuard(r.noOutputSec, r.noProgressSec, func(reason string) {
			select {
			case <-cancelCh:
			default:
				close(cancelCh)
			}
		})
		tg.Start()
		defer tg.Stop()
	}

	// Monitor stderr in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		r.readStderr(stderrPipe, stderrWriter, tg, progressCb, logCb)
	}()

	// Wait for either the process to finish or a timeout
	doneCh := make(chan error, 1)
	go func() {
		doneCh <- cmd.Wait()
	}()

	var result RunResult

	select {
	case err := <-doneCh:
		wg.Wait()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				result.ExitCode = exitErr.ExitCode()
			} else {
				result.ExitCode = -1
			}
			result.ErrorMessage = err.Error()
		} else {
			result.ExitCode = 0
		}
	case <-cancelCh:
		// Timeout triggered — kill the process
		if jo != nil {
			jo.Terminate(1)
		}
		killProcess(cmd)
		wg.Wait()
		<-doneCh
		if tg != nil {
			_, reason := tg.TimedOut()
			result.TimedOut = true
			result.TimeoutReason = reason
			result.ExitCode = -1
			result.ErrorMessage = fmt.Sprintf("timeout: %s", reason)
		}
	case <-ctx.Done():
		// Context cancelled (user abort/cancel)
		if jo != nil {
			jo.Terminate(1)
		}
		killProcess(cmd)
		wg.Wait()
		<-doneCh
		result.ExitCode = -1
		result.ErrorMessage = "cancelled"
	}

	result.UsedJobObject = usedJobObject
	return result
}

func (r *ProcessRunner) readStderr(pipe io.ReadCloser, writer io.Writer, tg *TimeoutGuard, progressCb ProgressCallback, logCb LogCallback) {
	// Use our custom scanner that handles \r and \n
	scanner := bufio.NewScanner(pipe)
	scanner.Split(scanCRLF)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var lastProgressEmit time.Time
	const progressThrottle = 500 * time.Millisecond

	for scanner.Scan() {
		line := scanner.Text()

		// Write raw bytes to stderr log
		if writer != nil {
			writer.Write([]byte(line + "\n"))
		}

		// Notify timeout guard
		if tg != nil {
			tg.NotifyOutput()
		}

		// Emit raw log line
		if logCb != nil {
			logCb(line)
		}

		// Parse progress
		progress := r.adapter.ParseProgress(line)
		if progress.Percent != nil {
			if tg != nil {
				tg.NotifyProgress(*progress.Percent)
			}

			// Throttle progress callbacks to ~500ms
			now := time.Now()
			if now.Sub(lastProgressEmit) >= progressThrottle {
				lastProgressEmit = now
				if progressCb != nil {
					progressCb(progress)
				}
			}
		}
	}
}

// scanCRLF is a bufio.SplitFunc that splits on \r and \n (NVEncC uses \r for progress).
func scanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	for i := 0; i < len(data); i++ {
		if data[i] == '\r' {
			if i+1 < len(data) && data[i+1] == '\n' {
				return i + 2, data[:i], nil
			}
			return i + 1, data[:i], nil
		}
		if data[i] == '\n' {
			return i + 1, data[:i], nil
		}
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

// killProcess terminates the process. Platform-specific kill is in jobobject_*.go files.
func killProcess(cmd *exec.Cmd) {
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
}
