package queue

import (
	"fmt"
	"os/exec"

	"github.com/yuta/enque/backend/logging"
)

// ExecutePostAction runs the post-complete action after session finishes.
func ExecutePostAction(action, customCommand string, logger *logging.AppLogger) error {
	switch action {
	case "none", "":
		return nil
	case "shutdown":
		if logger != nil {
			logger.Info("executing post-complete action: shutdown")
		}
		return platformShutdown()
	case "sleep":
		if logger != nil {
			logger.Info("executing post-complete action: sleep")
		}
		return platformSleep()
	case "custom":
		if customCommand == "" {
			return fmt.Errorf("custom post-complete command is empty")
		}
		if logger != nil {
			logger.Info("executing post-complete action: custom command: %s", customCommand)
		}
		return executeCustomCommand(customCommand)
	default:
		return fmt.Errorf("unknown post-complete action: %s", action)
	}
}

func executeCustomCommand(command string) error {
	cmd := exec.Command("cmd", "/C", command)
	return cmd.Start()
}
