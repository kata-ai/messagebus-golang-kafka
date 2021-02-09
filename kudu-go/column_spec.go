package kudu

// #cgo CXXFLAGS: --std=c++11 -Wall
// #cgo LDFLAGS: -lkudu_client
// #include "kudu.h"
// #include <stdlib.h>
import "C"

// ColumnSpec :: 
type ColumnSpec = C.CKuduColumnSpec

// SetType ::
func (cspec *ColumnSpec) SetType(colType int) {
	C.KuduColumnSpec_SetType(cspec, C.int(colType))
}

// SetNotNull ::
func (cspec *ColumnSpec) SetNotNull() {
	C.KuduColumnSpec_SetNotNull(cspec)
}