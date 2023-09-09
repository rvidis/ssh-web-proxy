package server

import (
	"fmt"
	"net/http"
)

func Start(
	handler *SSHProxyHandler,
	port int,
) {
	fmt.Printf("Local server on port %d\n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler); err != nil {
		fmt.Printf("ERROR: Could not start the server on port %d: %v\n", port, err)
	}
}
