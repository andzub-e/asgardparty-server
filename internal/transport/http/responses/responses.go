package responses

import "bitbucket.org/electronicjaw/asgardparty-server/internal/entities"

type StateResponse struct {
	entities.SessionState
}

type GetFreeSpinsResponse struct {
	FreeSpins []*entities.FreeSpin `json:"freespins"`
}

type GetFreeSpinsWithIntegratorBetResponse struct {
	FreeSpins map[string][]*entities.FreeSpin `json:"freespins"`
}

type HealthResponse struct {
	Success string `json:"success"`
}

type InfoResponse struct {
	Tag string `json:"tag"`
}

type NoContentResponse struct {
	Success bool `json:"success"`
}
