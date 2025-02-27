package services

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"log"
	"payment-gateway/internal/app/domain"
)

type TransactionPublisherServiceImpl struct {
	faultToleranceService  FaultTolerance
	rabbitService          RabbitMqService
	dataEncryptor          DataEncryptor
	failedTransactionRedis RedisService
}

func NewTransactionPublisherService(faultToleranceService FaultTolerance, rabbitService RabbitMqService, dataEncryptor DataEncryptor) TransactionPublisherServiceImpl {
	return TransactionPublisherServiceImpl{faultToleranceService: faultToleranceService, rabbitService: rabbitService, dataEncryptor: dataEncryptor}
}

func (tps TransactionPublisherServiceImpl) PublishTransaction(ctx context.Context, txId int64, statusToUpdate, dataFormat string) error {
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		err := func() error {

			switch dataFormat {
			case "application/json":
				var tx domain.TxStatusRabbitJson
				tx.TxID = txId
				tx.Status = statusToUpdate
				jsonData, err := json.Marshal(tx)
				if err != nil {
					cancel()
					return err
				}
				data, err := tps.dataEncryptor.MaskData(jsonData)
				if err != nil {
					cancel()
					return err
				}
				return tps.faultToleranceService.PublishWithCircuitBreaker(func() error {
					return tps.rabbitService.PublishJsonData(data)
				})
			case "text/xml", "application/xml":
				var tx domain.TxStatusXML
				tx.TxID = txId
				tx.Status = statusToUpdate
				soapData, err := xml.Marshal(tx)
				if err != nil {
					return domain.ErrTransactionStatus
				}

				data, err := tps.dataEncryptor.MaskData(soapData)
				if err != nil {
					cancel()
					return err
				}
				return tps.faultToleranceService.PublishWithCircuitBreaker(func() error {
					return tps.rabbitService.PublishSoapData(data)
				})
			default:
				return domain.ErrTransactionStatus
			}
		}()
		if err != nil {
			cancel()
			return
		}
		done <- struct{}{}
	}()
	select {
	case <-done:
		log.Printf("transaction status update published successfully")
		return nil
	case <-ctx.Done():
		log.Printf("failed to publish transaction status update")
		err := tps.failedTransactionRedis.CreateFailedCallbackTransaction(ctx, txId, statusToUpdate, dataFormat)
		if err != nil {
			return domain.ErrRedisCreateTransaction
		}
		return domain.ErrTransactionStatus
	}
}
