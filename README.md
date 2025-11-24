# cadastur-csv

Ferramenta de linha de comando em Go para extrair prestadores de serviço turístico do Cadastur, paginar automaticamente os resultados, normalizar campos e gerar um CSV pronto para análise.

Principais pontos:

- Recupera a lista de UFs e atividades do Cadastur

- Faz requisições paginadas para obter todos os prestadores

- Normaliza telefone e CEP (apenas dígitos)

- Corrige problemas comuns de encoding e grava CSV com BOM para compatibilidade com Excel no Windows

---

## Requisitos

- Go 1.25 ou superior (o `go.mod` indica 1.25.4)

- Conexão com a internet para acessar a API pública do Cadastur

Nota: o projeto pode baixar automaticamente dependências (ex.: `golang.org/x/net/html/charset`) ao executar `go build`/`go run`.

---

## Instruções rápidas (Windows / PowerShell)

1) Clonar o repositório

```powershell
git clone https://github.com/rsornellas/cadastur-csv
cd 'cadastur-to-csv'
```

2) Build (opcional)

```powershell
go build -v -o cadastur-csv ./cmd/cadastur-csv
.\cadastur-csv
```

3) Ou executar sem build

```powershell
go run ./cmd/cadastur-csv
```

Ao executar, o CLI fará perguntas interativas em Português:

- ID da UF (padrão 24 — Santa Catarina)

- ID da atividade (padrão 29 — Guia de Turismo)

- Cidade (opcional)

---

## Sobre o CSV gerado

- Nome: `prestadores-atividade-<ID>-<slug>.csv`
- O CSV contém todas as colunas retornadas pela API (id, tipoPessoa, numeroCadastro, inicioVigencia, fimVigencia, website, telefone, logradouro, complemento, cep, uf, bairro, nomePrestador, registroRf, nuAtividadeTuristica, atividade, nuSituacaoCadastral, situacao, nuUf, localidadeNuUf, municipio, localidade, noLocalidade, nuLocalidade, nuMunicipio, nuPessoa, possuiVeiculo, nuSitCadTramite, atividadeRedeSociais).
- O arquivo é gravado com BOM UTF-8 (EF BB BF) — isso ajuda o Excel no Windows a detectar corretamente UTF-8 e evitar exibição de caracteres corrompidos (ex.: "Ã¡").

Se for necessário gerar CSV em encoding Windows-1252 (ANSI) por compatibilidade com sistemas antigos, pode-se adaptar o writer para codificar o arquivo antes de gravar.

---

## Estrutura do projeto

```
cadastur-csv/
├── cmd/cadastur-csv/main.go        # entrypoint do CLI
├── internal/
│   ├── cadastur/                   # cliente HTTP, endpoints e service
│   ├── cli/                         # prompts e orquestração (Run)
│   ├── csvx/                        # writer CSV
│   └── normalize/                   # utilitários de normalização
├── README.md
├── LICENSE
└── go.mod
```

---

## Observações e notas técnicas

- O cliente HTTP agora respeita o charset declarado pelo servidor (usa `golang.org/x/net/html/charset`) e converte para UTF-8 quando necessário.
- Alguns campos na API podem retornar tipos inconsistentes (ex.: boolean em vez de string). O modelo foi ajustado para tolerar essas variações.
- Há heurística para corrigir mojibake já presente nos dados (caso raro) e a escrita com BOM ajuda consumidores como Excel.

---

## Quer adicionar automação? Sugestões

- Adicionar flags para execução não interativa (`--uf`, `--activity`, `--city`).
- Gerar também um `.xlsx` usando uma biblioteca como `excelize` para evitar problemas de importação no Excel.
- Testes unitários para `normalize` (OnlyDigits, MsToDate, FixMojibake).

---

## Licença

MIT — veja o arquivo `LICENSE`.

---

## Créditos

- Desenvolvido por Raphael Ornellas
- Auxílio técnico: ChatGPT
- Dados: API pública Cadastur (Ministério do Turismo)
