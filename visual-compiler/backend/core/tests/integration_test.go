package integration_test

import(
	"testing"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services/commands"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
)

func TestStartServer(t *testing.T) {
	go core.StartServer()
	time.Sleep(500*time.Millisecond)
}

func TestRegister(t *testing.T) {
	 res,err := user.RegisterUser("tiaharripersad","t@gmail.com","tia1234$$")
	 if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	t.Logf("Registration succeeded: %s", res)
}