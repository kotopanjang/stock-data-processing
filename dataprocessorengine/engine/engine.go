package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"

	"stock-data-processing/dataprocessorengine/storage/redis"
	"stock-data-processing/model"
)

type EngineHandler struct {
	storage *redis.RedisStorage
}

func NewEngine(storage *redis.RedisStorage) (*EngineHandler, error) {
	return &EngineHandler{
		storage: storage,
	}, nil
}

func (h *EngineHandler) Handle(ctx context.Context, message any) (err error) {
	msg, ok := message.(*sarama.ConsumerMessage)
	if !ok {
		log.Error().Msg("Not a kafka message")
		return errors.New("not a kafka message")
	}
	var rawData model.Raw

	if err = json.Unmarshal(msg.Value, &rawData); err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	err = h.processData(rawData)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	return
}

func (*EngineHandler) getPrice(rawData model.Raw) float64 {
	price, _ := strconv.ParseFloat(rawData.Price, 64)
	executionPrice, _ := strconv.ParseFloat(rawData.ExecutionPrice, 64)
	price += executionPrice

	return price
}
func (*EngineHandler) getQuantity(rawData model.Raw) float64 {
	quantity, _ := strconv.ParseFloat(rawData.Quantity, 64)
	executedQuantity, _ := strconv.ParseFloat(rawData.ExecutedQuantity, 64)
	quantity += executedQuantity

	return quantity
}

func (h *EngineHandler) processData(rawData model.Raw) error {
	price := h.getPrice(rawData)
	quantity := h.getQuantity(rawData)

	// get stored data from redis
	stockSum := &model.StockSummary{}
	if err := h.storage.GetByKey(rawData.StockCode, stockSum); err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	// calculate this => OpenPrice, LowestPrice, HighestPrice, PreviousPrice, ClosePrice
	stockSum.StockCode = rawData.StockCode
	switch cq := quantity; {
	case cq == 0:
		stockSum.PreviousPrice = price
	case cq > 0 && stockSum.OpenPrice == 0:
		stockSum.OpenPrice = price
		stockSum.LowestPrice = price
	default:
		stockSum.ClosePrice = price
		if price > stockSum.HighestPrice {
			stockSum.HighestPrice = price
		}
		if price < stockSum.LowestPrice {
			stockSum.LowestPrice = price
		}
	}

	// said => elements of a Transaction are Quantity and Price. Quantity and ExecutedQuantity are the same. Price and ExecutionPrice are the same.
	if price > 0 || quantity > 0 {
		stockSum.Transaction += 1
	}

	// said => Every Transaction with type of E and P are accountable for Volume, Value, and Average Price of a Stock.
	if rawData.Type == "E" || rawData.Type == "P" {
		stockSum.Volume += quantity
		stockSum.Value += quantity * price
		avg := stockSum.Value / stockSum.Volume
		stockSum.AveragePrice = math.Round(avg)
	}

	// unmarshall data
	stockSumJSONByte, err := json.Marshal(stockSum)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	msg := fmt.Sprintf("processing stock: %v, transaction: %v, qty: %v, price: %v", stockSum.StockCode, stockSum.Transaction, quantity, price)
	log.Info().Msg(msg)
	err = h.storage.Set(rawData.StockCode, stockSumJSONByte)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	return nil
}
