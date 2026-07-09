package service_test

import (
	"errors"
	"order-api/internal/dto"
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
		setup   func(*mocks.MockOrderRepository)
	}{
		{
			name: "success",
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
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.EXPECT().GetAll().Return([]model.Order{
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
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockOrderRepository(ctrl)
			tt.setup(mockRepo)
			s := service.NewOrderService(mockRepo)
			got, gotErr := s.GetAll()
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
	tests := []struct {
		name    string
		id      int
		want    model.Order
		wantErr bool
		setup   func(*mocks.MockOrderRepository)
	}{
		{
			name: "success",
			id:   1,
			want: model.Order{
				ID:       1,
				Customer: "Andi",
				Product:  "Laptop",
				Quantity: 2,
				Status:   "Pending",
			},
			wantErr: false,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.
					EXPECT().GetByID(1).Return(
					model.Order{
						ID:       1,
						Customer: "Andi",
						Product:  "Laptop",
						Quantity: 2,
						Status:   "Pending",
					}, nil)
			},
		},
		{
			name:    "order not found",
			id:      67,
			want:    model.Order{},
			wantErr: true,
			setup: func(mockrepo *mocks.MockOrderRepository) {
				mockrepo.
					EXPECT().
					GetByID(67).
					Return(model.Order{}, errors.New("order not found"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockOrderRepository(ctrl)

			tt.setup(mockRepo)

			s := service.NewOrderService(mockRepo)

			got, gotErr := s.GetByID(tt.id)
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
	tests := []struct {
		name    string
		req     dto.CreateOrderRequest
		want    model.Order
		wantErr error
		setup   func(*mocks.MockOrderRepository)
	}{
		{
			name: "success",
			req: dto.CreateOrderRequest{
				Customer: "Andi",
				Product:  "Laptop",
				Quantity: 1,
			},
			want: model.Order{
				ID:       1,
				Customer: "Andi",
				Product:  "Laptop",
				Quantity: 1,
				Status:   "Pending",
			},
			wantErr: nil,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.EXPECT().Create(model.Order{
					Customer: "Andi",
					Product:  "Laptop",
					Quantity: 1,
					Status:   "Pending",
				}).Return(model.Order{
					ID:       1,
					Customer: "Andi",
					Product:  "Laptop",
					Quantity: 1,
					Status:   "Pending",
				}, nil)
			},
		},
		{
			name: "customer empty",
			req: dto.CreateOrderRequest{
				Customer: "",
				Product:  "Laptop",
				Quantity: 1,
			},
			want:    model.Order{},
			wantErr: service.ErrCustomerRequired,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.EXPECT().Create(gomock.Any()).Times(0)
			},
		},
		{
			name: "product empty",
			req: dto.CreateOrderRequest{
				Customer: "Andi",
				Product:  "",
				Quantity: 1,
			},
			want:    model.Order{},
			wantErr: service.ErrProductRequired,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.EXPECT().Create(gomock.Any()).Times(0)
			},
		},
		{
			name: "invalid quantity",
			req: dto.CreateOrderRequest{
				Customer: "Andi",
				Product:  "Laptop",
				Quantity: 0,
			},
			want:    model.Order{},
			wantErr: service.ErrInvalidQuantity,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.EXPECT().Create(gomock.Any()).Times(0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockOrderRepository(ctrl)
			tt.setup(mockRepo)
			s := service.NewOrderService(mockRepo)
			got, gotErr := s.Create(tt.req)

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
	tests := []struct {
		name    string
		id      int
		req     dto.UpdateOrderRequest
		want    model.Order
		wantErr error
		setup   func(*mocks.MockOrderRepository)
	}{
		{
			name: "success",
			id:   1,
			req: dto.UpdateOrderRequest{
				Customer: "Budi",
				Product:  "Handphone",
				Quantity: 10,
			},
			want: model.Order{
				ID:       1,
				Customer: "Budi",
				Product:  "Handphone",
				Quantity: 10,
				Status:   "Pending",
			},
			wantErr: nil,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				oldOrder := model.Order{
					ID:       1,
					Customer: "Andi",
					Product:  "Laptop",
					Quantity: 2,
					Status:   "Pending",
				}

				mockRepo.EXPECT().GetByID(1).Return(oldOrder, nil)
				mockRepo.EXPECT().Update(1, model.Order{
					ID:       1,
					Customer: "Budi",
					Product:  "Handphone",
					Quantity: 10,
					Status:   "Pending",
				}).Return(model.Order{
					ID:       1,
					Customer: "Budi",
					Product:  "Handphone",
					Quantity: 10,
					Status:   "Pending",
				}, nil)
			},
		},
		{
			name: "order not found",
			id:   67,
			req: dto.UpdateOrderRequest{
				Customer: "Budi",
				Product:  "Handphone",
				Quantity: 10,
			},
			want:    model.Order{},
			wantErr: service.ErrOrderNotFound,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.EXPECT().GetByID(67).Return(model.Order{}, service.ErrOrderNotFound)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "customer empty",
			id:   1,
			req: dto.UpdateOrderRequest{
				Customer: "",
				Product:  "Handphone",
				Quantity: 10,
			},
			want:    model.Order{},
			wantErr: service.ErrCustomerRequired,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.EXPECT().GetByID(gomock.Any()).Times(0)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "product empty",
			id:   1,
			req: dto.UpdateOrderRequest{
				Customer: "Budi",
				Product:  "",
				Quantity: 10,
			},
			want:    model.Order{},
			wantErr: service.ErrProductRequired,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.EXPECT().GetByID(gomock.Any()).Times(0)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "invalid quantity",
			id:   1,
			req: dto.UpdateOrderRequest{
				Customer: "Budi",
				Product:  "Handphone",
				Quantity: -2,
			},
			want:    model.Order{},
			wantErr: service.ErrInvalidQuantity,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.EXPECT().GetByID(gomock.Any()).Times(0)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockOrderRepository(ctrl)
			tt.setup(mockRepo)
			s := service.NewOrderService(mockRepo)
			got, gotErr := s.Update(tt.id, tt.req)
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
	tests := []struct {
		name    string
		id      int
		wantErr bool
		setup   func(*mocks.MockOrderRepository)
	}{
		{
			name:    "success",
			id:      1,
			wantErr: false,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.EXPECT().Delete(1).Return(nil)
			},
		},
		{
			name:    "order not found",
			id:      67,
			wantErr: true,
			setup: func(mockRepo *mocks.MockOrderRepository) {
				mockRepo.EXPECT().Delete(67).Return(errors.New("error not found"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockOrderRepository(ctrl)
			tt.setup(mockRepo)
			s := service.NewOrderService(mockRepo)
			gotErr := s.Delete(tt.id)
			if (gotErr != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", gotErr, tt.wantErr)
			}
		})
	}
}
