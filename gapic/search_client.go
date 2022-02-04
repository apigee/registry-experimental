// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go_gapic. DO NOT EDIT.

package gapic

import (
	"context"
	"math"
	"time"

	"cloud.google.com/go/longrunning"
	lroauto "cloud.google.com/go/longrunning/autogen"
	rpcpb "github.com/apigee/registry-experimental/rpc"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

var newSearchClientHook clientHook

// SearchCallOptions contains the retry settings for each method of SearchClient.
type SearchCallOptions struct {
	Index []gax.CallOption
	Query []gax.CallOption
}

func defaultSearchGRPCClientOptions() []option.ClientOption {
	return []option.ClientOption{
		internaloption.WithDefaultEndpoint("apigeeregistry.googleapis.com:443"),
		internaloption.WithDefaultMTLSEndpoint("apigeeregistry.mtls.googleapis.com:443"),
		internaloption.WithDefaultAudience("https://apigeeregistry.googleapis.com/"),
		internaloption.WithDefaultScopes(DefaultAuthScopes()...),
		internaloption.EnableJwtWithScope(),
		option.WithGRPCDialOption(grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}

func defaultSearchCallOptions() *SearchCallOptions {
	return &SearchCallOptions{
		Index: []gax.CallOption{
		},
		Query: []gax.CallOption{
		},
	}
}

// internalSearchClient is an interface that defines the methods availaible from .
type internalSearchClient interface {
	Close() error
	setGoogleClientInfo(...string)
	Connection() *grpc.ClientConn
	Index(context.Context, *rpcpb.IndexRequest, ...gax.CallOption) (*IndexOperation, error)
	IndexOperation(name string) *IndexOperation
	Query(context.Context, *rpcpb.QueryRequest, ...gax.CallOption) *QueryResponse_ResultIterator
}

// SearchClient is a client for interacting with .
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
//
// Build and search an index of APIs.
type SearchClient struct {
	// The internal transport-dependent client.
	internalClient internalSearchClient

	// The call options for this service.
	CallOptions *SearchCallOptions

	// LROClient is used internally to handle long-running operations.
	// It is exposed so that its CallOptions can be modified if required.
	// Users should not Close this client.
	LROClient *lroauto.OperationsClient

}

// Wrapper methods routed to the internal client.

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *SearchClient) Close() error {
	return c.internalClient.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *SearchClient) setGoogleClientInfo(keyval ...string) {
	c.internalClient.setGoogleClientInfo(keyval...)
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *SearchClient) Connection() *grpc.ClientConn {
	return c.internalClient.Connection()
}

// Index add a resource to the search index.
func (c *SearchClient) Index(ctx context.Context, req *rpcpb.IndexRequest, opts ...gax.CallOption) (*IndexOperation, error) {
	return c.internalClient.Index(ctx, req, opts...)
}

// IndexOperation returns a new IndexOperation from a given name.
// The name must be that of a previously created IndexOperation, possibly from a different process.
func (c *SearchClient) IndexOperation(name string) *IndexOperation {
	return c.internalClient.IndexOperation(name)
}

// Query query the index.
func (c *SearchClient) Query(ctx context.Context, req *rpcpb.QueryRequest, opts ...gax.CallOption) *QueryResponse_ResultIterator {
	return c.internalClient.Query(ctx, req, opts...)
}

// searchGRPCClient is a client for interacting with  over gRPC transport.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type searchGRPCClient struct {
	// Connection pool of gRPC connections to the service.
	connPool gtransport.ConnPool

	// flag to opt out of default deadlines via GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE
	disableDeadlines bool

	// Points back to the CallOptions field of the containing SearchClient
	CallOptions **SearchCallOptions

	// The gRPC API client.
	searchClient rpcpb.SearchClient

	// LROClient is used internally to handle long-running operations.
	// It is exposed so that its CallOptions can be modified if required.
	// Users should not Close this client.
	LROClient **lroauto.OperationsClient

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewSearchClient creates a new search client based on gRPC.
// The returned client must be Closed when it is done being used to clean up its underlying connections.
//
// Build and search an index of APIs.
func NewSearchClient(ctx context.Context, opts ...option.ClientOption) (*SearchClient, error) {
	clientOpts := defaultSearchGRPCClientOptions()
	if newSearchClientHook != nil {
		hookOpts, err := newSearchClientHook(ctx, clientHookParams{})
		if err != nil {
			return nil, err
		}
		clientOpts = append(clientOpts, hookOpts...)
	}

	disableDeadlines, err := checkDisableDeadlines()
	if err != nil {
		return nil, err
	}

	connPool, err := gtransport.DialPool(ctx, append(clientOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	client := SearchClient{CallOptions: defaultSearchCallOptions()}

	c := &searchGRPCClient{
		connPool:    connPool,
		disableDeadlines: disableDeadlines,
		searchClient: rpcpb.NewSearchClient(connPool),
		CallOptions: &client.CallOptions,

	}
	c.setGoogleClientInfo()

	client.internalClient = c

	client.LROClient, err = lroauto.NewOperationsClient(ctx, gtransport.WithConnPool(connPool))
	if err != nil {
		// This error "should not happen", since we are just reusing old connection pool
		// and never actually need to dial.
		// If this does happen, we could leak connp. However, we cannot close conn:
		// If the user invoked the constructor with option.WithGRPCConn,
		// we would close a connection that's still in use.
		// TODO: investigate error conditions.
		return nil, err
	}
	c.LROClient = &client.LROClient
	return &client, nil
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *searchGRPCClient) Connection() *grpc.ClientConn {
	return c.connPool.Conn()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *searchGRPCClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *searchGRPCClient) Close() error {
	return c.connPool.Close()
}

func (c *searchGRPCClient) Index(ctx context.Context, req *rpcpb.IndexRequest, opts ...gax.CallOption) (*IndexOperation, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).Index[0:len((*c.CallOptions).Index):len((*c.CallOptions).Index)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.searchClient.Index(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &IndexOperation{
		lro: longrunning.InternalNewOperation(*c.LROClient, resp),
	}, nil
}

func (c *searchGRPCClient) Query(ctx context.Context, req *rpcpb.QueryRequest, opts ...gax.CallOption) *QueryResponse_ResultIterator {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).Query[0:len((*c.CallOptions).Query):len((*c.CallOptions).Query)], opts...)
	it := &QueryResponse_ResultIterator{}
	req = proto.Clone(req).(*rpcpb.QueryRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*rpcpb.QueryResponse_Result, string, error) {
		resp := &rpcpb.QueryResponse{}
		if pageToken != "" {
			req.PageToken = pageToken
		}
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else if pageSize != 0 {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.searchClient.Query(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}

		it.Response = resp
		return resp.GetResults(), resp.GetNextPageToken(), nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}

	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	it.pageInfo.MaxSize = int(req.GetPageSize())
	it.pageInfo.Token = req.GetPageToken()

	return it
}

// IndexOperation manages a long-running operation from Index.
type IndexOperation struct {
	lro *longrunning.Operation
}

// IndexOperation returns a new IndexOperation from a given name.
// The name must be that of a previously created IndexOperation, possibly from a different process.
func (c *searchGRPCClient) IndexOperation(name string) *IndexOperation {
	return &IndexOperation{
		lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.
//
// See documentation of Poll for error-handling information.
func (op *IndexOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*rpcpb.IndexResponse, error) {
	var resp rpcpb.IndexResponse
	if err := op.lro.WaitWithInterval(ctx, &resp, time.Minute, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully,
// op.Done will return true, and the response of the operation is returned.
// If Poll succeeds and the operation has not completed, the returned response and error are both nil.
func (op *IndexOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*rpcpb.IndexResponse, error) {
	var resp rpcpb.IndexResponse
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *IndexOperation) Metadata() (*rpcpb.IndexMetadata, error) {
	var meta rpcpb.IndexMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *IndexOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *IndexOperation) Name() string {
	return op.lro.Name()
}

// QueryResponse_ResultIterator manages a stream of *rpcpb.QueryResponse_Result.
type QueryResponse_ResultIterator struct {
	items    []*rpcpb.QueryResponse_Result
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// Response is the raw response for the current page.
	// It must be cast to the RPC response type.
	// Calling Next() or InternalFetch() updates this value.
	Response interface{}

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*rpcpb.QueryResponse_Result, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *QueryResponse_ResultIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *QueryResponse_ResultIterator) Next() (*rpcpb.QueryResponse_Result, error) {
	var item *rpcpb.QueryResponse_Result
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *QueryResponse_ResultIterator) bufLen() int {
	return len(it.items)
}

func (it *QueryResponse_ResultIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}

func (c *SearchClient) GrpcClient() rpcpb.SearchClient {
	return c.internalClient.(*searchGRPCClient).searchClient
}
