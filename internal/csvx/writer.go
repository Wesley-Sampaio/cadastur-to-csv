package csvx

import (
	"encoding/csv"
	"fmt"
	"os"

	"cadastur-csv/internal/cadastur"
	"cadastur-csv/internal/normalize"
)

// Writer handles CSV file writing with header and row normalization.
type Writer struct {
	file   *os.File
	writer *csv.Writer
}

// NewWriter creates a new CSV writer for the specified filename.
func NewWriter(filename string) (*Writer, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	w := csv.NewWriter(f)

	return &Writer{
		file:   f,
		writer: w,
	}, nil
}

// WriteHeader writes the CSV header once.
// Header must match exactly: id,tipoPessoa,numeroCadastro,inicioVigencia,fimVigencia,website,telefone,logradouro,complemento,cep,uf,bairro,nomePrestador,registroRf,nuAtividadeTuristica,atividade,nuSituacaoCadastral,situacao,nuUf,localidadeNuUf,municipio,localidade,noLocalidade,nuLocalidade,nuMunicipio,nuPessoa,possuiVeiculo,nuSitCadTramite,atividadeRedeSociais
func (w *Writer) WriteHeader() error {
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
	return w.writer.Write(header)
}

// WriteRow writes a normalized row for a Prestador.
// Normalizes telephone and CEP to digits only, handles dates, bools, and pointers.
func (w *Writer) WriteRow(p cadastur.Prestador) error {
	row := []string{
		fmt.Sprint(p.ID),
		p.TipoPessoa,
		p.NumeroCadastro,
		normalize.MsToDate(p.DtInicioVigencia),
		normalize.MsToDate(p.DtFimVigencia),
		normalize.EmptyIfNil(p.NoWebSite),
		normalize.OnlyDigits(p.NuTelefone),
		p.NoLogradouro,
		p.Complemento,
		normalize.OnlyDigits(p.NuCep),
		p.Sguf,
		p.NoBairro,
		p.NomePrestador,
		p.RegistroRf,
		fmt.Sprint(p.NuAtividadeTuristica),
		p.Atividade,
		fmt.Sprint(p.NuSituacaoCadastral),
		p.Situacao,
		fmt.Sprint(p.NuUf),
		normalize.IntPtrToStr(p.LocalidadeNuUf),
		p.Municipio,
		p.Localidade,
		p.NoLocalidade,
		fmt.Sprint(p.NuLocalidade),
		fmt.Sprint(p.NuMunicipio),
		fmt.Sprint(p.NuPessoa),
		normalize.BoolToStr(p.FlPossuiVeiculo),
		fmt.Sprint(p.NuSitCadTramite),
		normalize.EmptyIfNil(p.AtividadeRedeSociais),
	}
	return w.writer.Write(row)
}

// Flush flushes the CSV writer buffer.
func (w *Writer) Flush() error {
	w.writer.Flush()
	return w.writer.Error()
}

// Close closes the underlying file.
func (w *Writer) Close() error {
	return w.file.Close()
}

