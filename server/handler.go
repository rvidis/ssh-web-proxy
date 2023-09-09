package server

import (
	"fmt"
	"net/http"
	"sshwebproxy/sshservice"
)

type SSHProxyHandler struct {
	service *sshservice.SSHService
}

func NewHandler(
	sshservice *sshservice.SSHService,
) *SSHProxyHandler {
	handler := &SSHProxyHandler{
		service: sshservice,
	}

	return handler
}

func (handler *SSHProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	content, err := handler.service.GetPageContent(path)
	if err != nil {
		fmt.Fprintf(w, "Error resolving the path: %s, %v", path, err)
		return
	}
	fmt.Fprintf(w, content)
}
