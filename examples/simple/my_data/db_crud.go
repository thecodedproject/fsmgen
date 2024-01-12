package my_data

import (
	context "context"
	sql "database/sql"
	errors "errors"
	fmt "fmt"
	simple "github.com/thecodedproject/fsmgen/examples/simple"
)

func Insert(
	ctx context.Context,
	db *sql.DB,
	d simple.MyData,
) (int64, error) {

	r, err := db.ExecContext(
		ctx,
		"insert into my_data set field_1=?, field_2=?, field_3=?",
		d.Field1,
		d.Field2,
		d.Field3,
	)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SelectByID(
	ctx context.Context,
	db *sql.DB,
	id int64,
) (simple.MyData, error) {

	r, err := Select(
		ctx,
		db,
		map[string]any{
			"id": id,
		},
	)
	if err != nil {
		return simple.MyData{}, nil
	}

	if len(r) == 0 {
		return simple.MyData{}, errors.New("SelectByID: id not found - " + fmt.Sprint(id))
	}

	if len(r) > 1 {
		return simple.MyData{}, errors.New("found more than one entry with id")
	}

	return r[0], nil
}

func Select(
	ctx context.Context,
	db *sql.DB,
	queryParams map[string]any,
) ([]simple.MyData, error) {


	q := "select id, field_1, field_2, field_3 from my_data"

	if len(queryParams) > 0 {
		q += " where "
	}

	queryVals := make([]any, 0, len(queryParams))
	i := 0
	for k, v := range queryParams {
		if !modelContainsField(k) {
			return nil, errors.New("Select: no such field to query - " + k)
		}

		q += k + "=?"
		i++
		if i < len(queryParams) {
			q += " and "
		}
		queryVals = append(queryVals, v)
	}

	r, err := db.QueryContext(
		ctx,
		q,
		queryVals...,
	)
	if err != nil {
		return nil, nil
	}

	// TODO: make this a configurable param
	maxResponses := 1000
	res := make([]simple.MyData, 0, maxResponses)
	for r.Next() {

		if len(res) >= maxResponses {
			return nil, errors.New("select query exceeded max responses")
		}

		var d simple.MyData
		r.Scan(
			&d.ID,
			&d.Field1,
			&d.Field2,
			&d.Field3,
		)
		if err != nil {
			return nil, nil
		}

		res = append(res, d)
	}

	return res, nil
}

func Update(
	ctx context.Context,
	db *sql.DB,
	updates map[string]any,
	queryParams map[string]any,
) (int64, error) {

	if len(updates) == 0 {
		return 0, nil
	}

	query := "update my_data set "
	queryArgs := make([]any, 0, len(updates) + len(queryParams))
	i := 0
	for k, v := range updates {
		if !modelContainsField(k) {
			return 0, errors.New("Update: no such field to update - " + k)
		}

		query += k + "=?"
		i++
		if i < len(updates) {
			query += ", "
		}

		queryArgs = append(queryArgs, v)
	}

	if len(queryParams) > 0 {
		query += " where "
	}
	i = 0
	for k, v := range queryParams {
		if !modelContainsField(k) {
			return 0, errors.New("Update: no such field to query - " + k)
		}

		query += k + "=?"
		i++
		if i < len(queryParams) {
			query += " and "
		}

		queryArgs = append(queryArgs, v)
	}

	r, err := db.ExecContext(
		ctx,
		query,
		queryArgs...,
	)
	if err != nil {
		return 0, err
	}

	count, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func UpdateByID(
	ctx context.Context,
	db *sql.DB,
	id int64,
	updates map[string]any,
) error {

	if len(updates) == 0 {
		return nil
	}

	n, err := Update(
		ctx,
		db,
		updates,
		map[string]any{
			"id": id,
		},
	)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("UpdateByID: no such ID")
	}

	return nil
}

func Delete(
	ctx context.Context,
	db *sql.DB,
	queryParams map[string]any,
) (int64, error) {

	query := "delete from my_data"

	if len(queryParams) > 0 {
		query += " where "
	}
	i := 0
	queryArgs := make([]any, 0, len(queryParams))
	for k, v := range queryParams {
		if !modelContainsField(k) {
			return 0, errors.New("Delete: no such field to query - " + k)
		}

		query += k + "=?"
		i++
		if i < len(queryParams) {
			query += " and "
		}

		queryArgs = append(queryArgs, v)
	}

	r, err := db.ExecContext(
		ctx,
		query,
		queryArgs...,
	)
	if err != nil {
		return 0, err
	}

	count, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func DeleteByID(
	ctx context.Context,
	db *sql.DB,
	id int64,
) error {

	n, err := Delete(
		ctx,
		db,
		map[string]any{
			"id": id,
		},
	)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("DeleteByID: no such ID")
	}

	return nil
}

func modelContainsField(field string) bool {

	modelFields := map[string]bool{
		"id": true,
		"field_1": true,
		"field_2": true,
		"field_3": true,
	}

	return modelFields[field]
}

