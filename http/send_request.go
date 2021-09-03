package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (h *LuaModuleHTTP) sendRequest(args *requestArgs) (*response, error) {
	h.logger.Debug("http request", zap.Any("args", args))

	req, err := http.NewRequest(args.Method, args.URL, bytes.NewReader(args.Body))
	if err != nil {
		return nil, fmt.Errorf("error build request, %w", err)
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), args.Timeout)
	defer ctxCancel()

	req = req.WithContext(ctx)

	for name, value := range args.Headers {
		req.Header.Set(name, value)
	}

	span, _ := opentracing.StartSpanFromContext(h.ctx, "/")
	defer span.Finish()
	span.SetTag("http.method", args.Method)
	span.SetTag("span.kind", "client")
	err = opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		return nil, fmt.Errorf("inject: %s", err.Error())
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error send request, %w", err)
	}
	defer resp.Body.Close()

	span.SetTag("http.status_code", resp.StatusCode)

	res := newResponse()

	res.StatusCode = resp.StatusCode

	res.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read body, %w", err)
	}

	for name, values := range resp.Header {
		if len(values) > 1 {
			h.logger.Debug("the response header has multiple values", zap.String("name", name), zap.Strings("values", values))
		}
		for _, value := range values {
			res.Headers[name] = value
		}
	}

	return res, nil
}
