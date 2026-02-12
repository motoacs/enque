package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/motoacs/enque/backend/app"
)

type noopEmitter struct{}

func (n *noopEmitter) Emit(name string, payload any) {}

func main() {
	appData, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	baseDir := filepath.Join(appData, "Enque")
	_, err = app.New(baseDir, ".", &noopEmitter{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Enque backend initialized")
}
