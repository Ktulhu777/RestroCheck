package core

import resp "restrocheck/pkg/response"

type SaveCategoryResponse struct {
	resp.Response
	ID int64 `json:"id"`
}