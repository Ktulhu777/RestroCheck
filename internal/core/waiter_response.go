package core

import resp "restrocheck/pkg/response"

type SaveWaiterResponse struct {
	// TODO: поменять потом на SaveWaiterResponse
	resp.Response
	ID int64 `json:"id"`
}

type FetchWaiterResponse struct {
	resp.Response
	Waiter *Waiter `json:"waiter"`
}

type RemoveWaiterResponse struct {
	resp.Response
	ID int64 `json:"id"`
}

type ChangeWaiterResponse struct {
	resp.Response
	Waiter *Waiter `json:"waiter"`
}

type PartialWaiter struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
type FetchAllWaiterResponse struct {
	resp.Response
	Waiters []PartialWaiter `json:"waiters"`
}
