package service_test

import (
	"context"
	"errors"
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
	tests := []struct {
		name    string
		want    []model.Order
		wantErr bool
		setup   func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache)
	}{
		{
			name: "success with caching",
			want: []model.Order{
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
			},
			wantErr: false,
			setup: func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache) {
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
				mockCache.EXPECT().GetAll(gomock.Any(), "orders").Return(expected, nil)
				mockRepo.EXPECT().GetAll().Times(0)
				mockCache.EXPECT().SetAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "success from database and cache it",
			want: []model.Order{
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
			},
			wantErr: false,
			setup: func(mockRepo *mocks.MockOrderRepository, mockCache *mocks.MockOrderCache) {
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
				mockCache.EXPECT().GetByID(gomock.Any(), 1, gomock.Any()).Return(model.Order{}, errors.New("error not found"))
				mockRepo.EXPECT().GetByID(1).Return(expected, nil)
				mockCache.EXPECT().SetByID(gomock.Any(), gomock.Any(), expected, gomock.Any()).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockOrderRepository(ctrl)
			mockCache := mocks.NewMockOrderCache(ctrl)
			tt.setup(mockRepo, mockCache)
			s := service.NewOrderService(mockRepo, mockCache)
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
