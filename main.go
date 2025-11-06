// cadastur-csv
// EN: A command-line tool that fetches UFs and tourism activities from Cadastur,
// asks the user for filters (UF, activity, optional city), paginates through
// "obterDadosPrestadores", normalizes the data (phones, CEPs), and exports a CSV.
// No external dependencies — uses only Go's standard library.
//
// PT: Uma ferramenta de linha de comando que busca UFs e atividades turísticas
// no Cadastur, solicita filtros ao usuário (UF, atividade, cidade opcional),
// pagina os resultados de "obterDadosPrestadores", normaliza os dados
// (telefones, CEPs) e exporta tudo em um CSV.
// Sem dependências externas — utiliza apenas a biblioteca padrão do Go.
//
// Public repo friendly: comments in English, user prompts in Portuguese.
// Author: Raphael Ornellas (with help from ChatGPT assistant)
// License: MIT (suggested)
package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
)

// UF represents a Federative Unit (state) returned by Cadastur `tipoUfs`.
type UF struct {
	ID   int    `json:"id"`
	NoUf string `json:"noUf"`
	SgUf string `json:"sgUf"`
}

// Activity represents a tourism activity entry from `atividadesTuristica`.
type Activity struct {
	NuAtividadeTuristica   int    `json:"nuAtividadeTuristica"`
	NoAtividadeTuristica   string `json:"noAtividadeTuristica"`
	FlAtividadeObrigatoria bool   `json:"flAtividadeObrigatoria"`
	FlAtivo                bool   `json:"flAtivo"`
}

// Filtros matches the request body's "filtros" expected by `obterDadosPrestadores`.
type Filtros struct {
	NoPrestador          string `json:"noPrestador"`
	Localidade           int    `json:"localidade"`
	NuAtividadeTuristica string `json:"nuAtividadeTuristica"`
	SouPrestador         bool   `json:"souPrestador"`
	SouTurista           bool   `json:"souTurista"`
	LocalidadesUfs       string `json:"localidadesUfs"`
	LocalidadeNuUf       int    `json:"localidadeNuUf"`
	FlPossuiVeiculo      string `json:"flPossuiVeiculo"`
}

// RequestBody is the POST payload for `obterDadosPrestadores`.
type RequestBody struct {
	CurrentPage    int     `json:"currentPage"`
	PageSize       int     `json:"pageSize"`
	SortFields     string  `json:"sortFields"`
	SortDirections string  `json:"sortDirections"`
	Filtros        Filtros `json:"filtros"`
}

// Response models the paginated response from `obterDadosPrestadores`.
type Response struct {
	CurrentPage    int         `json:"currentPage"`
	PageSize       int         `json:"pageSize"`
	TotalResults   int         `json:"totalResults"`
	SortFields     string      `json:"sortFields"`
	SortDirections string      `json:"sortDirections"`
	Filtros        Filtros     `json:"filtros"`
	List           []Prestador `json:"list"`
	Start          int         `json:"start"`
}

// Prestador represents one provider row returned in `Response.List`.
type Prestador struct {
	ID                   int     `json:"id"`
	TipoPessoa           string  `json:"tipoPessoa"`
	NumeroCadastro       string  `json:"numeroCadastro"`
	DtInicioVigencia     int64   `json:"dtInicioVigencia"`
	DtFimVigencia        int64   `json:"dtFimVigencia"`
	NoWebSite            *string `json:"noWebSite"`
	NuTelefone           string  `json:"nuTelefone"`
	NoLogradouro         string  `json:"noLogradouro"`
	Complemento          string  `json:"complemento"`
	NuCep                string  `json:"nuCep"`
	Sguf                 string  `json:"sguf"`
	NoBairro             string  `json:"noBairro"`
	NomePrestador        string  `json:"nomePrestador"`
	RegistroRf           string  `json:"registroRf"`
	NuAtividadeTuristica int     `json:"nuAtividadeTuristica"`
	Atividade            string  `json:"atividade"`
	NuSituacaoCadastral  int     `json:"nuSituacaoCadastral"`
	Situacao             string  `json:"situacao"`
	NuUf                 int     `json:"nuUf"`
	LocalidadeNuUf       *int    `json:"localidadeNuUf"`
	Localidade           string  `json:"localidade"`
	NoLocalidade         string  `json:"noLocalidade"`
	NuLocalidade         int     `json:"nuLocalidade"`
	NuPessoa             int     `json:"nuPessoa"`
	NatJuridEspecial     *string `json:"natJuridEspecial"`
	Municipio            string  `json:"municipio"`
	NuMunicipio          int     `json:"nuMunicipio"`
	FlPossuiVeiculo      bool    `json:"flPossuiVeiculo"`
	NuSitCadTramite      int     `json:"nuSitCadTramite"`
	AtividadeRedeSociais *string `json:"atividadeRedeSociais"`
}

// fetchUFs retrieves the list of UFs (states) from Cadastur.
// It returns a slice of UF or an error on network/parse failures.
func fetchUFs() ([]UF, error) {
	url := "https://cadastur.turismo.gov.br/cadastur-backend/rest/dominios/tipoUfs"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ufs []UF
	if err := json.Unmarshal(body, &ufs); err != nil {
		return nil, err
	}

	return ufs, nil
}

// fetchActivities retrieves the list of tourism activities from Cadastur.
// Only active activities are displayed to the user in the prompt.
func fetchActivities() ([]Activity, error) {
	url := "https://cadastur.turismo.gov.br/cadastur-backend/rest/portal/atividadesTuristica"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var acts []Activity
	if err := json.Unmarshal(body, &acts); err != nil {
		return nil, err
	}
	return acts, nil
}

// section prints a simple visual separator with a title for better CLI UX.
func section(title string) {
	bar := strings.Repeat("─", 60)
	fmt.Printf("\n%s\n%s\n%s\n", bar, title, bar)
}

// slugify turns a human-readable string into a safe filename fragment:
// lowercased, spaces -> dashes, and only keeps letters/digits/-/_.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// keep only letters, digits, '-' and '_'
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			out = append(out, r)
		}
	}
	// remove consecutive dashes
	clean := make([]rune, 0, len(out))
	var prevDash bool
	for _, r := range out {
		if r == '-' {
			if !prevDash {
				clean = append(clean, r)
				prevDash = true
			}
			continue
		}
		prevDash = false
		clean = append(clean, r)
	}
	if len(clean) == 0 {
		return "atividade"
	}
	return string(clean)
}

// onlyDigits returns a version of s that contains digits only.
// Used to normalize phone numbers and CEP fields in the CSV.
func onlyDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func main() {
	ufs, err := fetchUFs()
	// 1) Fetch UFs and prompt the user to select a state (with default).
	if err != nil {
		panic(err)
	}

	section("Selecione um Estado (UF)")
	fmt.Println("UFs disponíveis:")
	for _, uf := range ufs {
		fmt.Printf("%d - %s (%s)\n", uf.ID, uf.NoUf, uf.SgUf)
	}

	// Ask the user to choose a UF
	var selectedUF int
	fmt.Print("▶ Digite o ID da UF (padrão 24 para Santa Catarina): ")
	n, err := fmt.Scanf("%d", &selectedUF)

	if n != 1 || err != nil {
		selectedUF = 24
	}

	fmt.Println("UF selecionada:", selectedUF)

	// 2) Fetch activities and prompt the user to select one (with default).
	acts, err := fetchActivities()
	if err != nil {
		panic(err)
	}
	section("Selecione uma Atividade Turística")
	fmt.Println("Atividades disponíveis (somente ativas):")
	for _, a := range acts {
		if a.FlAtivo {
			fmt.Printf("%d - %s\n", a.NuAtividadeTuristica, a.NoAtividadeTuristica)
		}
	}

	// Ask the user to choose an activity (by ID)
	var selectedActID int
	fmt.Print("▶ Digite o ID da atividade (padrão 29 para Guia de Turismo): ")
	n2, err := fmt.Scanf("%d", &selectedActID)
	if n2 != 1 || err != nil {
		selectedActID = 29
	}

	// Map ID -> Name (fallback to 'Guia de Turismo' if not found)
	selectedActName := "Guia de Turismo"
	for _, a := range acts {
		if a.NuAtividadeTuristica == selectedActID {
			selectedActName = a.NoAtividadeTuristica
			break
		}
	}
	fmt.Println("Atividade selecionada:", selectedActID, "-", selectedActName)

	// Build a descriptive CSV filename based on the chosen activity.
	fileName := fmt.Sprintf("prestadores-atividade-%d-%s.csv", selectedActID, slugify(selectedActName))

	// 3) Prompt for optional city (free-text). Leaving it blank is recommended for broader results.
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

	url := "https://cadastur.turismo.gov.br/cadastur-backend/rest/portal/obterDadosPrestadores"
	client := &http.Client{Timeout: 30 * time.Second}

	// Prepare CSV writer (single file) — header is written once.
	section(fmt.Sprintf("Salvando CSV em %s", fileName))
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Header
	header := []string{
		"id", "tipoPessoa", "numeroCadastro",
		"inicioVigencia", "fimVigencia",
		"website", "telefone",
		"logradouro", "complemento", "cep", "uf", "bairro",
		"nomePrestador", "registroRf",
		"nuAtividadeTuristica", "atividade",
		"nuSituacaoCadastral", "situacao",
		"nuUf", "localidadeNuUf",
		"municipio", "localidade", "noLocalidade",
		"nuLocalidade", "nuMunicipio", "nuPessoa",
		"possuiVeiculo", "nuSitCadTramite", "atividadeRedeSociais",
	}
	if err := w.Write(header); err != nil {
		panic(err)
	}

	// 4) Pagination loop — keep fetching pages until the last page is smaller than pageSize.
	pageSize := 1000
	totalFetched := 0
	pages := 0
	totalExpected := -1

	// Keep a few samples to show after
	samples := make([]Prestador, 0, 5)

	for currentPage := 1; ; currentPage++ {
		// Create and POST the request body for the current page.
		body := RequestBody{
			CurrentPage:    currentPage,
			PageSize:       pageSize,
			SortFields:     "nomePrestador",
			SortDirections: "ASC",
			Filtros: Filtros{
				NoPrestador:          "",
				Localidade:           8452,
				NuAtividadeTuristica: selectedActName,
				SouPrestador:         false,
				SouTurista:           true,
				LocalidadesUfs:       localidadesUfs,
				LocalidadeNuUf:       selectedUF,
				FlPossuiVeiculo:      "",
			},
		}

		payload, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if err != nil {
			panic(err)
		}
		req.Header.Set("accept", "application/json, text/plain, */*")
		req.Header.Set("content-type", "application/json;charset=UTF-8")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		resBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			panic(err)
		}

		var out Response
		if err := json.Unmarshal(resBytes, &out); err != nil {
			panic(err)
		}

		if totalExpected == -1 {
			totalExpected = out.TotalResults
		}

		// Status per page
		fmt.Printf("Página %d — HTTP: %s — recebidos: %d\n", currentPage, resp.Status, len(out.List))

		// Append each provider as one CSV row, normalizing phone/CEP.
		for _, p := range out.List {
			if len(samples) < 5 {
				samples = append(samples, p)
			}
			row := []string{
				fmt.Sprint(p.ID),
				p.TipoPessoa,
				p.NumeroCadastro,
				msToDate(p.DtInicioVigencia),
				msToDate(p.DtFimVigencia),
				emptyIfNil(p.NoWebSite),
				onlyDigits(p.NuTelefone),
				p.NoLogradouro,
				p.Complemento,
				onlyDigits(p.NuCep),
				p.Sguf,
				p.NoBairro,
				p.NomePrestador,
				p.RegistroRf,
				fmt.Sprint(p.NuAtividadeTuristica),
				p.Atividade,
				fmt.Sprint(p.NuSituacaoCadastral),
				p.Situacao,
				fmt.Sprint(p.NuUf),
				intPtrToStr(p.LocalidadeNuUf),
				p.Municipio,
				p.Localidade,
				p.NoLocalidade,
				fmt.Sprint(p.NuLocalidade),
				fmt.Sprint(p.NuMunicipio),
				fmt.Sprint(p.NuPessoa),
				boolToStr(p.FlPossuiVeiculo),
				fmt.Sprint(p.NuSitCadTramite),
				emptyIfNil(p.AtividadeRedeSociais),
			}
			if err := w.Write(row); err != nil {
				panic(err)
			}
		}

		w.Flush()
		if err := w.Error(); err != nil {
			panic(err)
		}

		totalFetched += len(out.List)
		pages++

		// Stop conditions: last page or server returned less than pageSize
		if len(out.List) < pageSize {
			break
		}
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
}

// msToDate converts a millisecond Unix timestamp into YYYY-MM-DD.
// Returns an empty string when ms == 0.
func msToDate(ms int64) string {
	if ms == 0 {
		return ""
	}
	t := time.Unix(0, ms*int64(time.Millisecond))
	return t.Format("2006-01-02")
}

// emptyIfNil safely dereferences optional string pointers for CSV output.
func emptyIfNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// boolToStr converts a boolean to "true"/"false" for CSV cells.
func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// intPtrToStr renders optional integer pointers as strings ("" when nil).
func intPtrToStr(v *int) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(*v)
}
