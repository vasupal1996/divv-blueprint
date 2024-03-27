package router

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type ErrorResp struct {
	ErrCode  string `json:"code,omitempty"`
	ErrMsg   string `json:"msg,omitempty"`
	ErrField string `json:"field,omitempty"`
}

func (er ErrorResp) Error() string {
	return er.ErrMsg
}

type Response struct {
	Success bool        `json:"success"`
	Payload interface{} `json:"payload,omitempty"`
	Error   []ErrorResp `json:"error,omitempty"`
}

func (r *Response) ToJSON() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func NewErrResponse(success bool, err ...ErrorResp) *Response {
	return &Response{
		Success: success,
		Error:   err,
	}
}

func NewErr(code, msg string) ErrorResp {
	return ErrorResp{
		ErrCode: code,
		ErrMsg:  msg,
	}
}

func NewJSONResp(success bool, payload interface{}) *Response {
	return &Response{
		Success: success,
		Payload: payload,
	}
}

func DecodeJSONBody(c *fiber.Ctx, dst interface{}) error {
	if c.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json"
		return NewErr("StatusUnsupportedMediaType", msg)
	}

	dec := json.NewDecoder(bytes.NewReader(c.Body()))
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf(
				"Request body contains badly-formed JSON (at position %d)",
				syntaxError.Offset,
			)
			return NewErr("StatusBadRequest", msg)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return NewErr("StatusBadRequest", "Request body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf(
				"Request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field,
				unmarshalTypeError.Offset,
			)
			return NewErr("StatusBadRequest", msg)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return NewErr("StatusBadRequest", msg)

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return NewErr("StatusBadRequest", msg)

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return NewErr("StatusRequestEntityTooLarge", msg)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return NewErr("StatusBadRequest", msg)
	}

	return nil
}
