package kudu

// #cgo CXXFLAGS: --std=c++11 -Wall
// #cgo LDFLAGS: -lkudu_client
// #include "kudu.h"
// #include <stdlib.h>
import "C"
import "unsafe"

// WriteOperation ::
type WriteOperation = C.CKuduWriteOperation

// SetBool ::
func (op *WriteOperation) SetBool(colName string , val bool) error {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return StatusToError(C.KuduWriteOperation_SetBool(op, cName, C.bool(val)))
}	

// SetInt8 ::
func (op *WriteOperation) SetInt8(colName string , val int8) error {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return StatusToError(C.KuduWriteOperation_SetInt8(op, cName, C.int8_t(val)))
}	

// SetInt16 ::
func (op *WriteOperation) SetInt16(colName string , val int16) error {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return StatusToError(C.KuduWriteOperation_SetInt16(op, cName, C.int16_t(val)))
}	

// SetInt32 ::
func (op *WriteOperation) SetInt32(colName string , val int32) error {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return StatusToError(C.KuduWriteOperation_SetInt32(op, cName, C.int32_t(val)))
}	

// SetInt64 ::
func (op *WriteOperation) SetInt64(colName string , val int64) error {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return StatusToError(C.KuduWriteOperation_SetInt64(op, cName, C.int64_t(val)))
}	

// SetUnixTimeMicros ::
func (op *WriteOperation) SetUnixTimeMicros(colName string , val int64) error {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return StatusToError(C.KuduWriteOperation_SetUnixTimeMicros(op, cName, C.int64_t(val)))
}	

// SetDouble ::
func (op *WriteOperation) SetDouble(colName string , val float64) error {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return StatusToError(C.KuduWriteOperation_SetDouble(op, cName, C.double(val)))
}	

// SetFloat ::
func (op *WriteOperation) SetFloat(colName string , val float32) error {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return StatusToError(C.KuduWriteOperation_SetFloat(op, cName, C.float(val)))
}	

// SetString ::
func (op *WriteOperation) SetString(colName string , val string) error {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return StatusToError(C.KuduWriteOperation_SetString(op, cName, C.CString(val)))
}	

// SetStringCopy ::
func (op *WriteOperation) SetStringCopy(colName string , val string) error {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return StatusToError(C.KuduWriteOperation_SetStringCopy(op, cName, C.CString(val)))
}	

// SetNull ::
func (op *WriteOperation) SetNull(colName string) error {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return StatusToError(C.KuduWriteOperation_SetNull(op, cName))
}	

