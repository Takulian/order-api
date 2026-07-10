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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockOrderRepository(ctrl)
			mockCache := mocks.NewMockOrderCache(ctrl)
			tt.setup(mockRepo, mockCache)
			s := service.NewOrderService(mockRepo, mockCache)
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
