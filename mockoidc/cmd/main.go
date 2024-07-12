package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mockoidc/mockoidc"
)

func main() {
	m, _ := mockoidc.Run(4044)
	m.ClientID = "foobar"
	m.ClientSecret = "frobozz"
	defer m.Shutdown()

	cfg := m.Config()
	cfgJson, _ := json.Marshal(cfg)
	fmt.Println(string(cfgJson))

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	os.Stderr.WriteString("Press CTRL-C to terminate\n")
	<-done
}
