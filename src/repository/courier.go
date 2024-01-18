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
func (r *Repository) SaveCouriers(courierList []models.CreateCourierDto) ([]models.CourierDto, error) {
	var query strings.Builder
	var args []any

	ordinal := 1

	query.WriteString("insert into couriers (courier_type, regions, working_hours) values ")

	for i, courier := range courierList {
		query.WriteString(fmt.Sprintf("($%d, $%d, $%d)", ordinal, ordinal+1, ordinal+2))
		args = append(args, courier.CourierType, pq.Array(courier.Regions), pq.Array(courier.WorkingHours))
		ordinal += 3

		if i != len(courierList)-1 {
			query.WriteString(", ")
		}
	}

	query.WriteString(" returning *;")

	rows, err := r.DB.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}

	return getCourierList(rows)
}

func (r *Repository) GetCouriers(limit int32, offset int32) ([]models.CourierDto, error) {
	rows, err := r.DB.Query("select * from couriers order by courier_id desc limit $1 offset $2;", limit, offset)
	if err != nil {
		return nil, err
	}

	return getCourierList(rows)
}

func (r *Repository) GetCourierById(courierId int64) (models.CourierDto, error) {
	var courier models.CourierDto
	row := r.DB.QueryRow("select * from couriers where courier_id=$1;", courierId)
	err := row.Scan(
		&courier.CourierId,
		&courier.CourierType,
		pq.Array(&courier.Regions),
		pq.Array(&courier.WorkingHours))

	return courier, err
}

func (r *Repository) GetCourierSumAndCount(id int64, startDate string, endDate string) (int32, int32, error) {
	query := `select count(oic.order_id), sum(oic.cost) 
				from (
					select o.order_id, o.cost 
						from orders as o, completed_orders as co
						where co.courier_id = $1 
						  and o.order_id = co.order_id 
						  and (o.completed_time >= $2::timestamp and o.completed_time < $3::timestamp)
					) as oic;`

	var count sql.NullInt32
	var sum sql.NullInt32

	row := r.DB.QueryRow(query, id, startDate, endDate)

	err := row.Scan(&count, &sum)

	if err != nil {
		return 0, 0, err
	}

	return count.Int32, sum.Int32, nil
}

func getCourierList(rows *sql.Rows) ([]models.CourierDto, error) {
	courierDtoList := make([]models.CourierDto, 0)

	defer rows.Close()

	for rows.Next() {
		var courier models.CourierDto

		err := rows.Scan(
			&courier.CourierId,
			&courier.CourierType,
			pq.Array(&courier.Regions),
			pq.Array(&courier.WorkingHours))

		if err != nil {
			return nil, err
		}

		courierDtoList = append(courierDtoList, courier)
	}

	return courierDtoList, nil
}
