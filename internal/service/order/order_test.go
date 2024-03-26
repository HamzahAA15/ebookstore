package order_test

// import (
// 	"context"
// 	"ebookstore/internal/model/response"
// 	"ebookstore/internal/repository/mocks"
// 	"ebookstore/internal/service/order"
// 	"ebookstore/utils/transactioner"
// 	"reflect"
// 	"testing"
// )

// func Test_orderService_GetUserOrders(t *testing.T) {
// 	type fields struct {
// 		orderRepository     mocks.IOrderRepository
// 		TransactionProvider transactioner.TransactionProvider
// 		bookRepository      mocks.IBookRepository
// 	}
// 	type args struct {
// 		ctx context.Context
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    []response.OrderData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			o := order.NewOrderService()
// 			got, err := o.GetUserOrders(tt.args.ctx)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("orderService.GetUserOrders() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("orderService.GetUserOrders() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
