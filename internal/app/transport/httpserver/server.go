package httpserver

import "payment-gateway/internal/app/services"

type HttpServer struct {
	userService                 UserService
	gatewayService              GatewayService
	transactionPublisherService services.TransactionPublisherService
}

func NewHttpServer(userService UserService, gtService GatewayService, txPublisherService services.TransactionPublisherService) HttpServer {
	return HttpServer{
		userService:                 userService,
		gatewayService:              gtService,
		transactionPublisherService: txPublisherService,
	}
}
