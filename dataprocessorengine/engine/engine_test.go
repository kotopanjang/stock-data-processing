package engine

import (
	"context"
	"reflect"
	"stock-data-processing/dataprocessorengine/storage/redis"
	"testing"

	"github.com/Shopify/sarama"
)

func TestNewEngine(t *testing.T) {
	redisStore := redis.NewTestRedis()

	type args struct {
		storage *redis.RedisStorage
	}
	tests := []struct {
		name    string
		args    args
		want    *EngineHandler
		wantErr bool
	}{
		{
			name: "newengine_ok",
			args: args{
				storage: redisStore,
			},
			want: &EngineHandler{
				storage: redisStore,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEngine(tt.args.storage)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEngine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEngineHandler_Handle(t *testing.T) {
	redisStore := redis.NewTestRedis()

	type fields struct {
		storage *redis.RedisStorage
	}
	type args struct {
		ctx     context.Context
		message any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "handle_success",
			fields: fields{
				storage: redisStore,
			},
			args: args{
				ctx: context.Background(),
				message: &sarama.ConsumerMessage{
					Value: []byte(`{"type":"P","executed_quantity":"5","order_book":"35","execution_price":"4530","stock_code":"UNVR"}`),
				},
			},
			wantErr: false,
		},
		{
			name: "handle_err_invalid_msg",
			fields: fields{
				storage: redisStore,
			},
			args: args{
				ctx:     context.Background(),
				message: "",
			},
			wantErr: true,
		},
		{
			name: "handle_err_unmarshall",
			fields: fields{
				storage: redisStore,
			},
			args: args{
				ctx: context.Background(),
				message: &sarama.ConsumerMessage{
					Value: []byte(`{test:false}`),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &EngineHandler{
				storage: tt.fields.storage,
			}
			if err := h.Handle(tt.args.ctx, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("EngineHandler.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestEngineHandler_getPrice(t *testing.T) {
// 	type fields struct {
// 		storage *redis.RedisStorage
// 	}
// 	type args struct {
// 		rawData model.Raw
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   float64
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e := &EngineHandler{
// 				storage: tt.fields.storage,
// 			}
// 			if got := e.getPrice(tt.args.rawData); got != tt.want {
// 				t.Errorf("EngineHandler.getPrice() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestEngineHandler_getQuantity(t *testing.T) {
// 	type fields struct {
// 		storage *redis.RedisStorage
// 	}
// 	type args struct {
// 		rawData model.Raw
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   float64
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e := &EngineHandler{
// 				storage: tt.fields.storage,
// 			}
// 			if got := e.getQuantity(tt.args.rawData); got != tt.want {
// 				t.Errorf("EngineHandler.getQuantity() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestEngineHandler_processData(t *testing.T) {
// 	type fields struct {
// 		storage *redis.RedisStorage
// 	}
// 	type args struct {
// 		rawData model.Raw
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &EngineHandler{
// 				storage: tt.fields.storage,
// 			}
// 			if err := h.processData(tt.args.rawData); (err != nil) != tt.wantErr {
// 				t.Errorf("EngineHandler.processData() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
