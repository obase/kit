package kit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/obase/conf"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

const HTTP_CKEY = "kit.http"

const (
	ProxyBufferPool_None = "none" // 没有缓存池
	ProxyBufferPool_Sync = "sync" // 采用sync.Pool

	ProxyErrorHandler_None = "none" // 没有错误处理
	ProxyErrorHandler_Body = "body" // 将错误写到body

	REVERSE_SCHEME = "x-rscheme"
	REVERSE_HOST   = "x-rhost"
	REVERSE_PATH   = "x-rpath"
)
const HTTP_BLOCK_SIZE = 32 * 1024

type HttpConfig struct {
	// Timeout is the maximum amount of time a dial will wait for
	// a connect to complete. If Deadline is also set, it may fail
	// earlier.
	//
	// The default is no timeout.
	//
	// When using TCP and dialing a host name with multiple IP
	// addresses, the timeout may be divided between them.
	//
	// With or without a timeout, the operating system may impose
	// its own earlier timeout. For instance, TCP timeouts are
	// often around 3 minutes.
	ConnectTimeout time.Duration `json:"connectTimeout" yaml:"connectTimeout"`

	// KeepAlive specifies the keep-alive period for an active
	// network connection.
	// If zero, keep-alives are enabled if supported by the protocol
	// and operating system. Network protocols or operating systems
	// that do not support keep-alives ignore this field.
	// If negative, keep-alives are disabled.
	KeepAlive time.Duration `json:"keepAlive" yaml:"keepAlive"`

	// MaxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	MaxIdleConns int `json:"maxIdleConns"  yaml:"maxIdleConns"`

	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// DefaultMaxIdleConnsPerHost is used.
	MaxIdleConnsPerHost int `json:"maxIdleConnsPerHost" yaml:"maxIdleConnsPerHost"`

	// MaxConnsPerHost optionally limits the total number of
	// connections per host, including connections in the dialing,
	// active, and idle states. On limit violation, dials will block.
	//
	// Zero means no limit.
	//
	// For HTTP/2, this currently only controls the number of new
	// connections being created at a time, instead of the total
	// number. In practice, hosts using HTTP/2 only have about one
	// idle connection, though.
	MaxConnsPerHost int `json:"maxConnsPerHost" yaml:"maxConnsPerHost"`

	// IdleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself.
	// Zero means no limit.
	IdleConnTimeout time.Duration `json:"idleConnTimeout" yaml:"idleConnTimeout"`
	// TLSHandshakeTimeout specifies the maximum amount of time waiting to
	// wait for a TLS handshake. Zero means no timeout.
	TLSHandshakeTimeout time.Duration `json:"tlsHandshakeTimeout" bson:"tlsHandshakeTimeout"`

	// DisableCompression, if true, prevents the Transport from
	// requesting compression with an "Accept-Encoding: gzip"
	// request header when the Request contains no existing
	// Accept-Encoding value. If the Transport requests gzip on
	// its own and gets a gzipped response, it's transparently
	// decoded in the Response.Body. However, if the user
	// explicitly requested gzip it is not automatically
	// uncompressed.
	DisableCompression    bool `json:"disableCompression" yaml:"disableCompression"`
	DisableCompressionSet bool

	// ResponseHeaderTimeout, if non-zero, specifies the amount of
	// time to wait for a server's response headers after fully
	// writing the request (including its body, if any). This
	// time does not include the time to read the response body.
	ResponseHeaderTimeout time.Duration `json:"responseHeaderTimeout" yaml:"responseHeaderTimeout"`

	// ExpectContinueTimeout, if non-zero, specifies the amount of
	// time to wait for a server's first response headers after fully
	// writing the request headers if the request has an
	// "Expect: 100-continue" header. Zero means no timeout and
	// causes the body to be sent immediately, without
	// waiting for the server to approve.
	// This time does not include the time to send the request header.
	ExpectContinueTimeout time.Duration `json:"expectContinueTimeout" yaml:"expectContinueTimeout"`

	// MaxResponseHeaderBytes specifies a limit on how many
	// response bytes are allowed in the server's response
	// header.
	//
	// Zero means to use a default limit.
	MaxResponseHeaderBytes int64 `json:"maxResponseHeaderBytes" yaml:"maxResponseHeaderBytes"`

	// WriteBufferSize specifies the size of the write buffer used
	// when writing to the transport.
	// If zero, a default (currently 4KB) is used.
	WriteBufferSize int `json:"writeBufferSize" yaml:"writeBufferSize"`

	// ReadBufferSize specifies the size of the read buffer used
	// when reading from the transport.
	// If zero, a default (currently 4KB) is used.
	ReadBufferSize int `json:"readBufferSize" bson:"readBufferSize"`

	// ForceAttemptHTTP2 controls whether HTTP/2 is enabled when a non-zero
	// Dial, DialTLS, or DialContext func or TLSClientConfig is provided.
	// By default, use of any those fields conservatively disables HTTP/2.
	// To use a custom dialer or TLS config and still attempt HTTP/2
	// upgrades, set this to true.
	ForceAttemptHTTP2    bool `json:"forceAttemptHTTP2" bson:"forceAttemptHTTP2"`
	ForceAttemptHTTP2Set bool

	// Timeout specifies a time limit for requests made by this
	// Client. The timeout includes connection time, any
	// redirects, and reading the response body. The timer remains
	// running after Get, Head, Post, or Do return and will
	// interrupt reading of the Response.Body.
	//
	// A Timeout of zero means no timeout.
	//
	// The Client cancels requests to the underlying Transport
	// as if the Request's Context ended.
	//
	// For compatibility, the Client will also use the deprecated
	// CancelRequest method on Transport if found. New
	// RoundTripper implementations should use the Request's Context
	// for cancelation instead of implementing CancelRequest.
	RequestTimeout time.Duration `json:"requestTimeout" yaml:"requestTimeout"`

	// FlushInterval specifies the flush interval
	// to flush to the client while copying the
	// response body.
	// If zero, no periodic flushing is done.
	// A negative value means to flush immediately
	// after each write to the client.
	// The FlushInterval is ignored when ReverseProxy
	// recognizes a response as a streaming response;
	// for such responses, writes are flushed to the client
	// immediately.
	ProxyFlushInterval time.Duration `json:"proxyFlushInterval" yaml:"proxyFlushInterval"`

	// BufferPool optionally specifies a buffer pool to
	// get byte slices for use by io.CopyBuffer when copying HTTP response bodies.
	// Values: none, sync
	ProxyBufferPool string `json:"proxyBufferPool" yaml:"proxyBufferPool"`

	// ErrorHandler is an optional function that handles errors
	// reaching the backend or errors from ModifyResponse.
	// Values: none, body
	ProxyErrorHandler string `json:"proxyErrorHandler" yaml:"proxyErrorHandler"`
}

func init() {
	var c HttpConfig
	if cnf, ok := conf.Get(HTTP_CKEY); ok {
		if err := conf.Convert(cnf, &c); err == nil {
			_, c.DisableCompressionSet = conf.Elem(cnf, "disableCompression")
			_, c.ForceAttemptHTTP2Set = conf.Elem(cnf, "forceAttemptHTTP2")
		}
	}
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = 30 * time.Second
	}
	if c.KeepAlive == 0 {
		c.KeepAlive = 30 * time.Second
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 100
	}
	if c.TLSHandshakeTimeout == 0 {
		c.TLSHandshakeTimeout = 10 * time.Second
	}
	if c.ExpectContinueTimeout == 0 {
		c.ExpectContinueTimeout = 1 * time.Second
	}

	if c.ProxyBufferPool == "" {
		c.ProxyBufferPool = ProxyBufferPool_Sync
	}
	if c.ProxyErrorHandler == "" {
		c.ProxyErrorHandler = ProxyErrorHandler_Body
	}
	ProxyFlushInterval = c.ProxyFlushInterval
	HttpTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   c.ConnectTimeout,
			KeepAlive: c.KeepAlive,
		}).DialContext,
		ForceAttemptHTTP2:      IfBool(c.ForceAttemptHTTP2Set || c.ForceAttemptHTTP2, c.ForceAttemptHTTP2, true),
		MaxIdleConns:           c.MaxIdleConns,
		MaxIdleConnsPerHost:    c.MaxIdleConnsPerHost,
		MaxConnsPerHost:        c.MaxConnsPerHost,
		IdleConnTimeout:        c.IdleConnTimeout,
		TLSHandshakeTimeout:    c.TLSHandshakeTimeout,
		DisableCompression:     IfBool(c.DisableCompressionSet || c.DisableCompression, c.DisableCompression, false),
		ResponseHeaderTimeout:  c.ResponseHeaderTimeout,
		ExpectContinueTimeout:  c.ExpectContinueTimeout,
		MaxResponseHeaderBytes: c.MaxResponseHeaderBytes,
		WriteBufferSize:        c.WriteBufferSize,
		ReadBufferSize:         c.ReadBufferSize,
	}
	HttpClient = &http.Client{
		Transport: HttpTransport,
		Timeout:   c.RequestTimeout,
	}

	ReverseProxy = &httputil.ReverseProxy{
		Transport:     HttpTransport,
		FlushInterval: c.ProxyFlushInterval,
		Director: func(req *http.Request) {
			req.URL.Scheme = req.Header.Get(REVERSE_SCHEME)
			req.URL.Host = req.Header.Get(REVERSE_HOST)
			req.URL.Path = req.Header.Get(REVERSE_PATH)
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to  value
				req.Header.Set("User-Agent", "")
			}
		},
		BufferPool:   proxyBufferPool(c.ProxyBufferPool),
		ErrorHandler: proxyErrorHandler(c.ProxyErrorHandler),
	}
}

var (
	HttpTransport      *http.Transport
	HttpClient         *http.Client
	ReverseProxy       *httputil.ReverseProxy
	ProxyFlushInterval time.Duration
)

type HttpError string

func (h HttpError) Error() string {
	return string(h)
}

type httpBufferPool struct {
	*sync.Pool
}

func (s *httpBufferPool) Get() []byte {
	ret := s.Pool.Get().([]byte)
	if len(ret) == 0 {
		ret = ret[:cap(ret)] // httpBufferPool的bytes长度必须大于,否则会抛panic
	}
	return ret
}
func (s *httpBufferPool) Put(v []byte) {
	s.Pool.Put(v)
}

func proxyBufferPool(name string) httputil.BufferPool {
	switch name {
	case ProxyBufferPool_None:
		return nil
	case ProxyBufferPool_Sync:
		return &httpBufferPool{Pool: &blockBufferPool}
	}
	panic("invalid proxy buffer pool type: " + name)
}

func proxyErrorHandler(name string) func(w http.ResponseWriter, r *http.Request, err error) {
	switch name {
	case ProxyErrorHandler_None:
		return nil
	case ProxyErrorHandler_Body:
		return func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, " proxy error: %v", err)
		}
	}
	panic("invalid proxy error handler type: " + name)
}

func JoinQuery(rurl string, params map[string]string) string {
	if len(params) > 0 {
		buf := GetBytesBuffer()
		buf.WriteString(rurl)
		first := true
		for k, v := range params {
			if first {
				buf.WriteByte('?')
				first = false
			} else {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
		rurl = buf.String()
		PutBytesBuffer(buf)
	}
	return rurl
}

func HttpRawRequest(method string, url string, header map[string]string, body io.Reader) (state int, content string, err error) {

	// 创建请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}
	rsp, err := HttpClient.Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	state = rsp.StatusCode
	buf := GetBytesBufferN(HTTP_BLOCK_SIZE)
	bss := GetBlockBufferN(HTTP_BLOCK_SIZE)
	if _, err = io.CopyBuffer(buf, rsp.Body, bss); err == nil {
		content = buf.String()
	}
	PutBlockBuffer(bss)
	PutBytesBuffer(buf)
	return
}

// 适用于大多数情况下的ContentType都是application/json,如果不需要请用HttpRawRequest
func HttpRequest(method string, url string, header map[string]string, body io.Reader) (state int, content string, err error) {
	// 创建请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	rsp, err := HttpClient.Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	state = rsp.StatusCode
	buf := GetBytesBufferN(HTTP_BLOCK_SIZE)
	bss := GetBlockBufferN(HTTP_BLOCK_SIZE)
	if _, err = io.CopyBuffer(buf, rsp.Body, bss); err == nil {
		content = buf.String()
	}
	PutBlockBuffer(bss)
	PutBytesBuffer(buf)
	return
}

func HttpJson(method string, url string, header map[string]string, reqobj interface{}, rspobj interface{}) (status int, err error) {
	var body io.Reader
	if reqobj != nil {
		var data []byte
		data, err = json.Marshal(reqobj)
		if err != nil {
			return
		}
		body = bytes.NewReader(data)
	}
	status, content, err := HttpRequest(method, url, header, body)
	if err != nil {
		return
	}
	if status < 200 || status > 299 {
		err = HttpError(content)
	} else {
		err = json.Unmarshal([]byte(content), &rspobj)
	}
	return
}

func HttpProxy(rurl string, writer http.ResponseWriter, request *http.Request) (err error) {
	purl, err := url.Parse(rurl)
	if err == nil {
		request.Header.Set(REVERSE_SCHEME, purl.Scheme)
		request.Header.Set(REVERSE_HOST, purl.Host)
		request.Header.Set(REVERSE_PATH, purl.Path)
		ReverseProxy.ServeHTTP(writer, request)
	} else {
		writer.WriteHeader(http.StatusBadGateway)
		writer.Write([]byte(err.Error()))
	}
	return
}

func HttpProxyHandler(rurl string) *httputil.ReverseProxy {
	purl, _ := url.Parse(rurl)
	return &httputil.ReverseProxy{
		Transport:     HttpTransport,
		FlushInterval: ProxyFlushInterval,
		Director: func(req *http.Request) {
			req.URL.Scheme = purl.Scheme
			req.URL.Host = purl.Host
			req.URL.Path = purl.Path
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
		},
		BufferPool:   ReverseProxy.BufferPool,
		ErrorHandler: ReverseProxy.ErrorHandler,
	}
}
