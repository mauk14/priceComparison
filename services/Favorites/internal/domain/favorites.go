package domain

import "time"

type Favorites struct {
	Id         int64     `json:"id"`
	User_id    int64     `json:"user_id"`
	Product_id int64     `json:"product_id"`
	Created_at time.Time `json:"created_at"`
}
