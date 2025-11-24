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

	// Write UTF-8 BOM so Excel on Windows detects UTF-8 encoding when opening the CSV.
	// This helps avoid mojibake like "Ã¡" when users open the CSV by double-clicking in Explorer.
	if _, err := f.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		f.Close()
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
	// Wrap numeric-like IDs in a formula string so Excel opens them as text
	// and doesn't convert to scientific notation or truncate digits.
	fmt.Sprintf("=\"%s\"", p.NumeroCadastro),
		normalize.MsToDate(p.DtInicioVigencia),
		normalize.MsToDate(p.DtFimVigencia),
		normalize.EmptyIfNil(p.NoWebSite),
		normalize.OnlyDigits(p.NuTelefone),
		normalize.FixMojibake(p.NoLogradouro),
		normalize.FixMojibake(p.Complemento),
		normalize.OnlyDigits(p.NuCep),
		p.Sguf,
		normalize.FixMojibake(p.NoBairro),
		normalize.FixMojibake(p.NomePrestador),
		p.RegistroRf,
		fmt.Sprint(p.NuAtividadeTuristica),
		normalize.FixMojibake(p.Atividade),
		fmt.Sprint(p.NuSituacaoCadastral),
		normalize.FixMojibake(p.Situacao),
		fmt.Sprint(p.NuUf),
		normalize.IntPtrToStr(p.LocalidadeNuUf),
		normalize.FixMojibake(p.Municipio),
		normalize.FixMojibake(p.Localidade),
		normalize.FixMojibake(p.NoLocalidade),
		fmt.Sprint(p.NuLocalidade),
		fmt.Sprint(p.NuMunicipio),
		fmt.Sprint(p.NuPessoa),
		normalize.BoolToStr(p.FlPossuiVeiculo),
		fmt.Sprint(p.NuSitCadTramite),
		normalize.FixMojibake(normalize.EmptyIfNil(p.AtividadeRedeSociais)),
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
