package services

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/rs/zerolog/log"
	"payment-gateway/internal/app/domain"
)

type CallbackServiceImpl struct {
	rabbitService RabbitMqService
	dataEncryptor DataEncryptor
	txRepo        TransactionRepository
}

func NewCallbackService(rabbitService RabbitMqService, dataEncryptor DataEncryptor, txRepo TransactionRepository) CallbackServiceImpl {
	return CallbackServiceImpl{
		rabbitService: rabbitService,
		dataEncryptor: dataEncryptor,
		txRepo:        txRepo,
	}
}

func (lps CallbackServiceImpl) StartListeningJsonMessages() {
	go func(ctx context.Context) {
		for jsonBytesMessage := range lps.rabbitService.GetJsonMessage() {
			unmaskedData, err := lps.dataEncryptor.UnmaskData(jsonBytesMessage)
			if err != nil {
				log.Error().Err(err).Msg("error while umasking json message")
				return
			}
			var txJsonMessage domain.TxStatusRabbitJson
			err = json.Unmarshal(unmaskedData, &txJsonMessage)
			err = lps.txRepo.UpdateTransactionStatus(ctx, txJsonMessage.TxID, txJsonMessage.Status)
			if err != nil {
				log.Error().Err(err).Msg("error while updating transaction status in json listener")
				return
			}
		}
		fmt.Println("json listener stopped")
	}(context.TODO())
}

func (lps CallbackServiceImpl) StartListeningSoapMessages() {
	go func(ctx context.Context) {
		for byteSoapMessage := range lps.rabbitService.GetSoapMessage() {
			unmaskedData, err := lps.dataEncryptor.UnmaskData(byteSoapMessage)
			if err != nil {
				log.Error().Err(err).Msg("error while unmasking xml message")
				return
			}
			var txJsonMessage domain.TxStatusXML
			err = xml.Unmarshal(unmaskedData, &txJsonMessage)
			if err != nil {
				log.Error().Err(err).Msg("error while unmarshalling xml message")
				return
			}
			err = lps.txRepo.UpdateTransactionStatus(ctx, txJsonMessage.TxID, txJsonMessage.Status)
			if err != nil {
				log.Error().Err(err).Msg("error while updating transaction status in xml listener")
				return
			}
		}
	}(context.TODO())
}
