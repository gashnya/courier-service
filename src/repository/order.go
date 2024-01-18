package repository

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"yandex-team.ru/bstask/models"
)

/*
 * There's a limit on a query length, but we don't consider it on the current prototyping stage.
 * The query could be splitted into multiple inserts in one transaction.
 */
func (r *Repository) SaveOrders(orderList []models.CreateOrderDto) ([]models.OrderDto, error) {
	var query strings.Builder
	var args []any

	ordinal := 1

	query.WriteString("insert into orders (weight, region, cost, delivery_hours) values ")

	for i, order := range orderList {
		query.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d)", ordinal, ordinal+1, ordinal+2, ordinal+3))
		args = append(args, order.Weight, order.Regions, order.Cost, pq.Array(order.DeliveryHours))
		ordinal += 4

		if i != len(orderList)-1 {
			query.WriteString(", ")
		}
	}

	query.WriteString(" returning *;")

	rows, err := r.DB.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}

	return getOrderList(rows)
}

func (r *Repository) GetOrders(limit int32, offset int32) ([]models.OrderDto, error) {
	rows, err := r.DB.Query("select * from orders order by order_id desc limit $1 offset $2;", limit, offset)
	if err != nil {
		return nil, err
	}

	return getOrderList(rows)
}

func (r *Repository) GetOrderById(id int64) (models.OrderDto, error) {
	var order models.OrderDto

	row := r.DB.QueryRow("select * from orders where order_id=$1;", id)
	err := row.Scan(
		&order.OrderId,
		&order.Weight,
		&order.Regions,
		pq.Array(&order.DeliveryHours),
		&order.Cost,
		&order.CompletedTime)

	return order, err
}

func (r *Repository) CompleteOrders(orderList []models.CompleteOrderDto) ([]models.OrderDto, error) {
	var queryInsert, queryUpdate strings.Builder
	var argsInsert, argsUpdate []any

	// insert

	ordinal := 1

	queryInsert.WriteString("insert into completed_orders (courier_id, order_id) values ")

	for i, order := range orderList {
		queryInsert.WriteString(fmt.Sprintf("($%d, $%d)", ordinal, ordinal+1))
		argsInsert = append(argsInsert, order.CourierId, order.OrderId)
		ordinal += 2

		if i != len(orderList)-1 {
			queryInsert.WriteString(", ")
		}
	}

	queryInsert.WriteString("; ")

	// update

	ordinal = 1

	queryUpdate.WriteString("update orders set completed_time = v.completed_time from (values ")

	for i, order := range orderList {
		queryUpdate.WriteString(fmt.Sprintf("($%d::integer, $%d::timestamp)", ordinal, ordinal+1))
		argsUpdate = append(argsUpdate, order.OrderId, order.CompleteTime)
		ordinal += 2

		if i != len(orderList)-1 {
			queryUpdate.WriteString(", ")
		}
	}

	queryUpdate.WriteString(") as v (order_id, completed_time) where v.order_id = orders.order_id returning orders.*;")

	// transaction

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	_, err = tx.Exec(queryInsert.String(), argsInsert...)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(queryUpdate.String(), argsUpdate...)
	if err != nil {
		return nil, err
	}

	result, err := getOrderList(rows)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getOrderList(rows *sql.Rows) ([]models.OrderDto, error) {
	orderDtoList := make([]models.OrderDto, 0)

	defer rows.Close()

	for rows.Next() {
		var order models.OrderDto

		err := rows.Scan(
			&order.OrderId,
			&order.Weight,
			&order.Regions,
			pq.Array(&order.DeliveryHours),
			&order.Cost,
			&order.CompletedTime)

		if err != nil {
			return nil, err
		}

		orderDtoList = append(orderDtoList, order)
	}

	return orderDtoList, nil
}
