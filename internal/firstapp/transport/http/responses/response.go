package responses

import "github.com/gofiber/fiber/v3"

type ResponseWrapper struct {
	c          fiber.Ctx
	statusCode int
	response   *Response
}

type Response struct {
	Success bool        `json:"success"`           // Always false for error
	Message string      `json:"message"`           // A user-friendly error message
	Data    interface{} `json:"data,omitempty"`    // the actual data that send
	Details interface{} `json:"details,omitempty"` // Optional: more specific error details (e.g., validation errors)
}

func NewResponse(c fiber.Ctx) *ResponseWrapper {
	return &ResponseWrapper{
		c:        c,
		response: &Response{},
	}
}

func (r *ResponseWrapper) Status(status int) *ResponseWrapper {
	r.statusCode = status
	return r
}

func (r *ResponseWrapper) Success(data interface{}) error {
	r.response.Success = true
	r.response.Data = data
	return r.c.Status(r.statusCode).JSON(r.response)
}

func (r *ResponseWrapper) Error(details interface{}) error {
	r.response.Success = false
	r.response.Details = details
	return r.c.Status(r.statusCode).JSON(r.response)
}
