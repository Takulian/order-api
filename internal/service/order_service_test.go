package service_test

import (
	"context"
	"errors"
	"log/slog"
	"order-api/internal/dto"
	"order-api/internal/mocks"
	"order-api/internal/model"
	"order-api/internal/repository"
	"order-api/internal/service"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

type testDeps struct {
	repo      *mocks.MockOrderRepository
	cache     *mocks.MockOrderCache
	publisher *mocks.MockPublisher
}

func setupTest(t *testing.T) (*service.OrderService, *testDeps) {
	t.Helper()

	ctrl := gomock.NewController(t)

	deps := &testDeps{
		repo:      mocks.NewMockOrderRepository(ctrl),
		cache:     mocks.NewMockOrderCache(ctrl),
		publisher: mocks.NewMockPublisher(ctrl),
	}

	logger := slog.New(slog.DiscardHandler)
	s := service.NewOrderService(deps.repo, deps.cache, deps.publisher, logger)

	t.Cleanup(func() {
		defer ctrl.Finish()
	})

	return s, deps
}

func TestOrderService_GetAll(t *testing.T) {
	expected := []model.Order{
		{
			ID:       1,
			Customer: "Andi",
			Product:  "Laptop",
			Quantity: 2,
			Status:   "Pending",
		},
		{
			ID:       2,
			Customer: "Budi",
			Product:  "Handphone",
			Quantity: 3,
			Status:   "Pending",
		},
	}
	tests := []struct {
		name    string
		want    []model.Order
		wantErr bool
		setup   func(d *testDeps)
	}{
		{
			name:    "get orders with all connection ok and have cache",
			want:    expected,
			wantErr: false,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetAll(gomock.Any(), "orders").Return(expected, nil)
				d.repo.EXPECT().GetAll().Times(0)
				d.cache.EXPECT().SetAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:    "get orders with all connection ok no cache then get from database and cache it",
			want:    expected,
			wantErr: false,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetAll(gomock.Any(), "orders").Return(nil, errors.New("no cache data"))
				d.repo.EXPECT().GetAll().Return(expected, nil)
				d.cache.EXPECT().SetAll(gomock.Any(), "orders", expected, gomock.Any()).Return(nil)
			},
		},
		{
			name:    "get orders when repository returns error",
			want:    nil,
			wantErr: true,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetAll(gomock.Any(), "orders").Return(nil, errors.New("no cache data"))
				d.repo.EXPECT().GetAll().Return(nil, errors.New("db error"))
				d.cache.EXPECT().SetAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:    "get orders with database ok but redis error and response ok",
			want:    expected,
			wantErr: false,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetAll(gomock.Any(), "orders").Return(nil, errors.New("no cache data"))
				d.repo.EXPECT().GetAll().Return(expected, nil)
				d.cache.EXPECT().SetAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("redis error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, deps := setupTest(t)
			tt.setup(deps)
			got, gotErr := s.GetAll(context.Background())
			if (gotErr != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %+v, want = %+v", got, tt.want)
			}
		})
	}
}

func TestOrderService_GetByID(t *testing.T) {
	expected := model.Order{
		ID:       1,
		Customer: "Andi",
		Product:  "Laptop",
		Quantity: 12,
		Status:   "Pending",
	}
	tests := []struct {
		name    string
		id      int
		want    model.Order
		wantErr bool
		setup   func(d *testDeps)
	}{
		{
			name:    "input id with normal condition with cache has data then succes response no error",
			id:      1,
			want:    expected,
			wantErr: false,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(expected, nil)
				d.repo.EXPECT().GetByID(1).Times(0)
				d.cache.EXPECT().SetByID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:    "input id with normal condition with cache no data then succes response no error and caching data",
			id:      1,
			want:    expected,
			wantErr: false,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(model.Order{}, errors.New("cache not found"))
				d.repo.EXPECT().GetByID(1).Return(expected, nil)
				d.cache.EXPECT().SetByID(gomock.Any(), gomock.Any(), expected, gomock.Any()).Return(nil)
			},
		},
		{
			name:    "input unknown id with normal condition then not found response",
			id:      67,
			want:    model.Order{},
			wantErr: true,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetByID(gomock.Any(), 67, gomock.Any()).Return(model.Order{}, errors.New("cache not found"))
				d.repo.EXPECT().GetByID(67).Return(model.Order{}, repository.ErrOrderNotFound)
				d.cache.EXPECT().SetByID(gomock.Any(), gomock.Any(), expected, gomock.Any()).Times(0)
			},
		},
		{
			name:    "input id but redis connection failed then response still success",
			id:      1,
			want:    expected,
			wantErr: false,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(model.Order{}, errors.New("connection error"))
				d.repo.EXPECT().GetByID(1).Return(expected, nil)
				d.cache.EXPECT().SetByID(gomock.Any(), gomock.Any(), expected, gomock.Any()).Return(errors.New("connection error"))
			},
		},
		{
			name:    "input id when repository returns error",
			id:      1,
			want:    model.Order{},
			wantErr: true,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(model.Order{}, errors.New("cache not found"))
				d.repo.EXPECT().GetByID(1).Return(model.Order{}, errors.New("db error"))
				d.cache.EXPECT().SetByID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:    "input id when repository ok but redis error with response ok",
			id:      1,
			want:    expected,
			wantErr: false,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(model.Order{}, errors.New("cache not found"))
				d.repo.EXPECT().GetByID(1).Return(expected, nil)
				d.cache.EXPECT().SetByID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("redis error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, deps := setupTest(t)
			tt.setup(deps)
			got, gotErr := s.GetByID(context.Background(), tt.id)
			if (gotErr != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %+v, want = %+v", got, tt.want)
			}
		})
	}
}

func TestOrderService_Create(t *testing.T) {
	dbErr := errors.New("db error")
	expected := model.Order{
		ID:       1,
		Customer: "Andi",
		Product:  "Laptop",
		Quantity: 12,
		Status:   "Pending",
	}
	tests := []struct {
		name    string
		req     dto.CreateOrderRequest
		want    model.Order
		wantErr error
		setup   func(d *testDeps)
	}{
		{
			name: "create order with all input correct all connection ok then response success",
			req: dto.CreateOrderRequest{
				Customer: "Andi",
				Product:  "Laptop",
				Quantity: 12,
			},
			want:    expected,
			wantErr: nil,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Create(model.Order{
					Customer: "Andi",
					Product:  "Laptop",
					Quantity: 12,
					Status:   "Pending",
				}).Return(expected, nil)
				d.cache.EXPECT().Del(gomock.Any(), "orders").Return(nil)
				d.publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "create order with input empty customer all connection ok then response failed with error customer is required",
			req: dto.CreateOrderRequest{
				Customer: "",
				Product:  "Laptop",
				Quantity: 12,
			},
			want:    model.Order{},
			wantErr: service.ErrCustomerRequired,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Create(gomock.Any()).Times(0)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
				d.publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "create order with input empty product all connection ok then response failed with error product is required",
			req: dto.CreateOrderRequest{
				Customer: "Andi",
				Product:  "",
				Quantity: 12,
			},
			want:    model.Order{},
			wantErr: service.ErrProductRequired,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Create(gomock.Any()).Times(0)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
				d.publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "create order with input zero or minus all connection ok then response failed with error invalid quantity",
			req: dto.CreateOrderRequest{
				Customer: "Andi",
				Product:  "Laptop",
				Quantity: 0,
			},
			want:    model.Order{},
			wantErr: service.ErrInvalidQuantity,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Create(gomock.Any()).Times(0)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
				d.publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "create order when repository returns error",
			req: dto.CreateOrderRequest{
				Customer: "Andi",
				Product:  "Laptop",
				Quantity: 12,
			},
			want:    model.Order{},
			wantErr: dbErr,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Create(model.Order{
					Customer: "Andi",
					Product:  "Laptop",
					Quantity: 12,
					Status:   "Pending",
				}).Return(model.Order{}, dbErr)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
				d.publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "create order when repository ok but redis error and response ok",
			req: dto.CreateOrderRequest{
				Customer: "Andi",
				Product:  "Laptop",
				Quantity: 12,
			},
			want:    expected,
			wantErr: nil,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Create(model.Order{
					Customer: "Andi",
					Product:  "Laptop",
					Quantity: 12,
					Status:   "Pending",
				}).Return(expected, nil)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Return(errors.New("redis error"))
				d.publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, deps := setupTest(t)
			tt.setup(deps)
			got, gotErr := s.Create(context.Background(), tt.req)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, gotErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %+v, want = %+v", got, tt.want)
			}
		})
	}
}

func TestOrderService_Update(t *testing.T) {
	dbErr := errors.New("db error")
	current := model.Order{
		ID:       1,
		Customer: "Andi",
		Product:  "Laptop",
		Quantity: 12,
		Status:   "Pending",
	}
	updated := model.Order{
		ID:       1,
		Customer: "Budi",
		Product:  "Monitor",
		Quantity: 5,
		Status:   "Pending",
	}

	tests := []struct {
		name    string
		id      int
		req     dto.UpdateOrderRequest
		want    model.Order
		wantErr error
		setup   func(d *testDeps)
	}{
		{
			name: "update order with valid input",
			id:   1,
			req: dto.UpdateOrderRequest{
				Customer: "Budi",
				Product:  "Monitor",
				Quantity: 5,
			},
			want:    updated,
			wantErr: nil,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(current, nil)
				d.repo.EXPECT().Update(1, updated).Return(updated, nil)
				d.cache.EXPECT().Del(gomock.Any(), "orders").Return(nil)
				d.cache.EXPECT().Del(gomock.Any(), "orders:1").Return(nil)
			},
		},
		{
			name: "update order with empty customer",
			id:   1,
			req: dto.UpdateOrderRequest{
				Customer: "",
				Product:  "Monitor",
				Quantity: 5,
			},
			want:    model.Order{},
			wantErr: service.ErrCustomerRequired,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "update order with empty product",
			id:   1,
			req: dto.UpdateOrderRequest{
				Customer: "Budi",
				Product:  "",
				Quantity: 5,
			},
			want:    model.Order{},
			wantErr: service.ErrProductRequired,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "update order with invalid quantity",
			id:   1,
			req: dto.UpdateOrderRequest{
				Customer: "Budi",
				Product:  "Monitor",
				Quantity: 0,
			},
			want:    model.Order{},
			wantErr: service.ErrInvalidQuantity,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "update order when get by id fails",
			id:   1,
			req: dto.UpdateOrderRequest{
				Customer: "Budi",
				Product:  "Monitor",
				Quantity: 5,
			},
			want:    model.Order{},
			wantErr: dbErr,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(model.Order{}, errors.New("cache not found"))
				d.repo.EXPECT().GetByID(1).Return(model.Order{}, dbErr)
				d.repo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "update order when repository update fails",
			id:   1,
			req: dto.UpdateOrderRequest{
				Customer: "Budi",
				Product:  "Monitor",
				Quantity: 5,
			},
			want:    model.Order{},
			wantErr: dbErr,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(current, nil)
				d.repo.EXPECT().Update(1, updated).Return(model.Order{}, dbErr)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "update order when repository ok but redis error and response ok",
			id:   1,
			req: dto.UpdateOrderRequest{
				Customer: "Budi",
				Product:  "Monitor",
				Quantity: 5,
			},
			want:    updated,
			wantErr: nil,
			setup: func(d *testDeps) {
				d.cache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(current, nil)
				d.repo.EXPECT().Update(1, updated).Return(updated, nil)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Return(errors.New("redis error"))
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Return(errors.New("redis error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, deps := setupTest(t)
			tt.setup(deps)
			got, gotErr := s.Update(context.Background(), tt.id, tt.req)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, gotErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %+v, want = %+v", got, tt.want)
			}
		})
	}
}

func TestOrderService_Delete(t *testing.T) {
	dbErr := errors.New("db error")
	tests := []struct {
		name    string
		id      int
		wantErr error
		setup   func(d *testDeps)
	}{
		{
			name:    "delete order successfully",
			id:      1,
			wantErr: nil,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Delete(1).Return(nil)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Return(nil)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:    "delete order with unknown id then not found response",
			id:      67,
			wantErr: service.ErrOrderNotFound,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Delete(67).Return(repository.ErrOrderNotFound)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:    "delete order when repository returns error",
			id:      1,
			wantErr: dbErr,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Delete(1).Return(dbErr)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:    "delete order when repository ok but redis error and response ok",
			id:      1,
			wantErr: nil,
			setup: func(d *testDeps) {
				d.repo.EXPECT().Delete(1).Return(nil)
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Return(errors.New("redis error"))
				d.cache.EXPECT().Del(gomock.Any(), gomock.Any()).Return(errors.New("redis error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, deps := setupTest(t)
			tt.setup(deps)
			gotErr := s.Delete(context.Background(), tt.id)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, gotErr)
			}
		})
	}
}

func TestOrderService_Checkout(t *testing.T) {
	tests := []struct {
		name    string
		req     dto.CreateOrderRequest
		wantErr error
		setup   func(*testDeps)
	}{
		{
			name: "hit checkout when publisher conn ok then response ok",
			req: dto.CreateOrderRequest{
				Customer: "Andi",
				Product:  "Laptop",
				Quantity: 10,
			},
			wantErr: nil,
			setup: func(d *testDeps) {
				d.publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "hit checkout with customer empty when publisher conn ok then response error",
			req: dto.CreateOrderRequest{
				Customer: "",
				Product:  "Laptop",
				Quantity: 10,
			},
			wantErr: service.ErrCustomerRequired,
			setup: func(d *testDeps) {
				d.publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "hit checkout with product empty when publisher conn ok then response error",
			req: dto.CreateOrderRequest{
				Customer: "Andi",
				Product:  "",
				Quantity: 10,
			},
			wantErr: service.ErrProductRequired,
			setup: func(d *testDeps) {
				d.publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "hit checkout with zero quantity when publisher conn ok then response error",
			req: dto.CreateOrderRequest{
				Customer: "Andi",
				Product:  "Laptop",
				Quantity: 0,
			},
			wantErr: service.ErrInvalidQuantity,
			setup: func(d *testDeps) {
				d.publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, deps := setupTest(t)
			tt.setup(deps)
			gotErr := s.Checkout(t.Context(), tt.req)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, gotErr)
			}
		})
	}
}
