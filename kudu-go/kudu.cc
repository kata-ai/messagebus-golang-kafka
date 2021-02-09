#include <kudu/client/client.h>
#include <kudu/client/write_op.h>
#include <kudu/util/status.h>

#include "kudu.h"

#include <vector>
#include <string>
#include <iostream>

using namespace kudu::client;
using namespace kudu;
using std::cout;
using std::vector;
using std::string;


namespace {
template<class C1, class C2>
vector<C1> FromCVector(const C2* c_elems, int len) {
  vector<C1> ret;
  ret.reserve(len);
  for (int i = 0; i < len; i++) {
    ret.emplace_back(c_elems[i]);
  }
  return ret;
}
}

extern "C" {
  struct CKuduStatus {
    Status status;
  };

  struct CKuduClient {
    sp::shared_ptr<KuduClient> impl;
  };

  struct CKuduTableCreator {
    KuduTableCreator* impl;
  };

  struct CKuduSchema {
    KuduSchema impl;
  };

  struct CKuduSchemaBuilder {
    KuduSchemaBuilder* impl;
  };

  struct CKuduTable {
    sp::shared_ptr<KuduTable> impl;
  };

  struct CKuduSession {
    sp::shared_ptr<KuduSession> impl;
  };

  struct CKuduScanner {
    KuduScanner* impl;
  };

  struct CKuduScanBatch {
    KuduScanBatch impl;
    KuduScanBatch::RowPtr cur_row;
    int cur_idx = -1;
  };
  
  struct CKuduColumnSpec;
  struct CKuduWriteOperation;
  
  // Get the message from the status. Must be freed with free().
  char* KuduStatus_Message(CKuduStatus* self) {
    string msg = self->status.ToString();
    return strdup(msg.c_str());
  }
  
  void KuduStatus_Free(CKuduStatus* self) {
    delete self;
  }

  static CKuduStatus* MakeStatus(Status s) {
    if (s.ok()) return nullptr;
    return new CKuduStatus { std::move(s) };
  }
  

  //////////////////

  struct CKuduClientBuilder {
    KuduClientBuilder* impl;
  };

  CKuduClientBuilder* KuduClientBuilder_Create() {
    auto ret = new CKuduClientBuilder();
    ret->impl = new KuduClientBuilder();
    return ret;
  }


  void KuduClientBuilder_Free(CKuduClientBuilder* self) {
    delete self->impl;
    delete self;
  }

  void KuduClientBuilder_add_master_server_addr(
      CKuduClientBuilder* self, const char* addr) {
    self->impl->add_master_server_addr(string(addr));
  }

  CKuduStatus* KuduClientBuilder_Build(
      CKuduClientBuilder* self, CKuduClient** client) {

    sp::shared_ptr<KuduClient> k_client;
    Status s = self->impl->Build(&k_client);
    if (!s.ok()) return MakeStatus(std::move(s));
    *client = new CKuduClient { std::move(k_client) };
    return nullptr;
  }

  /////////////////////////////////

  CKuduTableCreator* KuduClient_NewTableCreator(CKuduClient* self) {
    return new CKuduTableCreator { self->impl->NewTableCreator() };
  }

  void KuduClient_Free(CKuduClient* client) {
    delete client;
  }

  CKuduStatus* KuduClient_TableExists(CKuduClient* self, const char* table_name, int* exists) {
    bool exists_b;
    Status s = self->impl->TableExists(string(table_name), &exists_b);
    if (!s.ok()) return MakeStatus(std::move(s));
    *exists = exists_b;
    return nullptr;
  }

  CKuduStatus* KuduClient_OpenTable(CKuduClient* self, const char* table_name, CKuduTable** table) {
    sp::shared_ptr<KuduTable> c_table;
    Status s = self->impl->OpenTable(string(table_name), &c_table);
    if (!s.ok()) return MakeStatus(std::move(s));
    *table = new CKuduTable { std::move(c_table) };
    return nullptr;
  }

  CKuduSession* KuduClient_NewSession(CKuduClient* self) {
    return new CKuduSession { self->impl->NewSession() };
  }

  //============================================================================================
  void KuduSession_Close(CKuduSession* self) {
    delete self;
  }

  CKuduStatus* KuduSession_SetFlushMode(CKuduSession* self, CFlushMode mode) {
    return MakeStatus(self->impl->SetFlushMode(static_cast<KuduSession::FlushMode>(mode)));
  }

  CKuduStatus* KuduSession_Apply(CKuduSession* self, CKuduWriteOperation* op) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(self->impl->Apply(cpp_op));
  }

  CKuduStatus* KuduSession_Flush(CKuduSession* self) {
    return MakeStatus(self->impl->Flush());
  }

  //============================================================================================
  CKuduStatus* KuduWriteOperation_SetBool(CKuduWriteOperation* op, const char* col_name, bool val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetBool(col_name, val));
  }
  CKuduStatus* KuduWriteOperation_SetInt8(CKuduWriteOperation* op, const char* col_name, int8_t val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetInt8(col_name, val));
  }
  CKuduStatus* KuduWriteOperation_SetInt16(CKuduWriteOperation* op, const char* col_name, int16_t val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetInt16(col_name, val));
  }
  CKuduStatus* KuduWriteOperation_SetInt32(CKuduWriteOperation* op, const char* col_name, int32_t val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetInt32(col_name, val));
  }
  CKuduStatus* KuduWriteOperation_SetInt64(CKuduWriteOperation* op, const char* col_name, int64_t val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetInt64(col_name, val));
  }
  CKuduStatus* KuduWriteOperation_SetUnixTimeMicros(CKuduWriteOperation* op, const char* col_name,  int64_t micros_since_utc_epoch) {
     auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetUnixTimeMicros(col_name, micros_since_utc_epoch));
  }
  CKuduStatus* KuduWriteOperation_SetFloat(CKuduWriteOperation* op, const char* col_name, float val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetFloat(col_name, val));
  }
  CKuduStatus* KuduWriteOperation_SetDouble(CKuduWriteOperation* op, const char* col_name, double val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetDouble(col_name, val));
  }
  CKuduStatus* KuduWriteOperation_SetUnscaledDecimal(CKuduWriteOperation* op, const char* col_name, int128_t val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetUnscaledDecimal(col_name, val));
  }
  CKuduStatus* KuduWriteOperation_SetString(CKuduWriteOperation* op, const char* col_name, char* val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetString(col_name, val));
  }
  CKuduStatus* KuduWriteOperation_SetBinary(CKuduWriteOperation* op, const char* col_name, Slice val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetBinary(col_name, val));
  }
  // with index
  CKuduStatus* KuduWriteOperation_SetBool_WithIndex(CKuduWriteOperation* op, const int col_idx, bool val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetBool(col_idx, val));
  }
  CKuduStatus* KuduWriteOperation_SetInt8_WithIndex(CKuduWriteOperation* op, const int col_idx, int8_t val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetInt8(col_idx, val));
  }
  CKuduStatus* KuduWriteOperation_SetInt16_WithIndex(CKuduWriteOperation* op, const int col_idx, int16_t val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetInt16(col_idx, val));
  }
  CKuduStatus* KuduWriteOperation_SetInt32_WithIndex(CKuduWriteOperation* op, const int col_idx, int32_t val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetInt32(col_idx, val));
  }
  CKuduStatus* KuduWriteOperation_SetInt64_WithIndex(CKuduWriteOperation* op, const int col_idx, int64_t val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetInt64(col_idx, val));
  }
  CKuduStatus* KuduWriteOperation_SetUnixTimeMicros_WithIndex(CKuduWriteOperation* op, const int col_idx, int64_t micros_since_utc_epoch) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetUnixTimeMicros(col_idx, micros_since_utc_epoch));
  } 
  CKuduStatus* KuduWriteOperation_SetFloat_WithIndex(CKuduWriteOperation* op, const int col_idx, float val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetFloat(col_idx, val));
  }
  CKuduStatus* KuduWriteOperation_SetDouble_WithIndex(CKuduWriteOperation* op, const int col_idx, double val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetDouble(col_idx, val));
  }
  CKuduStatus* KuduWriteOperation_SetUnscaledDecimal_WithIndex(CKuduWriteOperation* op, const int col_idx, int128_t val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetUnscaledDecimal(col_idx, val));
  }
  CKuduStatus* KuduWriteOperation_SetString_WithIndex(CKuduWriteOperation* op, const int col_idx, char* val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetString(col_idx, val));
  }
  CKuduStatus* KuduWriteOperation_SetBinary_WithIndex(CKuduWriteOperation* op, const int col_idx, Slice val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetBinary(col_idx, val));
  }
  CKuduStatus* KuduWriteOperation_SetNull(CKuduWriteOperation* op, const char* col_name) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetNull(col_name));
  }
  CKuduStatus* KuduWriteOperation_SetNull_WithIndex(CKuduWriteOperation* op, const int col_idx) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetNull(col_idx));
  }
  CKuduStatus* KuduWriteOperation_SetStringCopy_WithIndex(CKuduWriteOperation* op, const int col_idx, char* val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetStringCopy(col_idx, val));
  }
  CKuduStatus* KuduWriteOperation_SetStringCopy(CKuduWriteOperation* op, const char* col_name, char* val) {
    auto cpp_op = reinterpret_cast<KuduWriteOperation*>(op);
    return MakeStatus(cpp_op->mutable_row()->SetStringCopy(col_name, val));
  }

  //============================================================================================
  void KuduTable_Close(CKuduTable* self) {
    delete self;
  }
  CKuduWriteOperation* KuduTable_NewInsert(CKuduTable* self) {
    return reinterpret_cast<CKuduWriteOperation*>(self->impl->NewInsert());
  }

  CKuduScanner* KuduTable_NewScanner(CKuduTable* self) {
    return new CKuduScanner { new KuduScanner(self->impl.get()) };
  }

  //============================================================================================

  void KuduScanner_Free(CKuduScanner* self) {
    delete self->impl;
    delete self;
  }

  CKuduStatus* KuduScanner_SetProjectedColumns(CKuduScanner* self, const char** col_names, int n_cols) {
    return MakeStatus(self->impl->SetProjectedColumnNames(FromCVector<string>(col_names, n_cols)));
  }

  CKuduStatus* KuduScanner_Open(CKuduScanner* self) {
    return MakeStatus(self->impl->Open());
  }

  int KuduScanner_HasMoreRows(CKuduScanner* self) {
    return self->impl->HasMoreRows();
  }

  CKuduStatus* KuduScanner_NextBatch(CKuduScanner* self, CKuduScanBatch** batch) {
    if (*batch == nullptr) {
      *batch = new CKuduScanBatch();
    }
    (*batch)->cur_idx = -1;
    return MakeStatus(self->impl->NextBatch(&(*batch)->impl));
  }

  //============================================================================================

  void KuduScanBatch_Free(CKuduScanBatch* self) {
    delete self;
  }

  int KuduScanBatch_HasNext(CKuduScanBatch* self) {
    return self->cur_idx + 1 < self->impl.NumRows();
  }

  void KuduScanBatch_SeekNext(CKuduScanBatch* self) {
    self->cur_idx++;
    self->cur_row = self->impl.Row(self->cur_idx);
  }

  const char* KuduScanBatch_Row_ToString(CKuduScanBatch* self) {
    return strdup(self->cur_row.ToString().c_str());
  }
  
  //============================================================================================
  void KuduSchema_Free(CKuduSchema* self) {
    delete self;
  }

  //============================================================================================
  CKuduColumnSpec* KuduColumnSpec_SetType(CKuduColumnSpec* self, int type) {
    reinterpret_cast<KuduColumnSpec*>(self)->Type(
        KuduColumnSchema::DataType(type));
    // TODO(todd) expose enum for datatype
    return self;
  }

  CKuduColumnSpec* KuduColumnSpec_SetNotNull(CKuduColumnSpec* self) {
    reinterpret_cast<KuduColumnSpec*>(self)->NotNull();
    return self;
  }

  
  //============================================================================================
  CKuduSchemaBuilder* KuduSchemaBuilder_Create() {
    return new CKuduSchemaBuilder { new KuduSchemaBuilder() };
  }
  void KuduSchemaBuilder_Free(CKuduSchemaBuilder* self) {
    delete self->impl;
    delete self;
  }

  CKuduColumnSpec* KuduSchemaBuilder_AddColumn(CKuduSchemaBuilder* self, const char* name) {
    return reinterpret_cast<CKuduColumnSpec*>(self->impl->AddColumn(string(name)));
  }
  void KuduSchemaBuilder_SetPrimaryKey(CKuduSchemaBuilder* self, const char** col_names, int n_cols) {
    auto names = FromCVector<string>(col_names, n_cols);
    self->impl->SetPrimaryKey(names);
  }
  CKuduStatus* KuduSchemaBuilder_Build(CKuduSchemaBuilder* self, CKuduSchema** out_schema) {
    KuduSchema schema;
    Status s = self->impl->Build(&schema);
    if (!s.ok()) return MakeStatus(std::move(s));
    *out_schema = new CKuduSchema { std::move(schema) };
    return nullptr;
  }

  //============================================================================================

  void KuduTableCreator_Free(CKuduTableCreator* self) {
    delete self->impl;
    delete self;
  }

  void KuduTableCreator_SetTableName(CKuduTableCreator* self, const char* name) {
    self->impl->table_name(string(name));
  }

  void KuduTableCreator_SetSchema(CKuduTableCreator* self, CKuduSchema* schema) {
    self->impl->schema(&schema->impl);
  }
  
  void KuduTableCreator_AddHashPartitions(CKuduTableCreator* self, const char** col_names, int n_cols, int num_buckets) {
    auto names = FromCVector<string>(col_names, n_cols);
    self->impl->add_hash_partitions(names, num_buckets);
    
  }

  
  CKuduStatus* KuduTableCreator_Create(CKuduTableCreator* self) {
    Status s = self->impl->Create();
    if (!s.ok()) return MakeStatus(std::move(s));
    return nullptr;
  }

}