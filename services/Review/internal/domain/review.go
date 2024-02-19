package domain

type Review struct {
	Id         int64  `json:"id"`
	User_id    int64  `json:"user_id"`
	Message    string `json:"message"`
	Rating     uint   `json:"rating"`
	Product_id int64  `json:"product_id"`
}
