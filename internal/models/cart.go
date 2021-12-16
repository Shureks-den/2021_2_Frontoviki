package models

//easyjson:json
type Cart struct {
	UserId   int64 `json:"user_id" example:"1"`
	AdvertId int64 `json:"advert_id" example:"1"`
	Amount   int64 `json:"amount" example:"1"`
}

//easyjson:json
type CartHandler struct {
	AdvertId int64 `json:"advert_id" valid:"int"`
	Amount   int64 `json:"amount" valid:"optional,int"`
}

//easyjson:json
func NewCart(userId int64, cartHandler *CartHandler) *Cart {
	return &Cart{
		UserId:   userId,
		AdvertId: cartHandler.AdvertId,
		Amount:   cartHandler.Amount,
	}
}

//easyjson:json
type CartList struct {
	UserId      int64
	AdvertsCart []*CartHandler
}

//easyjson:json
type CHs []*CartHandler
