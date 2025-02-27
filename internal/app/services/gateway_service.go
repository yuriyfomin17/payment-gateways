package services

import "context"

type GatewayServiceImpl struct {
	gatewayRepo GatewayRepository
}

func NewGatewayService(gatewayRepo GatewayRepository) GatewayServiceImpl {
	return GatewayServiceImpl{gatewayRepo: gatewayRepo}
}

func (gt GatewayServiceImpl) UpdateGatewayPriority(ctx context.Context, gatewayID int64, priority string) error {
	return gt.gatewayRepo.UpdateGatewayPriority(ctx, gatewayID, priority)
}
