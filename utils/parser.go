package utils

import "flag"

type Inputs struct {
	HostAddress    *string
	HostPort       *int
	HostUser       *string
	IdentityFile   *string
	HostBaseDir    *string
	HttpServerPort *int
}

func ParseInputs() *Inputs {

	hostAddress := flag.String("ha", "0.0.0.0", "Host Address")
	hostPort := flag.Int("hp", 22, "Host Port - port to use with remote host")
	hostUser := flag.String("hu", "", "Host User")
	identityFile := flag.String("i", "", "Identity file")
	hostBaseDir := flag.String("hd", "", "Host Base Directory - directory containing the content on the host")
	httpServerPort := flag.Int("sp", 8080, "Local HTTP server port")

	flag.Parse()

	return &Inputs{
		HostAddress:    hostAddress,
		HostPort:       hostPort,
		HostUser:       hostUser,
		IdentityFile:   identityFile,
		HostBaseDir:    hostBaseDir,
		HttpServerPort: httpServerPort,
	}
}
