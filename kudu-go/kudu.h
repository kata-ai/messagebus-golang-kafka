#pragma once

// #include <kudu/util/int128.h>
// #include <kudu/util/slice.h>
// #include "kudu/util/slice.h"

#include <stdint.h>
#include <stdbool.h>

// using namespace kudu;

#ifdef __cplusplus
extern "C" {
#endif

  typedef enum {
    KUDU_INT8 = 0, KUDU_INT16 = 1, KUDU_INT32 = 2, KUDU_INT64 = 3, KUDU_STRING = 4, KUDU_BOOL = 5, KUDU_FLOAT = 6, KUDU_DOUBLE = 7, KUDU_BINARY = 8, KUDU_UNIXTIME_MICROS = 9
  } CKuduType;

  typedef enum {
    KUDU_AUTO_FLUSH_SYNC = 0, KUDU_AUTO_FLUSH_BACKGROUND = 1, KUDU_MANUAL_FLUSH = 2
  } CFlushMode;

  typedef struct CKuduClient CKuduClient;
  typedef struct CKuduClientBuilder CKuduClientBuilder;
  typedef struct CKuduColumnSpec CKuduColumnSpec;
  typedef struct CKuduScanner CKuduScanner;
  typedef struct CKuduScanBatch CKuduScanBatch;
  typedef struct CKuduSchema CKuduSchema;
  typedef struct CKuduSchemaBuilder CKuduSchemaBuilder;
  typedef struct CKuduSession CKuduSession;
  typedef struct CKuduStatus CKuduStatus;
  typedef struct CKuduTable CKuduTable;
  typedef struct CKuduTableCreator CKuduTableCreator;
  typedef struct CKuduWriteOperation CKuduWriteOperation;

  char* KuduStatus_Message(CKuduStatus* self);
  void KuduStatus_Free(CKuduStatus* self);

  // ------------------------------------------------------------
  CKuduClientBuilder* KuduClientBuilder_Create();
  void KuduClientBuilder_Free(CKuduClientBuilder* self);
  void KuduClientBuilder_add_master_server_addr(CKuduClientBuilder* self, const char* addr);
  CKuduStatus* KuduClientBuilder_Build(CKuduClientBuilder* self, CKuduClient** client);

  // ------------------------------------------------------------
  CKuduTableCreator* KuduClient_NewTableCreator(CKuduClient* self);
  void KuduClient_Free(CKuduClient* self);

  CKuduStatus* KuduClient_TableExists(CKuduClient* self, const char* table_name, int* exists);

  CKuduStatus* KuduClient_OpenTable(CKuduClient* self, const char* table_name, CKuduTable** table);

  CKuduSession* KuduClient_NewSession(CKuduClient* self);

  // ------------------------------------------------------------
  void KuduSession_Close(CKuduSession* self);
  CKuduStatus* KuduSession_SetFlushMode(CKuduSession* self, CFlushMode mode);

  CKuduStatus* KuduSession_Apply(CKuduSession* self, CKuduWriteOperation* op);

  CKuduStatus* KuduSession_Flush(CKuduSession* self);

  // ------------------------------------------------------------
  void KuduTable_Close(CKuduTable* table);
  CKuduWriteOperation* KuduTable_NewInsert(CKuduTable* self);

  // ------------------------------------------------------------
  // // Setters
  // Slice
  CKuduStatus* KuduWriteOperation_SetBool(CKuduWriteOperation* op, const char* col_name, bool val);
  CKuduStatus* KuduWriteOperation_SetInt8(CKuduWriteOperation* op, const char* col_name, int8_t val);
  CKuduStatus* KuduWriteOperation_SetInt16(CKuduWriteOperation* op, const char* col_name, int16_t val);
  CKuduStatus* KuduWriteOperation_SetInt32(CKuduWriteOperation* op, const char* col_name, int32_t val);
  CKuduStatus* KuduWriteOperation_SetInt64(CKuduWriteOperation* op, const char* col_name, int64_t val);
  CKuduStatus* KuduWriteOperation_SetUnixTimeMicros(CKuduWriteOperation* op, const char* col_name, int64_t micros_since_utc_epoch);
  CKuduStatus* KuduWriteOperation_SetDouble(CKuduWriteOperation* op, const char* col_name, double val);
  CKuduStatus* KuduWriteOperation_SetFloat(CKuduWriteOperation* op, const char* col_name, float val);
  // // Integer
  CKuduStatus* KuduWriteOperation_SetInt8_WithIndex(CKuduWriteOperation* op, int col_idx, int8_t val);
  CKuduStatus* KuduWriteOperation_SetInt16_WithIndex(CKuduWriteOperation* op, int col_idx, int16_t val);
  CKuduStatus* KuduWriteOperation_SetInt32_WithIndex(CKuduWriteOperation* op, int col_idx, int32_t val);
  CKuduStatus* KuduWriteOperation_SetInt64_WithIndex(CKuduWriteOperation* op, int col_idx, int64_t val);
  CKuduStatus* KuduWriteOperation_SetUnixTimeMicros_WithIndex(CKuduWriteOperation* op, int col_idx, int64_t micros_since_utc_epoch);
  CKuduStatus* KuduWriteOperation_SetDouble_WithIndex(CKuduWriteOperation* op, int col_idx, double val);
  CKuduStatus* KuduWriteOperation_SetFloat_WithIndex(CKuduWriteOperation* op, int col_idx, float val);
  // // Without copy string
  CKuduStatus* KuduWriteOperation_SetString(CKuduWriteOperation* op, const char* col_name, char* val);
  CKuduStatus* KuduWriteOperation_SetString_WithIndex(CKuduWriteOperation* op, int col_idx, char* val);  
  CKuduStatus* KuduWriteOperation_SetStringCopy(CKuduWriteOperation* op, const char* col_name, char* val);
  CKuduStatus* KuduWriteOperation_SetStringCopy_WithIndex(CKuduWriteOperation* op, int col_idx, char* val);
  CKuduStatus* KuduWriteOperation_SetNull(CKuduWriteOperation* op, const char* col_name);
  CKuduStatus* KuduWriteOperation_SetNull_WithIndex(CKuduWriteOperation* op, int col_idx);
  CKuduStatus* KuduWriteOperation_Unset(CKuduWriteOperation* op, const char* col_name);
  CKuduStatus* KuduWriteOperation_Unset_WithIndex(CKuduWriteOperation* op, int col_idx);

  CKuduScanner* KuduTable_NewScanner(CKuduTable* self);

  // ------------------------------------------------------------
  void KuduScanner_Free(CKuduScanner* scanner);
  CKuduStatus* KuduScanner_SetProjectedColumns(CKuduScanner* self, const char** col_names, int n_cols);
  CKuduStatus* KuduScanner_Open(CKuduScanner* self);
  int KuduScanner_HasMoreRows(CKuduScanner* self);
  CKuduStatus* KuduScanner_NextBatch(CKuduScanner* self, CKuduScanBatch** batch);

  // ------------------------------------------------------------
  void KuduScanBatch_Free(CKuduScanBatch* self);

  int KuduScanBatch_HasNext(CKuduScanBatch* self);
  void KuduScanBatch_SeekNext(CKuduScanBatch* self);
  // allocates result (caller must free)
  const char* KuduScanBatch_Row_ToString(CKuduScanBatch* self);

  // ------------------------------------------------------------
  CKuduSchemaBuilder* KuduSchemaBuilder_Create();
  void KuduSchemaBuilder_Free(CKuduSchemaBuilder* self);
  CKuduColumnSpec* KuduSchemaBuilder_AddColumn(CKuduSchemaBuilder* self, const char* name);
  void KuduSchemaBuilder_SetPrimaryKey(CKuduSchemaBuilder* self, const char** col_names, int n_cols);
  CKuduStatus* KuduSchemaBuilder_Build(CKuduSchemaBuilder* self, CKuduSchema** out_schema);

  // ------------------------------------------------------------
  CKuduColumnSpec* KuduColumnSpec_SetType(CKuduColumnSpec* self, int type);
  CKuduColumnSpec* KuduColumnSpec_SetNotNull(CKuduColumnSpec* self);
  
  // ------------------------------------------------------------
  void KuduTableCreator_Free(CKuduTableCreator* self);

  void KuduTableCreator_SetTableName(CKuduTableCreator* self, const char* name);

  void KuduTableCreator_SetSchema(CKuduTableCreator* self, CKuduSchema* schema);

  void KuduTableCreator_AddHashPartitions(CKuduTableCreator* self, const char** col_names, int n_cols, int num_buckets);

  CKuduStatus* KuduTableCreator_Create(CKuduTableCreator* self);

#ifdef __cplusplus
}
#endif