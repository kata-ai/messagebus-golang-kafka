package kudu

// #cgo CXXFLAGS: --std=c++11 -Wall
// #cgo LDFLAGS: -lkudu_client
// #include "kudu.h"
// #include <stdlib.h>
import "C"
import "unsafe"
import "errors"

// INT32 ::
const INT32 = C.KUDU_INT32
// AutoFlushBackground ::
const AutoFlushBackground = C.KUDU_AUTO_FLUSH_BACKGROUND
// Status ::
type Status = C.CKuduStatus

// StatusToError ::
func StatusToError(status *Status) error {
	if status == nil {
		return nil
	}
	defer C.KuduStatus_Free(status)
	msg := C.KuduStatus_Message(status)
	defer C.free(unsafe.Pointer(msg))
	return errors.New(C.GoString(msg))
}