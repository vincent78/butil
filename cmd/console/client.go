package console

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/model"
)

type ApiHandleFunc func(ctx context.Context, ch chan string, params ...string) model.RespModel

type Client struct {
	apiHandler map[string]ApiHandleFunc
	printer    io.Writer
}

func NewClient(printer io.Writer) *Client {
	return &Client{
		apiHandler: map[string]ApiHandleFunc{},
		printer:    printer,
	}
}

func (c *Client) getKey(apiName string) string {
	if strings.Index(apiName, ".") > 0 {
		return strings.ToLower(apiName)
	} else {
		return strings.ToLower(fmt.Sprintf("%v.%v", "", apiName))
	}

}

func (c *Client) RegisterApiHandler(api, name string, f ApiHandleFunc) {
	c.apiHandler[c.getKey(fmt.Sprintf("%v.%v", api, name))] = f
}

func (c *Client) DoHandler(apName string, args ...string) model.RespModel {
	key := c.getKey(apName)
	logger1.Info("handler: %v - %v", key, args)
	if f, exist := c.apiHandler[key]; exist && f != nil {
		ch := make(chan string)
		ctx := context.Background()
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-ch:
					if c.printer != nil {
						str := fmt.Sprintf("%v%v", msg, string(Newline))
						fmt.Fprintf(c.printer, str, "", "")
						logger1.Info("cmd[%v]: %v", key, str)
					}
				}
			}
		}(ctx)
		m := f(ctx, ch, args...)
		logger1.Info("cmd[%v] result: %v", key, m)
		if m.Code == model.SUCCESS {
			if m.Data != nil {
				ch <- fmt.Sprintf("%v", m.Data)
			}
		} else {
			ch <- m.String()
		}
		ctx.Done()
		return m
	} else {
		m := ErrorNoHandler(key)
		logger1.Info("cmd[%v] result: %v", key, m)
		return model.RespWithErrModel(m)
	}
}
