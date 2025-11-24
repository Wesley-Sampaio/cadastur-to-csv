package cadastur

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
	// natJuridEspecial can come as string or boolean from the API; accept any to avoid unmarshal errors.
	NatJuridEspecial     any     `json:"natJuridEspecial"`
	Municipio            string  `json:"municipio"`
	NuMunicipio          int     `json:"nuMunicipio"`
	FlPossuiVeiculo      bool    `json:"flPossuiVeiculo"`
	NuSitCadTramite      int     `json:"nuSitCadTramite"`
	AtividadeRedeSociais *string `json:"atividadeRedeSociais"`
}
