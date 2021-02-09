package kudu

// #cgo CXXFLAGS: --std=c++11 -Wall
// #cgo LDFLAGS: -lkudu_client
// #include "kudu.h"
// #include <stdlib.h>
import "C"
import "unsafe"

// ScanBatch ::
type ScanBatch = C.CKuduScanBatch

// Free ::
func (b *ScanBatch) Free() {
	C.KuduScanBatch_Free(b)
}
// Next ::
func (b *ScanBatch) Next() bool {
	if C.KuduScanBatch_HasNext(b) == 0 {
		return false
	}
	C.KuduScanBatch_SeekNext(b)
	return true
}

// RowToString ::
func (b *ScanBatch) RowToString() string {
	cStr := C.KuduScanBatch_Row_ToString(b)
	defer C.free(unsafe.Pointer(cStr))
	return C.GoString(cStr)
}