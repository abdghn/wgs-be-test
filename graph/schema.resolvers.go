package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/abghn/kuncie-be-test/graph/generated"
	"github.com/abghn/kuncie-be-test/graph/model"
)

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input model.OrderInput) (*model.Order, error) {
	items := r.mapItemsFromInput(input.Items)
	if len(items) == 0 {
		return nil, fmt.Errorf("order item is empty")
	}
	totalAmount := calculateTotalAmout(items)

	order := model.Order{
		CustomerName: input.CustomerName,
		OrderAmount:  totalAmount,
		Items:        items,
	}
	err := r.DB.Create(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateOrder is the resolver for the updateOrder field.
func (r *mutationResolver) UpdateOrder(ctx context.Context, orderID int, input model.OrderInput) (*model.Order, error) {
	items := r.mapItemsFromInput(input.Items)
	if len(items) == 0 {
		return nil, fmt.Errorf("order item is empty")
	}

	totalAmount := calculateTotalAmout(items)

	updatedOrder := model.Order{
		ID:           orderID,
		CustomerName: input.CustomerName,
		OrderAmount:  totalAmount,
		Items:        items,
	}
	err := r.DB.Save(&updatedOrder).Error
	if err != nil {
		return nil, err
	}
	return &updatedOrder, nil
}

// DeleteOrder is the resolver for the deleteOrder field.
func (r *mutationResolver) DeleteOrder(ctx context.Context, orderID int) (bool, error) {
	r.DB.Where("id = ?", orderID).Delete(&model.Order{})
	return true, nil
}

// Orders is the resolver for the orders field.
func (r *queryResolver) Orders(ctx context.Context) ([]*model.Order, error) {
	var orders []*model.Order
	// err := r.DB.Preload("Items").Find(&orders).Error
	err := r.DB.Set("gorm:auto_preload", true).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) mapItemsFromInput(itemsInput []*model.ItemInput) []*model.OrderItem {
	var items []*model.OrderItem
	for _, itemInput := range itemsInput {
		var item = model.Item{}
		err := r.DB.Table("items").Where("product_code = ?", itemInput.ProductCode).First(&item).Error
		if err != nil {
			fmt.Printf("error execute query %v \n", err)
			continue
		}

		switch item.Promo {
		case "BundleFreeForEveryItemBought":
			items = append(items, &model.OrderItem{
				ProductCode: itemInput.ProductCode,
				ProductName: item.ProductName,
				Quantity:    itemInput.Quantity,
				Price:       float64(itemInput.Quantity) * item.Price,
			})

			items = append(items, &model.OrderItem{
				ProductCode: "234234",
				ProductName: "Raspberry Pi B",
				Quantity:    itemInput.Quantity,
				Price:       0,
			})

		case "BuyThreePayTwoOnly":
			if itemInput.Quantity%3 == 0 {
				items = append(items, &model.OrderItem{
					ProductCode: item.ProductCode,
					ProductName: item.ProductName,
					Quantity:    itemInput.Quantity,
					Price:       float64(itemInput.Quantity-1) * item.Price,
				})
			}

		case "DiscountMoreThanThree":

			if itemInput.Quantity >= 3 {
				discPrice := item.Price * 0.1
				items = append(items, &model.OrderItem{
					ProductCode: item.ProductCode,
					ProductName: item.ProductName,
					Quantity:    itemInput.Quantity,
					Price:       float64(itemInput.Quantity) * discPrice,
				})
			}
		default:
			items = append(items, &model.OrderItem{
				ProductCode: itemInput.ProductCode,
				ProductName: item.ProductName,
				Quantity:    itemInput.Quantity,
				Price:       float64(itemInput.Quantity) * item.Price,
			})
		}
	}

	return items
}
func calculateTotalAmout(items []*model.OrderItem) float64 {
	var totalAmount float64

	for _, item := range items {
		totalAmount += item.Price
	}

	return totalAmount
}
