package graceful

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/eapache/go-resiliency/breaker"
)

type Client struct {
	url     url.URL
	codec   Codec
	breaker *breaker.Breaker
	do      func(*http.Request) (*http.Response, error)
}

func NewClient(baseurl string, codec Codec, breaker *breaker.Breaker) (*Client, error) {
	if !strings.HasSuffix(baseurl, "/") {
		baseurl += "/"
	}
	url, err := url.Parse(baseurl)
	if err != nil {
		return nil, err
	}
	c := &Client{
		url:     *url,
		codec:   codec,
		breaker: breaker,
		do:      http.DefaultClient.Do,
	}
	return c, nil
}

func (c *Client) Call(ctx context.Context, method string, args, reply interface{}) error {
	if c.breaker == nil {
		return c.call(ctx, method, args, reply)
	}

	return c.breaker.Run(func() error {
		return c.call(ctx, method, args, reply)
	})
}

func (c *Client) call(ctx context.Context, method string, args, reply interface{}) error {
	var data []byte
	var err error
	if args != nil {
		data, err = c.codec.Marshal(args)
		if err != nil {
			return err
		}
	}
	// URL
	url := c.url
	url.Path += method
	// http.Request
	req := &http.Request{}
	req = req.WithContext(ctx)
	req.URL = &url
	if args != nil {
		req.Method = "POST"
		req.Body = ioutil.NopCloser(bytes.NewReader(data))
		req.ContentLength = int64(len(data))
		// req.Header.Set("Content-Type", c.codec.MIME())
	} else {
		req.Method = "GET"
	}
	// start request
	return c.request(reply, req)
}

func (c *Client) request(reply interface{}, req *http.Request) error {
	// c.timeoutToHeader(req)
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = resp.Body.Close(); err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		for len(data) != 0 && data[len(data)-1] == '\n' {
			data = data[:len(data)-1]
		}
		err = errors.New(string(data))
		return err
	}
	err = c.codec.Unmarshal(data, reply)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SetRequestProcessor(p func(*http.Request) (*http.Response, error)) {
	c.do = p
}

// timeoutToHeader write context timeout to http-header
func (c *Client) timeoutToHeader(req *http.Request) {
	ctx := req.Context()
	deadline, ok := ctx.Deadline()
	if ok {
		timeout := -1 * time.Since(deadline)
		str := strconv.FormatInt(int64(timeout), 16)
		if req.Header == nil {
			req.Header = make(http.Header, 1)
		}
		req.Header.Add(contextTimeoutHeader, str)
	}
}
