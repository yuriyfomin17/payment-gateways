package services

import (
	"context"
	"payment-gateway/internal/app/domain"
	"payment-gateway/internal/app/repository/model"
)

type UserServiceImpl struct {
	userRepo              UserRepository
	transactionRepo       TransactionRepository
	countryRepo           CountryRepository
	gatewayRepo           GatewayRepository
	faultToleranceService FaultTolerance
	redisService          RedisService
}

func NewUserService(userRepo UserRepository,
	transactionRepo TransactionRepository,
	countryRepo CountryRepository,
	gatewayRepo GatewayRepository,
	faultToleranceService FaultTolerance,
	redisService RedisService,
) UserServiceImpl {
	return UserServiceImpl{
		userRepo:              userRepo,
		transactionRepo:       transactionRepo,
		countryRepo:           countryRepo,
		gatewayRepo:           gatewayRepo,
		faultToleranceService: faultToleranceService,
		redisService:          redisService,
	}
}

func (us UserServiceImpl) FetchTxId(ctx context.Context) (int64, error) {
	id, err := us.transactionRepo.FetchTxId(ctx)
	if err != nil {
		return 0, domain.ErrTransactionNotCreated
	}
	return id, nil
}

func (us UserServiceImpl) ExecuteTransaction(ctx context.Context, tx domain.TransactionData) error {
	user, err := us.userRepo.GetUserByID(ctx, tx.UserID)
	if err != nil {
		return domain.ErrUserNotFound
	}
	country, err := us.countryRepo.GetCountryByID(ctx, user.CountryID(), tx.Currency)
	if err != nil {
		return domain.ErrTransactionNotCreated
	}
	tx.CountryID = country.ID()

	supportedGatewaysSortedByPriorities, err := us.gatewayRepo.GetSupportedGatewaysByCountrySortedByPriorities(ctx, country.ID(), tx.DataFormat)
	if err != nil {
		return domain.ErrTransactionNotCreated
	}

	for _, gateway := range supportedGatewaysSortedByPriorities {
		tx.GatewayID = gateway.ID()
		tx.Status = model.Pending.String()

		transaction, err := domain.NewTransaction(tx)
		if err != nil {
			return domain.ErrTransactionNotCreated
		}

		err = us.faultToleranceService.RetryOperation(func() error {
			return us.transactionRepo.CreateTransaction(ctx, transaction)
		}, 3)

		if err == nil {
			return nil
		}
	}

	return domain.ErrTransactionNotCreated
}
