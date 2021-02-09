package kudu

// #cgo CXXFLAGS: --std=c++11 -Wall
// #cgo LDFLAGS: -lkudu_client
// #include "kudu.h"
// #include <stdlib.h>
import "C"
import "unsafe"

// ClientBuilder ::
type ClientBuilder = C.CKuduClientBuilder


// NewClientBuilder ::
func NewClientBuilder() *ClientBuilder {
	return (*ClientBuilder)(C.KuduClientBuilder_Create())
}

// Free ::
func (b *ClientBuilder) Free() {
	C.KuduClientBuilder_Free(b)
}

// AddMasterServerAddr ::
func (b *ClientBuilder) AddMasterServerAddr(addr string) {
	cAddr := C.CString(addr)
	defer C.free(unsafe.Pointer(cAddr))
	C.KuduClientBuilder_add_master_server_addr(b, cAddr)
}

// Build ::
func (b *ClientBuilder) Build() (*Client, error) {
	var client *Client
	if err := StatusToError(C.KuduClientBuilder_Build(b, &client)); err != nil {
		return nil, err
	}
	return (*Client)(client), nil
}
