package main

import (
	"net/http"
	"os"

	routes "github.com/syedomair/ex-pay-gateway/approve_micro/routes"
	"github.com/syedomair/ex-pay-gateway/lib/container"
)

func main() {
	// Setting up container
	c := container.New(map[string]string{
		container.ServiceNameEnvVar: os.Getenv(container.ServiceNameEnvVar),
		container.LogLevelEnvVar:    os.Getenv(container.LogLevelEnvVar),
		container.DatabaseURLEnvVar: os.Getenv(container.DatabaseURLEnvVar),
		container.PortEnvVar:        os.Getenv(container.PortEnvVar),
		container.SigningKeyEnvVar:  os.Getenv(container.SigningKeyEnvVar),
	})

	httpPort := c.Port()
	router := routes.NewRouter(c)

	c.Logger().Info("", "%q API Server listening on port %v", c.ServiceName(), httpPort)
	c.Logger().Alert("", "%v", http.ListenAndServe(":"+httpPort, router))

}
