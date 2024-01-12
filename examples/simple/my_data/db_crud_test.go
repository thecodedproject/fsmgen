package my_data_test

import (
	context "context"
	fmt "fmt"
	require "github.com/stretchr/testify/require"
	simple "github.com/thecodedproject/fsmgen/examples/simple"
	my_data "github.com/thecodedproject/fsmgen/examples/simple/my_data"
	assert "github.com/thecodedproject/gotest/assert"
	time "github.com/thecodedproject/gotest/time"
	sqltest "github.com/thecodedproject/sqltest"
	testing "testing"
	time2 "time"
)

func populateDataModelFromNonce(nonce int64) simple.MyData {

	return simple.MyData{
		Field1: "some_str" + fmt.Sprint(nonce),
		Field2: "some_str" + fmt.Sprint(nonce),
		Field3: "some_str" + fmt.Sprint(nonce),
	}
}

func populateDataModelFromNonceWithIDAndTimestamp(
	nonce int64,
	id int64,
	t time2.Time,
) simple.MyData {

	d := populateDataModelFromNonce(nonce)
	d.ID = id
	return d
}

func queryFromNonce(nonce int64) map[string]any {

	return map[string]any{
		"field_1": "some_str" + fmt.Sprint(nonce),
		"field_2": "some_str" + fmt.Sprint(nonce),
		"field_3": "some_str" + fmt.Sprint(nonce),
	}
}

func TestInsertAndSelect(t *testing.T) {

	now := time.SetTimeNowForTesting(t)

	testCases := []struct{
		Name string
		ToInsert []simple.MyData
		Query map[string]any
		Expected []simple.MyData
		ExpectErr bool
	}{
		{
			Name: "selects nothing when nothing inserted",
		},
		{
			Name: "insert one and select",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(11),
			},
			Expected: []simple.MyData{
				populateDataModelFromNonceWithIDAndTimestamp(11, 1, now),
			},
		},
		{
			Name: "insert many and select",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(11),
				populateDataModelFromNonce(21),
				populateDataModelFromNonce(31),
				populateDataModelFromNonce(41),
			},
			Expected: []simple.MyData{
				populateDataModelFromNonceWithIDAndTimestamp(11, 1, now),
				populateDataModelFromNonceWithIDAndTimestamp(21, 2, now),
				populateDataModelFromNonceWithIDAndTimestamp(31, 3, now),
				populateDataModelFromNonceWithIDAndTimestamp(41, 4, now),
			},
		},
		{
			Name: "insert many and select with query",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(22),
				populateDataModelFromNonce(45),
				populateDataModelFromNonce(45),
				populateDataModelFromNonce(1),
				populateDataModelFromNonce(45),
			},
			Query: queryFromNonce(45),
			Expected: []simple.MyData{
				populateDataModelFromNonceWithIDAndTimestamp(45, 2, now),
				populateDataModelFromNonceWithIDAndTimestamp(45, 3, now),
				populateDataModelFromNonceWithIDAndTimestamp(45, 5, now),
			},
		},
		{
			Name: "select query field which is not in data model returns error",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(1),
			},
			Query: map[string]any{
				"some_field_not_in_MyData": 1,
			},
			ExpectErr: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			ctx := context.Background()
			db := sqltest.OpenMysql(t, "schema.sql")

			for _, d := range test.ToInsert {
				_, err := my_data.Insert(ctx, db, d)
				require.NoError(t, err)
			}

			actual, err := my_data.Select(ctx, db, test.Query)
			if test.ExpectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, len(test.Expected), len(actual))

			for i := range actual {
				assert.LogicallyEqual(t, test.Expected[i], actual[i], fmt.Sprint(i) + "th element not equal")
			}
		})
	}
}

func TestSelectByID(t *testing.T) {

	now := time.SetTimeNowForTesting(t)

	testCases := []struct{
		Name string
		ToInsert []simple.MyData
		ID int64
		Expected simple.MyData
		ExpectErr bool
	}{
		{
			Name: "when ID not found returns error",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(100),
				populateDataModelFromNonce(200),
			},
			ID: 12345,
			ExpectErr: true,
		},
		{
			Name: "when ID is found returns row",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(100),
				populateDataModelFromNonce(200),
				populateDataModelFromNonce(300),
			},
			ID: 2,
			Expected: populateDataModelFromNonceWithIDAndTimestamp(200, 2, now),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			ctx := context.Background()
			db := sqltest.OpenMysql(t, "schema.sql")

			for _, d := range test.ToInsert {
				_, err := my_data.Insert(ctx, db, d)
				require.NoError(t, err)
			}

			actual, err := my_data.SelectByID(ctx, db, test.ID)
			if test.ExpectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.LogicallyEqual(t, test.Expected, actual)
		})
	}
}

func TestUpdate(t *testing.T) {

	now := time.SetTimeNowForTesting(t)

	testCases := []struct{
		Name string
		ToInsert []simple.MyData
		Updates map[string]any
		Query map[string]any
		ExpectedNumUpdates int64
		Expected []simple.MyData
		ExpectErr bool
	}{
		{
			Name: "empty params does nothing",
		},
		{
			Name: "update unknown field throws error",
			Updates: map[string]any{
				"field_not_in_the_simple.MyData_type": "update",
			},
			ExpectErr: true,
		},
		{
			Name: "query unknown field throws error",
			Updates: queryFromNonce(1),
			Query: map[string]any{
				"field_not_in_MyData": "update",
			},
			ExpectErr: true,
		},
		{
			Name: "update all records",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(123),
				populateDataModelFromNonce(124),
				populateDataModelFromNonce(125),
				populateDataModelFromNonce(126),
			},
			Updates: queryFromNonce(111),
			ExpectedNumUpdates: 4,
			Expected: []simple.MyData{
				populateDataModelFromNonceWithIDAndTimestamp(111, 1, now),
				populateDataModelFromNonceWithIDAndTimestamp(111, 2, now),
				populateDataModelFromNonceWithIDAndTimestamp(111, 3, now),
				populateDataModelFromNonceWithIDAndTimestamp(111, 4, now),
			},
		},
		{
			Name: "update records with query",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(123),
				populateDataModelFromNonce(125),
				populateDataModelFromNonce(124),
				populateDataModelFromNonce(125),
				populateDataModelFromNonce(126),
				populateDataModelFromNonce(125),
			},
			Updates: queryFromNonce(999),
			Query: queryFromNonce(125),
			ExpectedNumUpdates: 3,
			Expected: []simple.MyData{
				populateDataModelFromNonceWithIDAndTimestamp(123, 1, now),
				populateDataModelFromNonceWithIDAndTimestamp(999, 2, now),
				populateDataModelFromNonceWithIDAndTimestamp(124, 3, now),
				populateDataModelFromNonceWithIDAndTimestamp(999, 4, now),
				populateDataModelFromNonceWithIDAndTimestamp(126, 5, now),
				populateDataModelFromNonceWithIDAndTimestamp(999, 6, now),
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			ctx := context.Background()
			db := sqltest.OpenMysql(t, "schema.sql")

			for _, d := range test.ToInsert {
				_, err := my_data.Insert(
					ctx,
					db,
					d,
				)
				require.NoError(t, err)
			}

			numUpdates, err := my_data.Update(
				ctx,
				db,
				test.Updates,
				test.Query,
			)
			if test.ExpectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, test.ExpectedNumUpdates, numUpdates)

			actual, err := my_data.Select(ctx, db, nil)
			require.NoError(t, err)

			require.Equal(t, len(test.Expected), len(actual))

			for i := range actual {
				assert.LogicallyEqual(t, test.Expected[i], actual[i], fmt.Sprint(i) + "th element not equal")
			}
		})
	}
}

func TestUpdateByID(t *testing.T) {

	now := time.SetTimeNowForTesting(t)

	testCases := []struct{
		Name string
		ToInsert []simple.MyData
		ID int64
		Updates map[string]any
		Expected []simple.MyData
		ExpectErr bool
	}{
		{
			Name: "no updates does not error - even if ID does not exist",
		},
		{
			Name: "when there are updates and ID not found throws error",
			ID: 1234,
			Updates: queryFromNonce(1),
			ExpectErr: true,
		},
		{
			Name: "when update field not in schema throws error",
			ID: 1,
			Updates: map[string]any{
				"field_not_in_the_simple.MyData_type": "update",
			},
			ExpectErr: true,
		},
		{
			Name: "insert many and update one by id",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(101),
				populateDataModelFromNonce(102),
				populateDataModelFromNonce(103),
				populateDataModelFromNonce(104),
			},
			ID: 3,
			Updates: queryFromNonce(555),
			Expected: []simple.MyData{
				populateDataModelFromNonceWithIDAndTimestamp(101, 1, now),
				populateDataModelFromNonceWithIDAndTimestamp(102, 2, now),
				populateDataModelFromNonceWithIDAndTimestamp(555, 3, now),
				populateDataModelFromNonceWithIDAndTimestamp(104, 4, now),
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			ctx := context.Background()
			db := sqltest.OpenMysql(t, "schema.sql")

			for _, d := range test.ToInsert {
				_, err := my_data.Insert(ctx, db, d)
				require.NoError(t, err)
			}

			err := my_data.UpdateByID(
				ctx,
				db,
				test.ID,
				test.Updates,
			)

			if test.ExpectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			actual, err := my_data.Select(ctx, db, nil)
			require.NoError(t, err)

			require.Equal(t, len(test.Expected), len(actual))

			for i := range actual {
				assert.LogicallyEqual(t, test.Expected[i], actual[i], fmt.Sprint(i) + "th element not equal")
			}
		})
	}
}

func TestDelete(t *testing.T) {

	now := time.SetTimeNowForTesting(t)

	testCases := []struct{
		Name string
		ToInsert []simple.MyData
		Query map[string]any
		ExpectedNumDeleted int64
		Expected []simple.MyData
		ExpectErr bool
	}{
		{
			Name: "empty query deletes all records",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(1000),
				populateDataModelFromNonce(1001),
				populateDataModelFromNonce(1002),
				populateDataModelFromNonce(1003),
				populateDataModelFromNonce(1004),
			},
			ExpectedNumDeleted: 5,
		},
		{
			Name: "delete records using query",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(1000),
				populateDataModelFromNonce(1001),
				populateDataModelFromNonce(1002),
				populateDataModelFromNonce(1002),
				populateDataModelFromNonce(1003),
				populateDataModelFromNonce(1004),
				populateDataModelFromNonce(1002),
			},
			Query: queryFromNonce(1002),
			ExpectedNumDeleted: 3,
			Expected: []simple.MyData{
				populateDataModelFromNonceWithIDAndTimestamp(1000, 1, now),
				populateDataModelFromNonceWithIDAndTimestamp(1001, 2, now),
				populateDataModelFromNonceWithIDAndTimestamp(1003, 5, now),
				populateDataModelFromNonceWithIDAndTimestamp(1004, 6, now),
			},
		},
		{
			Name: "when query contains field not in data model returns error",
			Query: map[string]any{
				"some_field_not_in_MyData": 1,
			},
			ExpectErr: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			ctx := context.Background()
			db := sqltest.OpenMysql(t, "schema.sql")

			for _, d := range test.ToInsert {
				_, err := my_data.Insert(ctx, db, d)
				require.NoError(t, err)
			}

			numDeleted, err := my_data.Delete(ctx, db, test.Query)
			if test.ExpectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, test.ExpectedNumDeleted, numDeleted)

			actual, err := my_data.Select(ctx, db, nil)
			require.NoError(t, err)

			require.Equal(t, len(test.Expected), len(actual))

			for i := range actual {
				assert.LogicallyEqual(t, test.Expected[i], actual[i], fmt.Sprint(i) + "th element not equal")
			}
		})
	}
}

func TestDeleteByID(t *testing.T) {

	now := time.SetTimeNowForTesting(t)

	testCases := []struct{
		Name string
		ToInsert []simple.MyData
		ID int64
		Expected []simple.MyData
		ExpectErr bool
	}{
		{
			Name: "when ID not found returns error",
			ExpectErr: true,
		},
		{
			Name: "insert many and delete by ID",
			ToInsert: []simple.MyData{
				populateDataModelFromNonce(101),
				populateDataModelFromNonce(102),
				populateDataModelFromNonce(103),
				populateDataModelFromNonce(104),
			},
			ID: 3,
			Expected: []simple.MyData{
				populateDataModelFromNonceWithIDAndTimestamp(101, 1, now),
				populateDataModelFromNonceWithIDAndTimestamp(102, 2, now),
				populateDataModelFromNonceWithIDAndTimestamp(104, 4, now),
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			ctx := context.Background()
			db := sqltest.OpenMysql(t, "schema.sql")

			for _, d := range test.ToInsert {
				_, err := my_data.Insert(ctx, db, d)
				require.NoError(t, err)
			}

			err := my_data.DeleteByID(
				ctx,
				db,
				test.ID,
			)

			if test.ExpectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			actual, err := my_data.Select(ctx, db, nil)
			require.NoError(t, err)

			require.Equal(t, len(test.Expected), len(actual))

			for i := range actual {
				assert.LogicallyEqual(t, test.Expected[i], actual[i], fmt.Sprint(i) + "th element not equal")
			}
		})
	}
}

