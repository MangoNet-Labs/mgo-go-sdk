package httpconn

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"net/http"
	"time"

	"github.com/mangonet-labs/mgo-go-sdk/model/request"

	"golang.org/x/time/rate"
)

const defaultTimeout = time.Second * 5

type HttpConn struct {
	c       *http.Client
	rl      *rate.Limiter
	rpcUrl  string
	timeout time.Duration
}

// newDefaultRateLimiter returns a new rate.Limiter with a rate of 10000 requests
// every 1 second.
func newDefaultRateLimiter() *rate.Limiter {
	rateLimiter := rate.NewLimiter(rate.Every(1*time.Second), 10000) // 10000 request every 1 seconds
	return rateLimiter
}

// NewHttpConn creates a new HttpConn with the specified RPC URL.
// It initializes the HTTP client and sets the default timeout for requests.
func NewHttpConn(rpcUrl string) *HttpConn {
	return &HttpConn{
		c:       &http.Client{},
		rpcUrl:  rpcUrl,
		timeout: defaultTimeout,
	}
}

// NewCustomHttpConn creates a new HttpConn with the specified RPC URL and a custom
// http.Client. This is useful if you want to set custom timeouts, transport, or other
// options for the HTTP client.
func NewCustomHttpConn(rpcUrl string, cli *http.Client) *HttpConn {
	return &HttpConn{
		c:       cli,
		rpcUrl:  rpcUrl,
		timeout: defaultTimeout,
	}
}

// Request sends a JSON-RPC request to the server using the specified operation.
// It constructs a JSON-RPC request payload with the given method and parameters,
// sends it via HTTP POST to the configured RPC URL, and returns the response body as bytes.
//
// Parameters:
// - ctx: The context to control cancellation and timeout.
// - op: The Operation containing method and parameters for the RPC call.
//
// Returns:
// - []byte: The response body from the server.
// - error: An error if the request fails or the response cannot be read.
func (h *HttpConn) Request(ctx context.Context, op Operation) ([]byte, error) {
	jsonRPCReq := request.JsonRPCRequest{
		JsonRPC: "2.0",
		ID:      time.Now().UnixMilli(),
		Method:  op.Method,
		Params:  op.Params,
	}
	reqBytes, err := json.Marshal(jsonRPCReq)
	if err != nil {
		return []byte{}, err
	}

	request, err := http.NewRequest("POST", h.rpcUrl, bytes.NewBuffer(reqBytes))
	if err != nil {
		return []byte{}, err
	}
	request = request.WithContext(ctx)
	request.Header.Add("Content-Type", "application/json")
	rsp, err := h.c.Do(request.WithContext(ctx))
	if err != nil {
		return []byte{}, err
	}
	defer rsp.Body.Close()

	bodyBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return []byte{}, err
	}
	return bodyBytes, nil
}
