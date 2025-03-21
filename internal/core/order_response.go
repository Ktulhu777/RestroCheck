package core

import resp "restrocheck/pkg/response"

type SaveOrderResponse struct {
	resp.Response
	ID int64 `json:"id"`
}
