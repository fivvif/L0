package repository

import (
	"L0"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type OrdersPostgres struct {
	db *sqlx.DB
}

func NewOrdersPostgres(db *sqlx.DB) *OrdersPostgres {
	return &OrdersPostgres{db: db}
}

func (s *OrdersPostgres) SaveOrder(order L0.Order) error {
	tx, err := s.db.Begin()
	insertIntoOrdersQuery := fmt.Sprintf("INSERT INTO %s (order_uid,track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", ordersTable)
	_, err = tx.Exec(insertIntoOrdersQuery, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SMID, order.DateCreated, order.OofShard)
	if err != nil {
		logrus.Errorf("Error while insert into orders : %s", err.Error())
		tx.Rollback()
		return err
	}
	insertIntoPaymentQuery := fmt.Sprintf("INSERT INTO %s (order_uid,transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", paymentTable)
	_, err = tx.Exec(insertIntoPaymentQuery, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		logrus.Errorf("Error while insert into payment : %s", err.Error())
		tx.Rollback()
		return err
	}
	insertIntoDeliveryQuery := fmt.Sprintf("INSERT INTO %s (order_uid,name,phone,zip,city,address,region,email) VALUES($1, $2, $3, $4, $5, $6, $7, $8) ", deliveryTable)
	_, err = tx.Exec(insertIntoDeliveryQuery, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		logrus.Errorf("Error while insert into delivery : %s", err.Error())
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		insertIntoItemsQuery := fmt.Sprintf("INSERT INTO %s (chrt_id,order_uid, track_number,price,rid,name,sale,size,total_price,nm_id,brand,status) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ", itemsTable)
		_, err = tx.Exec(insertIntoItemsQuery, item.ChrtID, order.OrderUID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			logrus.Errorf("Error while insert into items : %s", err.Error())
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()

}

func (s *OrdersPostgres) RecoverCache() ([]L0.Order, error) {
	var delivery L0.Delivery
	var payment L0.Payment
	var items []L0.Item
	var orders []L0.Order

	logrus.Println("hello")
	ordersQuery := fmt.Sprintf("SELECT * FROM %s", ordersTable)
	err := s.db.Select(&orders, ordersQuery)
	if err != nil {
		return nil, err
	}
	for k, order := range orders {
		UID := order.OrderUID
		paymentQuery := fmt.Sprintf("SELECT transaction,request_id,currency,provider,amount,payment_dt,bank ,delivery_cost,goods_total,custom_fee FROM %s WHERE order_uid = $1 ", paymentTable)
		err = s.db.Get(&payment, paymentQuery, UID)
		if err != nil {
			return nil, err
		}
		orders[k].Payment = payment
		deliveryQuery := fmt.Sprintf("SELECT name,phone,zip, city,address,region,email FROM %s WHERE order_uid = $1", deliveryTable)
		err := s.db.Get(&delivery, deliveryQuery, UID)
		if err != nil {
			return nil, err
		}
		orders[k].Delivery = delivery

		itemsQuery := fmt.Sprintf("SELECT chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status FROM %s WHERE order_uid = $1", itemsTable)
		err = s.db.Select(&items, itemsQuery, UID)
		if err != nil {
			return nil, err
		}
		orders[k].Items = items
	}

	return orders, nil
}
