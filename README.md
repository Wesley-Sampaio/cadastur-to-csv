# cadastur-csv

Ferramenta de linha de comando escrita em Go para buscar dados de prestadores de serviços turísticos no **Cadastur**, processar todas as páginas disponíveis da API pública e gerar um arquivo **CSV completo**, limpo e padronizado.

---

# 🇧🇷 Versão em Português

## 🎯 Objetivo do Projeto
O **cadastur-csv** foi criado para resolver uma limitação prática do site do Cadastur, que:

- não permite exportar dados para planilha,
- exibe somente listas resumidas,
- exige que cada item seja aberto manualmente para acessar detalhes.

Isso tornava extremamente lento o processo de coleta de informações — especialmente para equipes que dependem de grande volume de dados para operar.

O projeto nasceu de uma demanda do time comercial da **Natural Extremo**, que precisava:

- acessar rapidamente informações de **Guias de Turismo**,  
- analisar regiões, contatos e perfis dos profissionais,  
- montar parcerias estratégicas e ampliar canais de vendas,  
- agilizar ações sem depender de processos manuais no portal do Cadastur.

Com este CLI, todo o fluxo passou a ser **automático, rápido e confiável**, entregando planilhas completas com dados reais diretamente da API pública.

---

## ✅ Benefícios
- **Economia de tempo:** elimina a necessidade de abrir registros manualmente.
- **Padronização:** telefones e CEPs são normalizados automaticamente.
- **Automação completa:** busca todas as páginas disponíveis sem intervenção.
- **Flexibilidade:** permite filtrar por UF, atividade turística e cidade (opcional).
- **Transparência:** gera um CSV com todas as colunas retornadas pela API.
- **Zero dependências externas:** apenas Go puro (standard library).

---

## ✨ Funcionalidades
- Consulta lista de **UFs** diretamente do Cadastur.
- Consulta lista de **atividades turísticas** (ex.: Guia de Turismo).
- Interação via terminal:
  - escolha de UF,
  - escolha de atividade,
  - entrada opcional de cidade.
- Paginação automática até o último resultado.
- Normalização de campos:
  - telefone → somente números,
  - CEP → somente números.
- Geração de arquivo CSV com nome baseado na atividade escolhida.

---

## 🚀 Como Executar

### 1. Instale o Go (se necessário)
https://go.dev/dl/

### 2. Clone o repositório
```bash
git clone https://github.com/rsornellas/cadastur-csv
cd cadastur-csv
```

### 3. Execute o programa
```bash
go run ./cmd/cadastur-csv
```

O CLI fará perguntas em português:

- selecione a UF,
- selecione a atividade,
- informe (ou deixe em branco) a cidade.

---

## 📄 Estrutura do CSV de Saída
O CSV contém todas as colunas retornadas pela API, incluindo:

```
id,tipoPessoa,numeroCadastro,inicioVigencia,fimVigencia,website,telefone,logradouro,complemento,cep,uf,bairro,nomePrestador,registroRf,nuAtividadeTuristica,atividade,nuSituacaoCadastral,situacao,nuUf,localidadeNuUf,municipio,localidade,noLocalidade,nuLocalidade,nuMunicipio,nuPessoa,possuiVeiculo,nuSitCadTramite,atividadeRedeSociais
```

---

## 🏗️ Estrutura do Projeto
```
cadastur-csv/
├── cmd/
│   └── cadastur-csv/
│       └── main.go
├── internal/
│   ├── cadastur/
│   │   ├── client.go
│   │   ├── endpoints.go
│   │   ├── models.go
│   │   └── service.go
│   ├── cli/
│   │   ├── prompts.go
│   │   └── run.go
│   ├── csvx/
│   │   └── writer.go
│   └── normalize/
│       └── normalize.go
├── README.md
├── LICENSE
├── .gitignore
└── go.mod
```

---

## 📦 Gerar Binário
```bash
go build -o cadastur-csv ./cmd/cadastur-csv
./cadastur-csv
```

---

## 📜 Licença
MIT - Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

## 🙌 Créditos
- Desenvolvido por **Raphael Ornellas**
- Auxílio técnico com ChatGPT
- Dados públicos fornecidos pelo Cadastur (Ministério do Turismo)

---

## 🔍 Palavras‑chave (SEO)
cadastur csv, exportar cadastur, extrair dados cadastur, api cadastur, guia de turismo cadastur, como baixar dados do cadastur, prestadores de serviços turísticos, turismo brasil dados, automação cadastur, cadastur scraper, natural extremo, ferramentas para turismo, análise de dados turismo, lista de guias de turismo brasil, cadastur download, cadastur api csv, exportar guias de turismo.

---

# 🇺🇸 English Version

## 🎯 Project Purpose
**cadastur-csv** was created to solve a practical limitation of the Cadastur website, which:

- does not allow exporting data to spreadsheets,
- only displays summarized lists,
- requires opening each record individually to view details.

For teams that need large-scale data analysis, this process is slow and inefficient.

The commercial team at **Natural Extremo** needed a fast way to access:

- guide information (“Guias de Turismo”),
- contact and region details,
- data-driven insights for partnerships and sales expansion.

This tool fully automates that workflow, retrieving all data directly from the public API and exporting it into a clean CSV file.

---

## ✅ Benefits
- **Time-saving:** no more manual data collection.
- **Standardization:** phone numbers and ZIP codes normalized automatically.
- **Full automation:** fetches all pages until no more results.
- **Flexible filters:** choose UF, activity, and optional city.
- **Complete data export:** all fields returned by the API.
- **Zero dependencies:** pure Go (standard library).

---

## ✨ Features
- Fetches Cadastur **UF list**
- Fetches **tourism activity list**
- Terminal prompts for UF, activity, and optional city
- Automatic pagination
- Normalizes:
  - phone numbers → digits only
  - ZIP code (CEP) → digits only
- Saves CSV with descriptive file name:
  - `prestadores-atividade-<ID>-<slug>.csv`

---

## 🚀 How to Run
1. Install Go: https://go.dev/dl/  
2. Clone the repo:
```bash
git clone https://github.com/rsornellas/cadastur-csv
cd cadastur-csv
```
3. Run:
```bash
go run ./cmd/cadastur-csv
```

---

## 📜 License
MIT (recommended)

---

## 🙌 Credits
- Developed by **Raphael Ornellas**
- Guidance with ChatGPT
- Public data by Cadastur (Brazilian Ministry of Tourism)