package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Starting MASSTCLI execution via Go...")

	cmd := exec.Command("cmd", "/C", "scripts\\run_masstcli.bat")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("MASSTCLI execution failed:", err)
		os.Exit(1)
	}

	fmt.Println("MASSTCLI execution finished successfully")
}
