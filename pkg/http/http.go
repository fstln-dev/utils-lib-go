package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
)

// func httpsClient(method, url string, headers map[string]string, request, response interface{}) error {

// 	jsonStr, err := json.Marshal(request)
// 	if err != nil {
// 		return err
// 	}

// 	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
// 	for k, v := range headers {
// 		req.Header.Set(k, v)
// 	}

// 	tr := &http.Transport{
// 		TLSClientConfig: &tls.Config{
// 			// 跳过客户端 HTTPS 验证
// 			InsecureSkipVerify: true,
// 		},
// 	}
// 	client := &http.Client{Transport: tr}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	defer client.CloseIdleConnections()

// 	if resp.Status != "200" {
// 		return errors.New("response error: " + resp.Status)
// 	}

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	return json.Unmarshal(body, response)
// }

func httpsClient(ctx context.Context, uri, method string, opt *options, response interface{}, buf io.Reader) error {

	req, err := http.NewRequestWithContext(ctx, method, uri, buf)
	if err != nil {
		return err
	}

	if opt.basicAuth != nil {
		req.SetBasicAuth(opt.basicAuth.username, opt.basicAuth.password)
	}

	for k, v := range opt.headers {
		req.Header.Set(k, v)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opt.insecure,
		},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer client.CloseIdleConnections()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(resp.StatusCode, resp.Status, string(body))
	}

	if response != nil {
		return json.Unmarshal(body, response)
	}
	return nil
}

type Option func(o *options)
type options struct {
	headers   map[string]string
	params    *url.Values
	request   interface{}
	basicAuth *basicAuth
	timeout   time.Duration
	insecure  bool
}

type basicAuth struct {
	username string
	password string
}

func Header(h map[string]string) Option {
	return func(o *options) {
		o.headers = h
	}
}

func Params(u *url.Values) Option {
	return func(o *options) {
		o.params = u
	}
}

func Request(r interface{}) Option {
	return func(o *options) {
		o.request = r
	}
}

func BasicAuth(username, password string) Option {
	return func(o *options) {
		o.basicAuth = &basicAuth{
			username: username,
			password: password,
		}
	}
}

func Timeout(t time.Duration) Option {
	return func(o *options) {
		o.timeout = t
	}
}

// Insecure 显式跳过 TLS 服务器证书验证（默认验证）。
// 仅在开发环境或自签名证书环境下使用，生产环境不建议。
func Insecure() Option {
	return func(o *options) {
		o.insecure = true
	}
}

func post(ctx context.Context, method, path string, response interface{}, opt *options) error {
	if _, ok := opt.headers["Content-Type"]; !ok {
		opt.headers["Content-Type"] = "application/json"
	}

	var buf io.Reader

	switch opt.headers["Content-Type"] {
	case "application/json":
		req, err := json.Marshal(opt.request)
		if err != nil {
			return err
		}
		buf = bytes.NewReader(req)
	case "application/x-www-form-urlencoded":
		uv := &url.Values{}
		err := StructToUrl(opt.request, uv)
		if err != nil {
			return err
		}
		buf = strings.NewReader(uv.Encode())
	default:
	}

	uri, err := url.Parse(path)
	if err != nil {
		return err
	}

	if opt.params != nil {
		uri.RawQuery = opt.params.Encode()
	}

	return httpsClient(ctx, uri.String(), method, opt, response, buf)
}

func HttpsPost(ctx context.Context, path string, response interface{}, opts ...Option) error {

	// 后续可校验 https:// 头

	opt := &options{
		headers: make(map[string]string),
	}

	for _, o := range opts {
		o(opt)
	}

	return post(ctx, http.MethodPost, path, response, opt)
}

func HttpsPut(ctx context.Context, path string, response interface{}, opts ...Option) error {

	// 后续可校验 https:// 头

	opt := &options{
		headers: make(map[string]string),
	}

	for _, o := range opts {
		o(opt)
	}

	return post(ctx, http.MethodPut, path, response, opt)
}

func HttpsPatch(ctx context.Context, path string, response interface{}, opts ...Option) error {

	opt := &options{
		headers: make(map[string]string),
	}

	for _, o := range opts {
		o(opt)
	}

	return post(ctx, http.MethodPatch, path, response, opt)
}

func get(ctx context.Context, method, path string, response interface{}, opt *options) error {
	uri, err := url.Parse(path)
	if err != nil {
		return err
	}

	if opt.params != nil {
		uri.RawQuery = opt.params.Encode()
	}

	return httpsClient(ctx, uri.String(), method, opt, response, nil)
}

func HttpsGet(ctx context.Context, path string, response interface{}, opts ...Option) error {

	// 后续可校验 https:// 头

	opt := &options{
		headers: make(map[string]string),
		params:  &url.Values{},
	}

	for _, o := range opts {
		o(opt)
	}

	return get(ctx, http.MethodGet, path, response, opt)
}

func HttpsDelete(ctx context.Context, path string, response interface{}, opts ...Option) error {

	// 后续可校验 https:// 头

	opt := &options{
		headers: make(map[string]string),
		params:  &url.Values{},
	}

	for _, o := range opts {
		o(opt)
	}

	return get(ctx, http.MethodDelete, path, response, opt)
}
