package kudu

// #cgo CXXFLAGS: --std=c++11 -Wall
// #cgo LDFLAGS: -lkudu_client
// #include "kudu.h"
// #include <stdlib.h>
import "C"

// Table ::
type Table = C.CKuduTable

// Close ::
func (t *Table) Close() {
	C.KuduTable_Close(t)
}

// NewInsert ::
func (t *Table) NewInsert() *WriteOperation {
	return (*WriteOperation)(C.KuduTable_NewInsert(t))
}

// NewScanner ::
func (t *Table) NewScanner() *Scanner {
	return (*Scanner)(C.KuduTable_NewScanner(t))
}


