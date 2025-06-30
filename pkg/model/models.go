package model

import "time"

type Delivery struct {
	DeliveryID int    `json:"-" db:"deliveryid"`
	FirstName  string `json:"name" db:"firstname"`
	Phone      string `json:"phone" db:"phone"`
	Zip        string `json:"zip" db:"zip"`
	City       string `json:"city" db:"city"`
	Address    string `json:"address" db:"address"`
	Region     string `json:"region" db:"region"`
	Email      string `json:"email" db:"email"`
}

type Payment struct {
	Transaction   string `json:"transaction" db:"transactionname"`
	RequestID     string `json:"request_id" db:"requestid"`
	Currency      string `json:"currency" db:"currency"`
	Provider      string `json:"provider" db:"providername"`
	Amount        int    `json:"amount" db:"amount"`
	PaymentDt     int64  `json:"payment_dt" db:"paymentdt"`
	Bank          string `json:"bank" db:"bank"`
	DeliveryCost  int    `json:"delivery_cost" db:"deliverycost"`
	GoodsTotal    int    `json:"goods_total" db:"goodstotal"`
	CustomFee     int    `json:"custom_fee" db:"customfee"`
}

type Item struct {
	ChrtID     int    `json:"chrt_id" db:"chrtid"`
	TrackNumber string `json:"track_number" db:"tracknumber"`
	Price       int    `json:"price" db:"price"`
	RID         string `json:"rid" db:"rid"`
	Name        string `json:"name" db:"itemname"`
	Sale        int    `json:"sale" db:"sale"`
	Size        string `json:"size" db:"itemsize"`
	TotalPrice  int    `json:"total_price" db:"totalprice"`
	NmID        int    `json:"nm_id" db:"nmid"`
	Brand       string `json:"brand" db:"brand"`
	Status      int    `json:"status" db:"itemstatus"`
}

type Order struct {
	OrderUID        string    `json:"order_uid" db:"orderuid"`
	TrackNumber     string    `json:"track_number" db:"tracknumber"`
	Entry           string    `json:"entry" db:"entry"`
	Locale          string    `json:"locale" db:"locale"`
	CustomerID      string    `json:"customer_id" db:"customerid"`
	DeliveryService string    `json:"delivery_service" db:"deliveryservice"`
	ShardKey        string    `json:"shardkey" db:"shardkey"`
	SmID            int       `json:"sm_id" db:"smid"`
	DateCreated     time.Time `json:"date_created" db:"datacreated"`
	OofShard        string    `json:"oof_shard" db:"oofshard"`

	Delivery *Delivery `json:"delivery"`
	Payment  *Payment  `json:"payment"`
	Items    []Item    `json:"items"`
}