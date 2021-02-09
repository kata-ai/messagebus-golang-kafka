package kudu

// #cgo CXXFLAGS: --std=c++11 -Wall
// #cgo LDFLAGS: -lkudu_client
// #include "kudu.h"
// #include <stdlib.h>
import "C"
import 	"unsafe"

// Client ::
type Client = C.CKuduClient

// Close ::
func (c *Client) Close() {
	C.KuduClient_Free(c)
}

// NewTableCreator ::
func (c *Client) NewTableCreator() *TableCreator {
	return (*TableCreator)(C.KuduClient_NewTableCreator(c))
}

// TableExists ::
func (c *Client) TableExists(name string) (bool, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var exists C.int
	if err := StatusToError(C.KuduClient_TableExists(c, cName, &exists)); err != nil {
		return false, err
	}
	return exists != 0, nil
}

// OpenTable ::
func (c *Client) OpenTable(name string) (*Table, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var table *Table
	if err := StatusToError(C.KuduClient_OpenTable(c, cName, &table)); err != nil {
		return nil, err
	}
	return (*Table)(table), nil
}

// NewSession ::
func (c *Client) NewSession() *Session {
	return (*Session)(C.KuduClient_NewSession(c))
}
