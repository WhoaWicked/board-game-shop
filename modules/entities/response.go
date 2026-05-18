package entities

import (
	"github.com/WhoaWicked/board-game-shop/pkg/shoplogger"
	"github.com/gofiber/fiber/v3"
)

type IResponse interface {
	Success(code int, data any) IResponse
	Error(code int, tractId string, msg string) IResponse
	Res() error
}

type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	Context    fiber.Ctx
	IsError    bool
}

type ErrorResponse struct {
	TractId string `json:"tract_id"`
	Msg     string `json:"message"`
}

func NewResponse(c fiber.Ctx) IResponse {
	return &Response{
		Context: c,
	}
}

func (r *Response) Success(code int, data any) IResponse {
	r.StatusCode = code
	r.Data = data
	shoplogger.InitShopLogger(r.Context, r.Data).Print().Save()
	return r
}
func (r *Response) Error(code int, tractId string, msg string) IResponse {
	r.StatusCode = code
	r.ErrorRes = &ErrorResponse{
		TractId: tractId,
		Msg:     msg,
	}
	r.IsError = true
	shoplogger.InitShopLogger(r.Context, r.ErrorRes).Print()
	return r
}

func (r *Response) Res() error {
	return r.Context.Status(r.StatusCode).JSON(func() any {
		if r.IsError {
			return r.ErrorRes
		}
		return &r.Data
	}())
}
