package core

import resp "restrocheck/pkg/response"

type SavePriceResponse struct {
	resp.Response
	ID int64 `json:"id"`
}