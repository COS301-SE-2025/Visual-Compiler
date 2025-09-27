//go:build !coverage
// +build !coverage

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func checkAPIhealth() bool {
	resp, err := http.Get("http://localhost:8080/api/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func startAPI() *exec.Cmd {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("Backend started",)
	return cmd
}

func killAPI(cmd *exec.Cmd) {
	if cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println("Error killing API:", err)
		} else {
			fmt.Println("Backend stopped")
		}
	}
}

func waitUntilAlive(timeout time.Duration) bool {
	start := time.Now()
	for time.Since(start) < timeout {
		if checkAPIhealth() {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

func main() {
	fmt.Println("API no started")
	time.Sleep(3 * time.Second)
	api_cmd := startAPI()

	for {
		time.Sleep(3 * time.Second)

		killAPI(api_cmd)

		time.Sleep(3 * time.Second)

		api_cmd = startAPI()
		if waitUntilAlive(3 * time.Second) {
			fmt.Println("API restarted")
		} else {
			fmt.Println("API could not restart")
		}
	}
}
