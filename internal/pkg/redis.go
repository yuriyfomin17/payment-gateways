package pkg

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"payment-gateway/internal/app/domain"
	"strconv"
	"time"
)

const statusKey = "status"
const dataFormatKey = "dataFormat"

type RedisClient struct {
	rd *redis.Client
}

func ConnectToRedis(redisAddr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
}

func NewRedisClient(rd *redis.Client) RedisClient {
	return RedisClient{rd: rd}
}

func (r RedisClient) CreateFailedCallbackTransaction(ctx context.Context, txId int64, statusToUpdate, dataFormat string) error {
	txIdStr := strconv.FormatInt(txId, 2)

	jsonData, err := json.Marshal(map[string]any{
		statusKey:     statusToUpdate,
		dataFormatKey: dataFormat,
	})
	if err != nil {
		return domain.ErrRedisTransactionNotCreated
	}
	r.rd.Set(ctx, txIdStr, string(jsonData), time.Minute*10)

	return nil
}

func (r RedisClient) GetListOfFailedCallbackTransactionsToProcess(ctx context.Context) ([]domain.TransactionData, error) {
	keys, err := r.rd.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	txDataList := make([]domain.TransactionData, 0, cap(keys))
	for _, key := range keys {
		intKey, err := strconv.Atoi(key)
		if err != nil {
			return nil, domain.ErrRedisCouldNotExtractTransactions
		}
		jsonData, err := r.rd.Get(ctx, key).Result()
		if err != nil {
			return nil, domain.ErrRedisCouldNotExtractTransactions
		}
		var data map[string]any
		err = json.Unmarshal([]byte(jsonData), &data)
		txDataList = append(txDataList, domain.TransactionData{
			ID:         int64(intKey),
			Status:     data[statusKey].(string),
			DataFormat: data[dataFormatKey].(string),
		})
	}
	return txDataList, nil
}

func (r RedisClient) DeleteFailedCallbackTransaction(ctx context.Context, txId int64) error {
	txIdStr := strconv.FormatInt(txId, 10)
	return r.rd.Del(ctx, txIdStr).Err()
}
