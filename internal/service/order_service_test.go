package service_test

import (
	"context"
	"errors"
	"order-api/internal/dto"
	"order-api/internal/mocks"
	"order-api/internal/model"
	"order-api/internal/service"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func setupTest(t *testing.T) (*service.OrderService, *mocks.MockOrderRepository, *mocks.MockOrderCache) {
	t.Helper()

	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	mockCache := mocks.NewMockOrderCache(ctrl)

	s := service.NewOrderService(mockRepo, mockCache)

	t.Cleanup(func() {
		defer ctrl.Finish()
	})

	return s, mockRepo, mockCache
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
		setup   func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache)
	}{
		{
			name:    "get orders with all connection ok and have cache",
			want:    expected,
			wantErr: false,
			setup: func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache) {
				mockCache.EXPECT().GetAll(gomock.Any(), "orders").Return(expected, nil)
				mockRepo.EXPECT().GetAll().Times(0)
				mockCache.EXPECT().SetAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:    "get orders with all connection ok no cache then get from database and cache it",
			want:    expected,
			wantErr: false,
			setup: func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache) {
				mockCache.EXPECT().GetAll(gomock.Any(), "orders").Return(nil, errors.New("no cache data"))
				mockRepo.EXPECT().GetAll().Return(expected, nil)
				mockCache.EXPECT().SetAll(gomock.Any(), "orders", expected, gomock.Any()).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo, mockCache := setupTest(t)
			tt.setup(mockRepo, mockCache)
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
		setup   func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache)
	}{
		{
			name:    "input id with normal condition with cache has data then succes response no error",
			id:      1,
			want:    expected,
			wantErr: false,
			setup: func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache) {
				mockCache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(expected, nil)
				mockRepo.EXPECT().GetByID(1).Times(0)
				mockCache.EXPECT().SetByID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:    "input id with normal condition with cache no data then succes response no error and caching data",
			id:      1,
			want:    expected,
			wantErr: false,
			setup: func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache) {
				mockCache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(model.Order{}, errors.New("cache not found"))
				mockRepo.EXPECT().GetByID(1).Return(expected, nil)
				mockCache.EXPECT().SetByID(gomock.Any(), gomock.Any(), expected, gomock.Any()).Return(nil)
			},
		},
		{
			name:    "input id but redis connection failed then response still success",
			id:      1,
			want:    expected,
			wantErr: false,
			setup: func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache) {
				mockCache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(model.Order{}, errors.New("connection error"))
				mockRepo.EXPECT().GetByID(1).Return(expected, nil)
				mockCache.EXPECT().SetByID(gomock.Any(), gomock.Any(), expected, gomock.Any()).Return(errors.New("connection error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo, mockCache := setupTest(t)
			tt.setup(mockRepo, mockCache)
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
		setup   func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache)
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
			setup: func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache) {
				mockRepo.EXPECT().Create(model.Order{
					Customer: "Andi",
					Product:  "Laptop",
					Quantity: 12,
					Status:   "Pending",
				}).Return(expected, nil)
				mockCache.EXPECT().Del(gomock.Any(), "orders").Return(nil)
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
			setup: func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache) {
				mockRepo.EXPECT().Create(gomock.Any()).Times(0)
				mockCache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
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
			setup: func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache) {
				mockRepo.EXPECT().Create(gomock.Any()).Times(0)
				mockCache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
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
			setup: func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache) {
				mockRepo.EXPECT().Create(gomock.Any()).Times(0)
				mockCache.EXPECT().Del(gomock.Any(), gomock.Any()).Times(0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo, mockCache := setupTest(t)
			tt.setup(mockRepo, mockCache)
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
