package restserver

import (
	"fmt"

	"github.com/benchkram/bobc/application"
)

type Option func(s *S)

func WithHost(hostname, port string) Option {
	return func(s *S) {
		s.address = fmt.Sprintf("%s:%s", hostname, port)
	}
}

func WithUploadDir(uploadPath string) Option {
	return func(s *S) {
		s.uploadDir = uploadPath
	}
}

func WithArtifactService(srv application.Application) Option {
	return func(s *S) {
		s.app = srv
	}
}

func WithAuthenticator(authn Authenticator) Option {
	return func(s *S) {
		s.authenticator = authn
	}
}
