package routes

import (
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/lib"
)

type HelloRoutes struct {
	logger          lib.Logger
	env             lib.Env
	hellohandler    handlers.HelloHandler
	request_handler lib.RequestHandler
}

func NewHelloRoutes(
	logger lib.Logger,
	env lib.Env,
	handler lib.RequestHandler,
	hellohandler handlers.HelloHandler,
) HelloRoutes {
	return HelloRoutes{
		logger:          logger,
		env:             env,
		request_handler: handler,
		hellohandler:    hellohandler,
	}
}

func (r HelloRoutes) Setup() {
	r.request_handler.Mux.Get("/hello", r.hellohandler.Hello)
}
