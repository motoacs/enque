package encoder

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type RunSpec struct {
	Executable           string
	Argv                 []string
	NoOutputTimeoutSec   int
	NoProgressTimeoutSec int
}

type RunResult struct {
	ExitCode      int
	UsedJobObject bool
	TimedOut      bool
	Err           error
}

type ProgressObserver func(line string)
type LogObserver func(line string)

type ProcessRunner interface {
	Run(ctx context.Context, spec RunSpec, onProgress ProgressObserver, onLog LogObserver) RunResult
}

type OSProcessRunner struct{}

func NewOSProcessRunner() *OSProcessRunner {
	return &OSProcessRunner{}
}

func (r *OSProcessRunner) Run(ctx context.Context, spec RunSpec, onProgress ProgressObserver, onLog LogObserver) RunResult {
	cmd := exec.CommandContext(ctx, spec.Executable, spec.Argv...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return RunResult{ExitCode: -1, Err: err}
	}
	if err := cmd.Start(); err != nil {
		return RunResult{ExitCode: -1, Err: err}
	}
	job, usedJobObject, jobErr := attachJobObject(cmd.Process)
	if jobErr != nil {
		usedJobObject = false
	}
	defer closeJobObject(job)

	guard := NewTimeoutGuard(time.Duration(spec.NoOutputTimeoutSec)*time.Second, time.Duration(spec.NoProgressTimeoutSec)*time.Second)
	guard.MarkLine()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		scanner.Split(splitCRLF)
		for scanner.Scan() {
			line := scanner.Text()
			guard.MarkLine()
			onLog(line)
			onProgress(line)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	for {
		select {
		case err := <-done:
			wg.Wait()
			if err == nil {
				return RunResult{ExitCode: 0, UsedJobObject: usedJobObject}
			}
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					return RunResult{ExitCode: status.ExitStatus(), UsedJobObject: usedJobObject, Err: err}
				}
			}
			return RunResult{ExitCode: 1, UsedJobObject: usedJobObject, Err: err}
		case <-ticker.C:
			if guard.IsOutputTimedOut() || guard.IsProgressTimedOut() {
				_ = terminateWithJobObject(job, cmd.Process.Pid)
				<-done
				wg.Wait()
				return RunResult{ExitCode: -1, UsedJobObject: usedJobObject, TimedOut: true, Err: fmt.Errorf("process timed out")}
			}
		case <-ctx.Done():
			_ = terminateWithJobObject(job, cmd.Process.Pid)
			<-done
			wg.Wait()
			return RunResult{ExitCode: -1, UsedJobObject: usedJobObject, Err: ctx.Err()}
		}
	}
}

func splitCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		j := i + 1
		for j < len(data) && (data[j] == '\r' || data[j] == '\n') {
			j++
		}
		return j, bytes.TrimSpace(data[:i]), nil
	}
	if atEOF {
		return len(data), bytes.TrimSpace(data), nil
	}
	return 0, nil, nil
}

func terminateProcessTree(pid int) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(pid))
		return cmd.Run()
	}
	cmd := exec.Command("kill", "-9", fmt.Sprint(pid))
	if out, err := cmd.CombinedOutput(); err != nil {
		if strings.TrimSpace(string(out)) == "" {
			return err
		}
		return fmt.Errorf("kill failed: %w (%s)", err, string(out))
	}
	return nil
}

func drain(r io.Reader) {
	_, _ = io.Copy(io.Discard, r)
}
