package v1

import (
	"github.com/go-server-template/internal/middleware"
	"github.com/go-server-template/internal/service"
	"github.com/go-server-template/internal/transport/rest/v1/auth"
	"github.com/go-server-template/internal/transport/rest/v1/onboarding"
)

type HttpServer struct {
	serv       *service.Service
	auth       *auth.Handler
	onboarding *onboarding.Handler
}

func NewHttpServer(s *service.Service) *HttpServer {
	handlerServer := &HttpServer{serv: s}

	// /v1
	v1 := s.Router.Group("/v1")
	{
		// /auth
		handlerServer.auth = auth.NewHandler(s.AuthService)
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", handlerServer.auth.Signup)
			auth.POST("/signin", handlerServer.auth.Signin)
			auth.POST("/refresh", handlerServer.auth.Refresh)
			auth.POST("/logout", handlerServer.auth.Logout)
		}

		// /onboarding
		handlerServer.onboarding = onboarding.NewHandler(*s)
		onboarding := v1.Group("/onboarding")
		onboarding.Use(middleware.ValidationMiddleware(s.Repo))
		{
			onboarding.GET("/test", handlerServer.onboarding.Test)
		}
	}

	return handlerServer
}

func (h *HttpServer) RunServer() error {
	if err := h.serv.Router.Run(":8080"); err != nil {
		return err
	}

	return nil
}
