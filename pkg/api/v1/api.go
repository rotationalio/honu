package api

import "context"

//go:generate msgp

//===========================================================================
// Service Interface
//===========================================================================

// Client defines the service interface for interacting with the HonuDB service.
type Client interface {
	Status(context.Context) (*StatusReply, error)
}

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Reply contains standard fields that are used for generic API responses and errors.
type Reply struct {
	Success     bool        `json:"success" msg:"success"`
	Error       string      `json:"error,omitempty" msg:"error,omitempty"`
	ErrorDetail ErrorDetail `json:"errors,omitempty" msg:"errors,omitempty"`
}

// Returned on status requests.
type StatusReply struct {
	Status  string `json:"status" msg:"status"`
	Uptime  string `json:"uptime,omitempty" msg:"uptime,omitempty"`
	Version string `json:"version,omitempty" msg:"version,omitempty"`
}

// PageQuery manages paginated list requests.
type PageQuery struct {
	PageSize      int    `json:"page_size,omitempty" msg:"page_size,omitempty" url:"page_size,omitempty" form:"page_size"`
	NextPageToken string `json:"next_page_token,omitempty" msg:"next_page_token,omitempty" url:"next_page_token,omitempty" form:"next_page_token"`
	PrevPageToken string `json:"prev_page_token,omitempty" msg:"prev_page_token,omitempty" url:"prev_page_token,omitempty" form:"prev_page_token"`
}

//===========================================================================
// Detail Error
//===========================================================================

type ErrorDetail []*DetailError

type DetailError struct {
	Field string `json:"field" msg:"field"`
	Error string `json:"error" msg:"error"`
}
