package main

import (
	"fmt"
	"os"
	"sshwebproxy/server"
	"sshwebproxy/sshservice"
	"sshwebproxy/utils"
)

func main() {
	inputs := utils.ParseInputs()

	utils.PrintLogo()

	service, err := sshservice.NewService(
		*inputs.HostAddress,
		*inputs.HostPort,
		*inputs.HostUser,
		*inputs.IdentityFile,
		*inputs.HostBaseDir,
	)
	if err != nil {
		fmt.Printf("error when creating ssh servie: %+v\n", err)
		os.Exit(1)
	}

	handler := server.NewHandler(service)

	server.Start(handler, *inputs.HttpServerPort)
}
