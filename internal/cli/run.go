package cli

import (
	"context"
	"fmt"
	"net/http"

	"cadastur-csv/internal/cadastur"
	"cadastur-csv/internal/csvx"
	"cadastur-csv/internal/normalize"
)

// Run orchestrates the full workflow: prompts → API → CSV writer → summary.
func Run(ctx context.Context, service *cadastur.Service) error {
	// 1) Fetch UFs and prompt the user to select a state (with default).
	ufs, err := service.FetchUFs(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch UFs: %w", err)
	}

	selectedUF, err := PromptUF(ufs)
	if err != nil {
		return fmt.Errorf("failed to prompt UF: %w", err)
	}

	// 2) Fetch activities and prompt the user to select one (with default).
	acts, err := service.FetchActivities(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch activities: %w", err)
	}

	selectedActID, selectedActName, err := PromptActivity(acts)
	if err != nil {
		return fmt.Errorf("failed to prompt activity: %w", err)
	}

	// Build a descriptive CSV filename based on the chosen activity.
	fileName := fmt.Sprintf("prestadores-atividade-%d-%s.csv", selectedActID, normalize.Slugify(selectedActName))

	// 3) Prompt for optional city (free-text). Leaving it blank is recommended for broader results.
	localidadesUfs, err := PromptCity()
	if err != nil {
		return fmt.Errorf("failed to prompt city: %w", err)
	}

	// Prepare CSV writer (single file) — header is written once.
	section(fmt.Sprintf("Salvando CSV em %s", fileName))
	csvWriter, err := csvx.NewWriter(fileName)
	if err != nil {
		return fmt.Errorf("failed to create CSV writer: %w", err)
	}
	defer csvWriter.Close()

	// Write header
	if err := csvWriter.WriteHeader(); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// 4) Pagination loop — keep fetching pages until the last page is smaller than pageSize.
	pageSize := 1000
	totalFetched := 0
	pages := 0
	totalExpected := -1

	// Keep a few samples to show after
	samples := make([]cadastur.Prestador, 0, 5)

	// Build filters
	filters := cadastur.BuildFilters(selectedUF, selectedActName, localidadesUfs)

	// Fetch all pages
	err = service.FetchPrestadoresPaged(ctx, filters, pageSize, func(prestadores []cadastur.Prestador, currentPage, totalResults int) error {
		if totalExpected == -1 {
			totalExpected = totalResults
		}

		// Status per page (simulating HTTP status for consistency with original)
		fmt.Printf("Página %d — HTTP: %s — recebidos: %d\n", currentPage, http.StatusText(http.StatusOK), len(prestadores))

		// Append each provider as one CSV row, normalizing phone/CEP.
		for _, p := range prestadores {
			if len(samples) < 5 {
				samples = append(samples, p)
			}
			if err := csvWriter.WriteRow(p); err != nil {
				return fmt.Errorf("failed to write CSV row: %w", err)
			}
		}

		// Flush after each page
		if err := csvWriter.Flush(); err != nil {
			return fmt.Errorf("failed to flush CSV: %w", err)
		}

		totalFetched += len(prestadores)
		pages++

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to fetch prestadores: %w", err)
	}

	// 5) Final summary and a small sample for visual verification in the terminal.
	section("Resumo")
	fmt.Printf("Total de resultados: %d | Páginas: %d | Retornados: %d\n", totalExpected, pages, totalFetched)

	// Show first samples
	max := len(samples)
	for i := 0; i < max; i++ {
		p := samples[i]
		fmt.Printf("%d) %s | %s | %s\n", i+1, p.NomePrestador, p.Municipio, p.NuTelefone)
	}

	return nil
}

