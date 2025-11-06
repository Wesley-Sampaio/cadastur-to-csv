package cadastur

import (
	"context"
	"encoding/json"
)

// Service provides methods to interact with the Cadastur API.
type Service struct {
	client *Client
}

// NewService creates a new Service instance.
func NewService() *Service {
	return &Service{
		client: NewClient(),
	}
}

// FetchUFs retrieves the list of UFs (states) from Cadastur.
// It returns a slice of UF or an error on network/parse failures.
func (s *Service) FetchUFs(ctx context.Context) ([]UF, error) {
	body, err := s.client.Get(ctx, EndpointUFs)
	if err != nil {
		return nil, err
	}

	var ufs []UF
	if err := json.Unmarshal(body, &ufs); err != nil {
		return nil, err
	}

	return ufs, nil
}

// FetchActivities retrieves the list of tourism activities from Cadastur.
func (s *Service) FetchActivities(ctx context.Context) ([]Activity, error) {
	body, err := s.client.Get(ctx, EndpointActivities)
	if err != nil {
		return nil, err
	}

	var acts []Activity
	if err := json.Unmarshal(body, &acts); err != nil {
		return nil, err
	}

	return acts, nil
}

// FetchPrestadoresPaged fetches providers data with pagination.
// It calls onPage callback for each page of results.
// Continues fetching while the page size is full (len(List) >= pageSize).
func (s *Service) FetchPrestadoresPaged(ctx context.Context, filters Filtros, pageSize int, onPage func([]Prestador, int, int) error) error {
	for currentPage := 1; ; currentPage++ {
		// Create and POST the request body for the current page.
		body := RequestBody{
			CurrentPage:    currentPage,
			PageSize:       pageSize,
			SortFields:     "nomePrestador",
			SortDirections: "ASC",
			Filtros:        filters,
		}

		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}

		respBody, err := s.client.Post(ctx, EndpointPrestadores, payload)
		if err != nil {
			return err
		}

		var out Response
		if err := json.Unmarshal(respBody, &out); err != nil {
			return err
		}

		// Call the callback with the current page's providers, page number, and total results
		if err := onPage(out.List, currentPage, out.TotalResults); err != nil {
			return err
		}

		// Stop if we got less than pageSize results (last page)
		if len(out.List) < pageSize {
			break
		}
	}

	return nil
}

// BuildFilters creates a Filtros struct with the provided parameters.
// Note: Localidade is hardcoded to 8452 as per original implementation.
func BuildFilters(selectedUF int, selectedActName string, localidadesUfs string) Filtros {
	return Filtros{
		NoPrestador:          "",
		Localidade:           8452, // Hardcoded value from original implementation
		NuAtividadeTuristica: selectedActName,
		SouPrestador:         false,
		SouTurista:           true,
		LocalidadesUfs:       localidadesUfs,
		LocalidadeNuUf:       selectedUF,
		FlPossuiVeiculo:      "",
	}
}

