package messagebus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/linkedin/goavro/v2"
)

// schemaRegistryClient allows interactions with
// Schema Registry over HTTP. Applications using
// this client can retrieve data about schemas,
// which in turn can be used to Serialize and
// Deserialize data.
type schemaRegistryClient struct {
	schemaRegistryURL      string
	credentials            *credentials
	httpClient             *http.Client
	cachingEnabled         bool
	codecCreationEnabled   bool
	idSchemaCache          map[int]*Schema
	idSchemaCacheLock      sync.RWMutex
	subjectSchemaCache     map[string]*Schema
	subjectSchemaCacheLock sync.RWMutex
}

type SchemaType string

const (
	Protobuf SchemaType = "PROTOBUF"
	Avro     SchemaType = "AVRO"
	Json     SchemaType = "JSON"
)

func (s SchemaType) String() string {
	return string(s)
}

// Schema is a data structure that holds all
// the relevant information about schemas.
type Schema struct {
	id      int
	schema  string
	version int
	codec   *goavro.Codec
}

type credentials struct {
	username string
	password string
}

type schemaRequest struct {
	Schema string `json:"schema"`
}

type schemaResponse struct {
	Subject string `json:"subject"`
	Version int    `json:"version"`
	Schema  string `json:"schema"`
	ID      int    `json:"id"`
}

const (
	schemaByID       = "/schemas/ids/%d"
	subjectVersions  = "/subjects/%s/versions"
	subjectByVersion = "/subjects/%s/versions/%s"
	contentType      = "application/vnd.schemaregistry.v1+json"
)

// createSchemaRegistryClient creates a client that allows
// interactions with Schema Registry over HTTP. Applications
// using this client can retrieve data about schemas, which
// in turn can be used to Serialize and Deserialize records.
func createSchemaRegistryClient(schemaRegistryURL string) *schemaRegistryClient {
	return &schemaRegistryClient{schemaRegistryURL: schemaRegistryURL,
		httpClient:     &http.Client{Timeout: 5 * time.Second},
		cachingEnabled: true, codecCreationEnabled: true,
		idSchemaCache:      make(map[int]*Schema),
		subjectSchemaCache: make(map[string]*Schema)}
}

// getSchema gets the schema associated with the given id.
func (client *schemaRegistryClient) getSchema(schemaID int) (*Schema, error) {

	if client.cachingEnabled {
		client.idSchemaCacheLock.RLock()
		cachedSchema := client.idSchemaCache[schemaID]
		client.idSchemaCacheLock.RUnlock()
		if cachedSchema != nil {
			return cachedSchema, nil
		}
	}

	resp, err := client.httpRequest("GET", fmt.Sprintf(schemaByID, schemaID), nil)
	if err != nil {
		return nil, err
	}

	var schemaResp = new(schemaResponse)
	err = json.Unmarshal(resp, &schemaResp)
	if err != nil {
		return nil, err
	}
	var codec *goavro.Codec
	if client.codecCreationEnabled {
		codec, err = goavro.NewCodec(schemaResp.Schema)
		if err != nil {
			return nil, err
		}
	}
	var schema = &Schema{
		id:     schemaID,
		schema: schemaResp.Schema,
		codec:  codec,
	}

	if client.cachingEnabled {
		client.idSchemaCacheLock.Lock()
		client.idSchemaCache[schemaID] = schema
		client.idSchemaCacheLock.Unlock()
	}

	return schema, nil
}

// getLatestSchema gets the schema associated with the given subject.
// The schema returned contains the last version for that subject.
func (client *schemaRegistryClient) getLatestSchema(subject string, isKey bool) (*Schema, error) {

	// In order to ensure consistency, we need
	// to temporarily disable caching to force
	// the retrieval of the latest release from
	// Schema Registry.
	cachingEnabled := client.cachingEnabled
	client.isCachingEnabled(false)
	schema, err := client.getVersion(subject, "latest", isKey)
	client.isCachingEnabled(cachingEnabled)

	return schema, err
}

// getSchemaVersions returns a list of versions from a given subject.
func (client *schemaRegistryClient) getSchemaVersions(subject string, isKey bool) ([]int, error) {

	concreteSubject := subject
	resp, err := client.httpRequest("GET", fmt.Sprintf(subjectVersions, concreteSubject), nil)
	if err != nil {
		return nil, err
	}

	var versions = []int{}
	err = json.Unmarshal(resp, &versions)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

// getSchemaByVersion gets the schema associated with the given subject.
// The schema returned contains the version specified as a parameter.
func (client *schemaRegistryClient) getSchemaByVersion(subject string, version int, isKey bool) (*Schema, error) {
	return client.getVersion(subject, strconv.Itoa(version), isKey)
}

// createSchema creates a new schema in Schema Registry and associates
// with the subject provided. It returns the newly created schema with
// all its associated information.
func (client *schemaRegistryClient) createSchema(subject string, schema string,
	schemaType SchemaType, isKey bool) (*Schema, error) {

	concreteSubject := subject

	switch schemaType {
	case Avro, Json:
		compiledRegex := regexp.MustCompile(`\r?\n`)
		schema = compiledRegex.ReplaceAllString(schema, " ")
	case Protobuf:
		break
	default:
		return nil, fmt.Errorf("invalid schema type. valid values are Avro, Json, or Protobuf")
	}
	schemaReq := schemaRequest{Schema: schema}
	schemaBytes, err := json.Marshal(schemaReq)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(schemaBytes)
	resp, err := client.httpRequest("POST", fmt.Sprintf(subjectVersions, concreteSubject), payload)
	if err != nil {
		return nil, err
	}

	schemaResp := new(schemaResponse)
	err = json.Unmarshal(resp, &schemaResp)
	if err != nil {
		return nil, err
	}
	// Conceptually, the schema returned below will be the
	// exactly same one created above. However, since Schema
	// Registry can have multiple concurrent clients writing
	// schemas, this may produce an incorrect result. Thus,
	// this logic strongly relies on the idempotent guarantees
	// from Schema Registry, as well as in the best practice
	// that schemas don't change very often.
	newSchema, err := client.getSchema(schemaResp.ID)
	if err != nil {
		return nil, err
	}

	if client.cachingEnabled {

		// Update the subject-2-schema cache
		cacheKey := cacheKey(concreteSubject,
			strconv.Itoa(newSchema.version))
		client.subjectSchemaCacheLock.Lock()
		client.subjectSchemaCache[cacheKey] = newSchema
		client.subjectSchemaCacheLock.Unlock()

		// Update the id-2-schema cache
		client.idSchemaCacheLock.Lock()
		client.idSchemaCache[newSchema.id] = newSchema
		client.idSchemaCacheLock.Unlock()

	}

	return newSchema, nil
}

// setCredentials allows users to set credentials to be
// used with Schema Registry, for scenarios when Schema
// Registry has authentication enabled.
func (client *schemaRegistryClient) setCredentials(username string, password string) {
	if len(username) > 0 && len(password) > 0 {
		credentials := credentials{username, password}
		client.credentials = &credentials
	}
}

// setTimeout allows the client to be reconfigured about
// how much time internal HTTP requests will take until
// they timeout. FYI, It defaults to five seconds.
func (client *schemaRegistryClient) setTimeout(timeout time.Duration) {
	client.httpClient.Timeout = timeout
}

// isCachingEnabled allows the client to cache any values
// that have been returned, which may speed up performance
// if these values rarely changes.
func (client *schemaRegistryClient) isCachingEnabled(value bool) {
	client.cachingEnabled = value
}

// isCodecCreationEnabled allows the application to enable/disable
// the automatic creation of codec's when schemas are returned.
func (client *schemaRegistryClient) isCodecCreationEnabled(value bool) {
	client.codecCreationEnabled = value
}

func (client *schemaRegistryClient) getVersion(subject string,
	version string, isKey bool) (*Schema, error) {

	concreteSubject := subject

	if client.cachingEnabled {
		cacheKey := cacheKey(concreteSubject, version)
		client.subjectSchemaCacheLock.RLock()
		cachedResult := client.subjectSchemaCache[cacheKey]
		client.subjectSchemaCacheLock.RUnlock()
		if cachedResult != nil {
			return cachedResult, nil
		}
	}

	resp, err := client.httpRequest("GET", fmt.Sprintf(subjectByVersion, concreteSubject, version), nil)
	if err != nil {
		return nil, err
	}

	schemaResp := new(schemaResponse)
	err = json.Unmarshal(resp, &schemaResp)
	if err != nil {
		return nil, err
	}
	var codec *goavro.Codec
	if client.codecCreationEnabled {
		codec, err = goavro.NewCodec(schemaResp.Schema)
		if err != nil {
			return nil, err
		}
	}
	var schema = &Schema{
		id:      schemaResp.ID,
		schema:  schemaResp.Schema,
		version: schemaResp.Version,
		codec:   codec,
	}

	if client.cachingEnabled {

		// Update the subject-2-schema cache
		cacheKey := cacheKey(concreteSubject, version)
		client.subjectSchemaCacheLock.Lock()
		client.subjectSchemaCache[cacheKey] = schema
		client.subjectSchemaCacheLock.Unlock()

		// Update the id-2-schema cache
		client.idSchemaCacheLock.Lock()
		client.idSchemaCache[schema.id] = schema
		client.idSchemaCacheLock.Unlock()

	}

	return schema, nil
}

func (client *schemaRegistryClient) httpRequest(method, uri string, payload io.Reader) ([]byte, error) {

	url := fmt.Sprintf("%s%s", client.schemaRegistryURL, uri)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	if client.credentials != nil {
		req.SetBasicAuth(client.credentials.username, client.credentials.password)
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, createError(resp)
	}

	return ioutil.ReadAll(resp.Body)
}

// ID ensures access to ID
func (schema *Schema) ID() int {
	return schema.id
}

// Schema ensures access to Schema
func (schema *Schema) Schema() string {
	return schema.schema
}

// Version ensures access to Version
func (schema *Schema) Version() int {
	return schema.version
}

// Codec ensures access to Codec
func (schema *Schema) Codec() *goavro.Codec {
	return schema.codec
}

func (schema Schema) FullName() string {
	var dat map[string]string
	_ = json.Unmarshal([]byte(schema.schema), &dat)
	switch dat["namespace"] {
	case "":
		return dat["name"]
	default:
		return fmt.Sprintf("%s.%s", dat["namespace"], dat["name"])
	}
}

func cacheKey(subject string, version string) string {
	return fmt.Sprintf("%s-%s", subject, version)
}

func createError(resp *http.Response) error {
	decoder := json.NewDecoder(resp.Body)
	var errorResp struct {
		ErrorCode int    `json:"error_code"`
		Message   string `json:"message"`
	}
	err := decoder.Decode(&errorResp)
	if err == nil {
		return fmt.Errorf("%s: %s", resp.Status, errorResp.Message)
	}
	return fmt.Errorf("%s", resp.Status)
}
