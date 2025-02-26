package core

import resp "restrocheck/pkg/response"

type SaveResponse struct {
	resp.Response
	ID int64 `json:"id"`
}

type FetchResponse struct {
	resp.Response
	Waiter *Waiter `json:"waiter"`
}

type RemoveResponse struct {
	resp.Response
	ID int64 `json:"id"`
}

type ChangeResponse struct {
	resp.Response
	Waiter *Waiter `json:"waiter"`
}

type PartialWaiter struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
type FetchAllResponse struct {
	resp.Response
	Waiters []PartialWaiter `json:"waiters"`
}
