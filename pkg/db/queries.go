package db

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
	"orders/pkg/model"
)


func InsertOrder(db *sql.DB, order model.Order) error {
	tx, err := db.Begin()

	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	defer tx.Rollback()

	// Inserting Delivery
	var deliveryID int
	err = tx.QueryRow(`
		INSERT INTO delivery (FirstName, Phone, Zip, City, Address, Region, Email)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING DeliveryId`,
		order.Delivery.FirstName,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	).Scan(&deliveryID)

	if err != nil {
		return fmt.Errorf("inserting delivery: %w", err)
	}

	// Inserting Payment
	_, err = tx.Exec(`
		INSERT INTO payment (
			TransactionName, RequestId, Currency, ProviderName,
			Amount, PaymentDt, Bank, DeliveryCost, GoodsTotal, CustomFee)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)

	if err != nil {
		return fmt.Errorf("inserting payment: %w", err)
	}

	// Inserting Order
	_, err = tx.Exec(`
		INSERT INTO orders (
			OrderUid, TrackNumber, Entry, Locale, CustomerId,
			DeliveryService, ShardKey, SmId, DataCreated, OofShard,
			DeliveryId, PaymentId
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
		deliveryID,
		order.Payment.Transaction,
	)

	if err != nil {
		return fmt.Errorf("inserting order: %w", err)
	}

	// Inserting Items and linking them
	for _, item := range order.Items {
		var itemID int
		err = tx.QueryRow(`
			INSERT INTO item (
				TrackNumber, Price, RId, ItemName,
				Sale, ItemSize, TotalPrice, NmId, Brand, ItemStatus
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
			RETURNING ChrtId`,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		).Scan(&itemID)

		if err != nil {
			return fmt.Errorf("inserting item: %w", err)
		}

		_, err = tx.Exec(`
			INSERT INTO items (ItemId, OrderId) VALUES ($1, $2)`,
			itemID, order.OrderUID,
		)

		if err != nil {
			return fmt.Errorf("linking item to order: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	
	return nil

}

func GetOrderById(db *sql.DB, orderUID string) (*model.Order, error) {
	var order model.Order
	var deliveryID int
	var paymentID string

	// Order
	err := db.QueryRow(`
		SELECT
			OrderUid, TrackNumber, Entry, Locale, CustomerId,
			DeliveryService, ShardKey, SmId, DataCreated, OofShard,
			DeliveryId, PaymentId
		FROM orders
		WHERE OrderUid = $1`, orderUID).
		Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.CustomerID, &order.DeliveryService, &order.ShardKey,
			&order.SmID, &order.DateCreated, &order.OofShard,
			&deliveryID, &paymentID,
		)
	
	if err != nil {
		return nil, fmt.Errorf("getting order: %w", err)
	}

	// Delivery
	order.Delivery = &model.Delivery{}
	err = db.QueryRow(`
		SELECT FirstName, Phone, Zip, City, Address, Region, Email
		FROM delivery WHERE DeliveryId = $1`, deliveryID).
		Scan(
			&order.Delivery.FirstName, &order.Delivery.Phone,
			&order.Delivery.Zip, &order.Delivery.City,
			&order.Delivery.Address, &order.Delivery.Region,
			&order.Delivery.Email,
		)
	
	if err != nil {
		return nil, fmt.Errorf("getting delivery: %w", err)
	}
	
	// Payment
	order.Payment = &model.Payment{}
	err = db.QueryRow(`
		SELECT TransactionName, RequestId, Currency, ProviderName,
		       Amount, PaymentDt, Bank, DeliveryCost, GoodsTotal, CustomFee
		FROM payment WHERE TransactionName = $1`, paymentID).
		Scan(
			&order.Payment.Transaction, &order.Payment.RequestID,
			&order.Payment.Currency, &order.Payment.Provider,
			&order.Payment.Amount, &order.Payment.PaymentDt,
			&order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee,
		)
	
	if err != nil {
		return nil, fmt.Errorf("getting payment: %w", err)
	}

	// Items
	rows, err := db.Query(`
		SELECT i.ChrtId, i.TrackNumber, i.Price, i.RId, i.ItemName,
		       i.Sale, i.ItemSize, i.TotalPrice, i.NmId, i.Brand, i.ItemStatus
		FROM item i
		JOIN items it ON i.ChrtId = it.ItemId
		WHERE it.OrderId = $1`, orderUID)

	if err != nil {
		return nil, fmt.Errorf("getting items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price,
			&item.RID, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

func GetAllOrders(db *sql.DB) ([]*model.Order, error) {
	query := `
		SELECT 
			o.orderuid, o.tracknumber, o.entry, o.locale, o.customerid,
			o.deliveryservice, o.shardkey, o.smid, o.datacreated, o.oofshard,

			d.deliveryid, d.firstname, d.phone, d.zip, d.city, d.address, d.region, d.email,

			p.transactionname, p.requestid, p.currency, p.providername, p.amount, 
			p.paymentdt, p.bank, p.deliverycost, p.goodstotal, p.customfee

		FROM orders o
		JOIN delivery d ON o.deliveryid = d.deliveryid
		JOIN payment p ON o.paymentid = p.transactionname;
		`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query orders failed: %w", err)
	}
	defer rows.Close()

	var orders []*model.Order

	for rows.Next() {
		var o model.Order
		var d model.Delivery
		var p model.Payment

		err := rows.Scan(
			&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.CustomerID,
			&o.DeliveryService, &o.ShardKey, &o.SmID, &o.DateCreated, &o.OofShard,

			&d.DeliveryID, &d.FirstName, &d.Phone, &d.Zip, &d.City,
			&d.Address, &d.Region, &d.Email,

			&p.Transaction, &p.RequestID, &p.Currency, &p.Provider,
			&p.Amount, &p.PaymentDt, &p.Bank, &p.DeliveryCost, &p.GoodsTotal, &p.CustomFee,
		)
		if err != nil {
			return nil, fmt.Errorf("scan order row failed: %w", err)
		}

		items, err := getItemsForOrder(db, o.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("fetching items for order %s failed: %w", o.OrderUID, err)
		}

		o.Delivery = &d
		o.Payment = &p
		o.Items = items

		orders = append(orders, &o)
	}

	return orders, nil
}


func getItemsForOrder(db *sql.DB, orderUID string) ([]model.Item, error) {
	query := `
	SELECT 
		i.chrtid, i.tracknumber, i.price, i.rid, i.itemname, i.sale,
		i.itemsize, i.totalprice, i.nmid, i.brand, i.itemstatus
	FROM items it
	JOIN item i ON it.itemid = i.chrtid
	WHERE it.orderid = $1;
	`

	rows, err := db.Query(query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("query items failed: %w", err)
	}
	defer rows.Close()

	var items []model.Item

	for rows.Next() {
		var it model.Item
		if err := rows.Scan(
			&it.ChrtID, &it.TrackNumber, &it.Price, &it.RID, &it.Name, &it.Sale,
			&it.Size, &it.TotalPrice, &it.NmID, &it.Brand, &it.Status,
		); err != nil {
			return nil, fmt.Errorf("scan item row failed: %w", err)
		}
		items = append(items, it)
	}

	return items, nil
}