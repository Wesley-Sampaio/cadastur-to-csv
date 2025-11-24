package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"cadastur-csv/internal/cadastur"
)

// section prints a simple visual separator with a title for better CLI UX.
func section(title string) {
	bar := strings.Repeat("─", 60)
	fmt.Printf("\n%s\n%s\n%s\n", bar, title, bar)
}

// PromptUF displays available UFs and prompts the user to select one.
// Returns the selected UF ID, defaulting to 24 (Santa Catarina) if input is invalid.
func PromptUF(ufs []cadastur.UF) (int, error) {
	section("Selecione um Estado (UF)")
	fmt.Println("UFs disponíveis:")
	for _, uf := range ufs {
		fmt.Printf("%d - %s (%s)\n", uf.ID, uf.NoUf, uf.SgUf)
	}

	// Ask the user to choose a UF using line-based input to avoid scanf issues
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("▶ Digite o ID da UF (padrão 24 para Santa Catarina): ")
	line, err := reader.ReadString('\n')
	if err != nil {
		// On read error, fall back to default
		return 24, nil
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return 24, nil
	}
	v, err := strconv.Atoi(line)
	if err != nil {
		return 24, nil
	}
	fmt.Println("UF selecionada:", v)
	return v, nil
}

// PromptActivity displays available activities and prompts the user to select one.
// Returns the selected activity ID and name, defaulting to 29 (Guia de Turismo) if input is invalid.
func PromptActivity(activities []cadastur.Activity) (int, string, error) {
	section("Selecione uma Atividade Turística")
	fmt.Println("Atividades disponíveis (somente ativas):")
	for _, a := range activities {
		if a.FlAtivo {
			fmt.Printf("%d - %s\n", a.NuAtividadeTuristica, a.NoAtividadeTuristica)
		}
	}

	// Ask the user to choose an activity (by ID) using line-based input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("▶ Digite o ID da atividade (padrão 29 para Guia de Turismo): ")
	line, err := reader.ReadString('\n')
	if err != nil {
		// fallback to default
		return 29, "Guia de Turismo", nil
	}
	line = strings.TrimSpace(line)
	var selectedActID int
	if line == "" {
		selectedActID = 29
	} else {
		v, err := strconv.Atoi(line)
		if err != nil {
			selectedActID = 29
		} else {
			selectedActID = v
		}
	}

	// Map ID -> Name (fallback to 'Guia de Turismo' if not found)
	selectedActName := "Guia de Turismo"
	for _, a := range activities {
		if a.NuAtividadeTuristica == selectedActID {
			selectedActName = a.NoAtividadeTuristica
			break
		}
	}
	fmt.Println("Atividade selecionada:", selectedActID, "-", selectedActName)

	return selectedActID, selectedActName, nil
}

// PromptCity prompts the user for an optional city input.
// Returns the city string (can be empty), with a Portuguese warning about leaving it blank.
func PromptCity() (string, error) {
	section("Opcional: localidadesUfs")
	// Ask for optional localidadesUfs (free text; can be empty)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("▶ Digite a cidade (ex: \"Florianópolis, SC\") ou pressione ENTER para deixar em branco.\n  Atenção: preencher incorretamente pode reduzir os resultados. Recomenda-se deixar vazio na primeira tentativa.\n  Cidade: ")
	localidadesUfs, _ := reader.ReadString('\n')
	localidadesUfs = strings.TrimSpace(localidadesUfs)
	if localidadesUfs == "" {
		fmt.Println("Nenhuma cidade informada. (Recomendado para obter o maior número de resultados.)")
	} else {
		fmt.Println("Cidade selecionada:", localidadesUfs)
	}

	return localidadesUfs, nil
}
