package kudu

// #cgo CXXFLAGS: --std=c++11 -Wall
// #cgo LDFLAGS: -lkudu_client
// #include "kudu.h"
// #include <stdlib.h>
import "C"
import "unsafe"

// TableCreator ::
type TableCreator = C.CKuduTableCreator

// Free ::
func (tc *TableCreator) Free() {
	C.KuduTableCreator_Free(tc)
}

// SetSchema ::
func (tc *TableCreator) SetSchema(schema *Schema) {
	C.KuduTableCreator_SetSchema(tc, schema)
}

// SetTableName ::
func (tc *TableCreator) SetTableName(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.KuduTableCreator_SetTableName(tc, cName)
}

// Create ::
func (tc *TableCreator) Create() error {
	return StatusToError(C.KuduTableCreator_Create(tc))
}

// AddHashPartitions ::
func (tc *TableCreator) AddHashPartitions(colNames []string, numBuckets int) {
	cNames := make([]*C.char, len(colNames))
	for i, s := range colNames {
		cNames[i] = C.CString(s)
		defer C.free(unsafe.Pointer(cNames[i]))
	}
	C.KuduTableCreator_AddHashPartitions(tc, &cNames[0], C.int(len(cNames)), C.int(numBuckets))
}