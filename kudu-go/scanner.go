
package kudu

// #cgo CXXFLAGS: --std=c++11 -Wall
// #cgo LDFLAGS: -lkudu_client
// #include "kudu.h"
// #include <stdlib.h>
import "C"

// Scanner ::
type Scanner = C.CKuduScanner

// Open ::
func (s *Scanner) Open() error {
	return StatusToError(C.KuduScanner_Open(s))
}

// Free ::
func (s *Scanner) Free() {
	C.KuduScanner_Free(s)
}

// NextBatch ::
func (s *Scanner) NextBatch(batch **ScanBatch) error {
	var cBatch *ScanBatch
	err := StatusToError(C.KuduScanner_NextBatch(s, &cBatch))
	*batch = (*ScanBatch)(cBatch)
	return err
}

// HasMoreRows ::
func (s *Scanner) HasMoreRows() bool {
	return C.KuduScanner_HasMoreRows(s) != 0
}
