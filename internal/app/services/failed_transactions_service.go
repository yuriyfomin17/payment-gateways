package services

import (
	"context"
	"fmt"
	"log"
	"payment-gateway/internal/pkg"
	"time"
)

type FailedTransactionsService struct {
	rd                           pkg.RedisClient
	rb                           RabbitMqService
	faultToleranceService        FaultTolerance
	dataEncryptor                DataEncryptor
	transactionCallbackPublisher TransactionPublisherService
}

func NewFailedTransactionsService(rd pkg.RedisClient, rb RabbitMqService, faultTolerance FaultTolerance, dataEncryptor DataEncryptor, transactionCallbackPublisher TransactionPublisherService) FailedTransactionsService {
	return FailedTransactionsService{
		rd:                           rd,
		rb:                           rb,
		faultToleranceService:        faultTolerance,
		dataEncryptor:                dataEncryptor,
		transactionCallbackPublisher: transactionCallbackPublisher,
	}
}

func (fts FailedTransactionsService) ProcessFailedTransactionsEveryTimePeriod(timePeriod time.Duration) {
	go func(ctx context.Context) {
		ticker := time.NewTicker(timePeriod)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				failedTransactionsToProcess, err := fts.rd.GetListOfFailedCallbackTransactionsToProcess(ctx)
				if err != nil {
					log.Printf("Publishing failed transactions failed: %v", err)
					return
				}
				log.Printf("Publishing failed transactions: %v", failedTransactionsToProcess)
				for _, failedTransaction := range failedTransactionsToProcess {
					fmt.Println(failedTransaction.ID)
					err = fts.transactionCallbackPublisher.PublishTransaction(context.TODO(), failedTransaction.ID, failedTransaction.Status, failedTransaction.DataFormat)
					if err != nil {
						log.Printf("failed to publish transaction status update: %v", err)
						continue
					}
					err := fts.rd.DeleteFailedCallbackTransaction(ctx, failedTransaction.ID)
					if err != nil {
						log.Printf("failed to delete transaction from redis: %v", err)
					}

				}
			case <-ctx.Done():
				return
			}
		}
	}(context.TODO())
}
