package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"proxy/pkg/cast"
	"proxy/pkg/database"
	"proxy/pkg/server"
	"strconv"

	"github.com/aerogo/aero"
)

type Handler struct {
	db *database.DB
}

type ParsedRequest struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Params  map[string][]string `json:"query_params,omitempty"`
	Headers map[string][]string `json:"headers"`
	Cookies map[string]string   `json:"cookies,omitempty"`
	Body    map[string]string   `json:"body,omitempty"`
}

type ParsedResponse struct {
	Code    uint64              `json:"code"`
	Message string              `json:"message"`
	Body    string              `json:"body"`
	Headers map[string][]string `json:"headers"`
}

type Record struct {
	Request  *ParsedRequest  `json:"request"`
	Response *ParsedResponse `json:"response"`
}

func SelectRecord(h *Handler, id int) (*Record, error) {
	req, err := h.db.Select("requests", id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	r := &ParsedRequest{
		Method:  req[1].(string),
		Path:    req[2].(string),
		Params:  cast.ToMapStringArray(req[3].(map[interface{}]interface{})),
		Headers: cast.ToMapStringArray(req[4].(map[interface{}]interface{})),
		Cookies: cast.ToMapString(req[5].(map[interface{}]interface{})),
	}
	if req[6] != nil {
		r.Body = cast.ToMapString(req[6].(map[interface{}]interface{}))
	}

	resp, err := h.db.Select("responses", id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	w := &ParsedResponse{
		Code:    resp[2].(uint64),
		Message: resp[1].(string),
		Body:    resp[3].(string),
		Headers: cast.ToMapStringArray(resp[4].(map[interface{}]interface{})),
	}
	return &Record{r, w}, nil
}

func (h *Handler) RequestAll(ctx aero.Context) error {
	result := make([]*Record, 0)
	id := 1
	for {
		record, err := SelectRecord(h, id)
		if err != nil {
			break
		}
		result = append(result, record)
		id += 1
	}
	return ctx.JSON(result)
}

func (h *Handler) RequestOne(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		return ctx.Error(http.StatusBadRequest)
	}

	record, err := SelectRecord(h, id)
	if err != nil {
		ctx.Error(http.StatusInternalServerError)
	}
	return ctx.JSON(record)
}

func (h *Handler) Repeat(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		return ctx.Error(http.StatusBadRequest)
	}

	record, err := SelectRecord(h, id)
	if err != nil {
		ctx.Error(http.StatusInternalServerError)
	}
	req := &http.Request{
		Method: record.Request.Method,
		Header: record.Request.Headers,
	}
	for k, v := range record.Request.Cookies {
		req.AddCookie(&http.Cookie{
			Name:   k,
			Value:  v,
			Path:   "/",
			MaxAge: 900,
		})
	}
	req.URL, _ = url.Parse(record.Request.Path)

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return ctx.Error(http.StatusServiceUnavailable)
	}
	defer resp.Body.Close()

	ctx.SetStatus(resp.StatusCode)
	server.CopyHeader(ctx.Response().Internal().Header(), resp.Header)

	body := make([]byte, resp.ContentLength)
	buf := bytes.NewBuffer(body)
	io.Copy(buf, resp.Body)
	return ctx.Bytes(buf.Bytes())
}

func (h *Handler) Scan(ctx aero.Context) error {
	return nil
}
