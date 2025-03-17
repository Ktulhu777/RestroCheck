package core

import resp "restrocheck/pkg/response"

type SaveMenuResponse struct {
	resp.Response
	ID int64 `json:"id"`
}