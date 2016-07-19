package models

import "github.com/raggaer/pigo"

// Category holds information about shop categories
type Category struct {
	ID          int64
	Name        string
	Description string
	Active      bool
}

// Offer holds information about a shop offer
type Offer struct {
	ItemID      int64
	Price       int
	Name        string
	Description string
}

// GetCategories returns all shop categories
func GetCategories() ([]*Category, error) {
	rows, err := pigo.Database.Query("SELECT id, name, description FROM cloaka_shop_categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []*Category{}
	for rows.Next() {
		category := &Category{}
		rows.Scan(&category.ID, &category.Name, &category.Description)
		categories = append(categories, category)
	}
	return categories, nil
}

// GetOffers returns all the offers related to a category
func (cat *Category) GetOffers() ([]*Offer, error) {
	rows, err := pigo.Database.Query("SELECT item_id, price, name, description FROM cloaka_shop_items WHERE category_id = ?", cat.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Offer{}
	for rows.Next() {
		item := &Offer{}
		rows.Scan(&item.ItemID, &item.Price, &item.Name, &item.Description)
		items = append(items, item)
	}
	return items, nil
}
