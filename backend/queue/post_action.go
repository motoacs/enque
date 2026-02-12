package queue

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/motoacs/enque/backend/model"
)

func runPostAction(cfg model.AppConfig) error {
	switch cfg.PostCompleteAction {
	case model.PostActionNone:
		return nil
	case model.PostActionShutdown:
		if runtime.GOOS != "windows" {
			return fmt.Errorf("shutdown action is only supported on windows")
		}
		return exec.Command("shutdown", "/s", "/t", "0").Run()
	case model.PostActionSleep:
		return setSystemSleep()
	case model.PostActionCustom:
		if cfg.PostCompleteCommand == "" {
			return fmt.Errorf("post_complete_command is empty")
		}
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", cfg.PostCompleteCommand).Run()
		}
		return exec.Command("sh", "-c", cfg.PostCompleteCommand).Run()
	default:
		return fmt.Errorf("unknown post action: %s", cfg.PostCompleteAction)
	}
}
