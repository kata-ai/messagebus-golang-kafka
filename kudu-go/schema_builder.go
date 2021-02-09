package kudu

// #cgo CXXFLAGS: --std=c++11 -Wall
// #cgo LDFLAGS: -lkudu_client
// #include "kudu.h"
// #include <stdlib.h>
import "C"
import "unsafe"

// Schema ::
type Schema = C.CKuduSchema
// SchemaBuilder :: 
type SchemaBuilder  = C.CKuduSchemaBuilder

// NewSchemaBuilder ::
func NewSchemaBuilder() *SchemaBuilder {
	return (*SchemaBuilder)(C.KuduSchemaBuilder_Create())
}

// Free ::
func (sb *SchemaBuilder) Free() {
	C.KuduSchemaBuilder_Free(sb)
}

// SetPrimaryKey ::
func (sb *SchemaBuilder) SetPrimaryKey(colNames []string) {
	cNames := make([]*C.char, len(colNames))
	for i, s := range colNames {
		cNames[i] = C.CString(s)
		defer C.free(unsafe.Pointer(cNames[i]))
	}
	C.KuduSchemaBuilder_SetPrimaryKey(sb, &cNames[0],
		C.int(len(cNames)))
}

// AddColumn ::
func (sb *SchemaBuilder) AddColumn(colName string) *ColumnSpec {
	cName := C.CString(colName)
	defer C.free(unsafe.Pointer(cName))
	return (*ColumnSpec)(C.KuduSchemaBuilder_AddColumn(sb, cName))
}

// Build ::
func (sb *SchemaBuilder) Build() (*Schema, error) {
	var schema *Schema
	if err := StatusToError(C.KuduSchemaBuilder_Build(sb, &schema)); err != nil {
		return nil, err
	}
	return (*Schema)(schema), nil
}