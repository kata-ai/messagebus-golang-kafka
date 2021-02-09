package kudu

// #cgo CXXFLAGS: --std=c++11 -Wall
// #cgo LDFLAGS: -lkudu_client
// #include "kudu.h"
// #include <stdlib.h>
import "C"

// Session ::
type Session = C.CKuduSession

// Close ::
func (s *Session) Close() {
	C.KuduSession_Close(s)
}

// Apply ::
func (s *Session) Apply(op *WriteOperation) error {
	return StatusToError(C.KuduSession_Apply(s, op))
}

// SetFlushMode ::
func (s *Session) SetFlushMode(mode C.CFlushMode) error {
	return StatusToError(C.KuduSession_SetFlushMode(s, mode))
}

// Flush ::
func (s *Session) Flush() error {
	return StatusToError(C.KuduSession_Flush(s))
}

