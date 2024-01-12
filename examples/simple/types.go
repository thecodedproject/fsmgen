package simple

import (

	"github.com/thecodedproject/dbcrudgen/dbcrudgen"
)

type MyData struct {
	dbcrudgen.DataModel

	ID int64

	//State 

	Field1 string
	Field2 string
	Field3 string
}
