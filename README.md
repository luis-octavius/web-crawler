# Web Crawler em Go 🕷️

Um web crawler de alta performance desenvolvido em Go, com suporte a crawling concorrente, rate limiting, retry automático e rotação de proxies.

## Características

### Funcionalidades Principais
- **Crawling Concorrente**: Utiliza goroutines para crawling paralelo eficiente
- **Rate Limiting**: Controle de taxa de requisições para evitar bloqueios
- **Retry Automático**: Backoff exponencial para lidar com falhas temporárias
- **Rotação de Proxies**: Suporte para múltiplos proxies com health check
- **Extração de Dados**: Coleta H1, primeiro parágrafo, links e imagens
- **Relatório CSV**: Exportação dos dados coletados em formato CSV
- **Headers Realistas**: User-Agent e headers que imitam navegadores reais

### Tecnologias Utilizadas
- **Go 1.25+**: Linguagem principal
- **goquery**: Parsing e extração de dados HTML (similar ao jQuery)
- **golang.org/x/time/rate**: Rate limiting robusto
- **Concorrência nativa**: Goroutines e channels para alta performance

## Pré-requisitos

- Go 1.25 ou superior
- Conexão com a internet

## Instalação

```bash
# Clone o repositório
git clone https://github.com/luis-octavius/web-crawler.git
cd web-crawler

# Instale as dependências
go mod download

# Compile o projeto
go build -o crawler
```

## Uso

### Uso Básico

```bash
# Sintaxe
./crawler <URL_BASE> <MAX_PÁGINAS> <MAX_CONCORRÊNCIA>

# Exemplo
./crawler https://example.com 100 5
```

### Parâmetros

- **URL_BASE**: URL inicial para começar o crawling
- **MAX_PÁGINAS**: Número máximo de páginas a serem coletadas
- **MAX_CONCORRÊNCIA**: Número máximo de requisições simultâneas

### Exemplos de Uso

```bash
# Crawl de até 50 páginas com 3 workers concorrentes
./crawler https://blog.example.com 50 3

# Crawl mais agressivo: 200 páginas com 10 workers
./crawler https://site.com 200 10

# Crawl conservador: 20 páginas com 2 workers
./crawler https://api-docs.example.com 20 2
```

## Arquitetura

### Estrutura de Arquivos

```
web-crawler/
├── main.go              # Ponto de entrada da aplicação
├── config.go            # Configuração e gerenciamento de estado
├── crawl_page.go        # Lógica principal de crawling
├── get_html.go          # Requisições HTTP com retry
├── rate_limiter.go      # Controle de taxa de requisições
├── retry.go             # Lógica de retry com backoff exponencial
├── proxy.go             # Rotação de proxies
├── parse_html.go        # Parsing e extração de dados
├── page_data.go         # Estruturas de dados
├── normalize_url.go     # Normalização de URLs
├── csv_report.go        # Geração de relatórios CSV
└── *_test.go           # Testes unitários
```

### Fluxo de Execução

```
1. Inicialização
   ↓
2. Configuração (rate limiter, proxies, concorrência)
   ↓
3. Crawling Recursivo
   ├── Rate Limiting
   ├── Delay entre requisições
   ├── Retry automático (se falhar)
   ├── Extração de dados
   └── Crawl de links encontrados
   ↓
4. Geração de Relatório CSV
```

## Funcionalidades Avançadas

### Rate Limiting

O crawler implementa rate limiting para evitar sobrecarga nos servidores:

```go
// Configuração padrão: 2 requisições/segundo, burst de 5
rateLimiter := NewRateLimiter(2.0, 5)
```

### Retry Automático

Implementa backoff exponencial para lidar com falhas temporárias:

```go
// Configuração padrão
RetryConfig{
    MaxRetries:     3,
    InitialBackoff: 1 * time.Second,
    MaxBackoff:     30 * time.Second,
    Multiplier:     2.0,
}
```

### Rotação de Proxies

Suporte para múltiplos proxies com health check automático:

```go
proxyURLs := []string{
    "http://proxy1.example.com:8080",
    "http://proxy2.example.com:8080",
}
proxyRotator, _ := NewProxyRotator(proxyURLs, true)
```

### Headers Realistas

O crawler usa headers que imitam navegadores reais para evitar bloqueios:

```
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36...
Accept: text/html,application/xhtml+xml,application/xml;q=0.9...
Accept-Language: pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7
```

## Dados Coletados

Para cada página, o crawler extrai:

- **URL**: URL normalizada da página
- **H1**: Primeiro heading H1 encontrado
- **Primeiro Parágrafo**: Primeiro parágrafo de texto
- **Links de Saída**: Todos os links encontrados na página
- **Imagens**: URLs de todas as imagens

## Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com cobertura
go test -cover ./...

# Executar testes de um arquivo específico
go test -v normalize_url_test.go normalize_url.go
```

## Performance

### Benchmarks

- **Concorrência**: Suporta centenas de goroutines simultâneas
- **Rate Limiting**: Controle preciso de requisições/segundo
- **Memória**: Uso eficiente com estruturas de dados otimizadas

### Otimizações Implementadas

1. **Controle de Concorrência**: Channel-based semaphore
2. **Mutex para Estado Compartilhado**: Evita race conditions
3. **Rate Limiting**: Previne sobrecarga e bloqueios
4. **Retry Inteligente**: Backoff exponencial para falhas temporárias
5. **Delay entre Requisições**: Comportamento mais "humano"

## Tratamento de Erros

O crawler lida com diversos cenários de erro:

- **429 (Rate Limit)**: Retry automático com backoff
- **5xx (Erro de Servidor)**: Retry automático
- **4xx (Erro de Cliente)**: Não retenta (exceto 429)
- **Timeout**: Retry com backoff exponencial
- **Proxy Inválido**: Remove proxy da rotação

## Boas Práticas Implementadas

1. **Respeito ao robots.txt**: (pode ser implementado)
2. **Rate Limiting**: Evita sobrecarga nos servidores
3. **User-Agent Identificável**: Headers realistas
4. **Delay entre Requisições**: Comportamento não-agressivo
5. **Tratamento de Erros Robusto**: Retry inteligente

## 🚧 Melhorias Futuras

- [ ] Suporte a robots.txt
- [ ] Suporte a sitemap.xml
- [ ] Detecção de paginação dinâmica
- [ ] Suporte a JavaScript rendering (Puppeteer/Playwright)
- [ ] Integração com Redis para fila de URLs
- [ ] Suporte a autenticação (cookies, tokens)
- [ ] Detecção de bloqueios anti-bot
- [ ] Métricas e monitoramento (Prometheus)
- [ ] Suporte a bancos de dados (PostgreSQL/MongoDB)
- [ ] API REST para controle do crawler

## Exemplos de Saída

### Console

```
starting crawl of: https://example.com...
✓ Crawled: https://example.com (encontrados 15 links)
✓ Crawled: https://example.com/about (encontrados 8 links)
✓ Crawled: https://example.com/contact (encontrados 5 links)
...
Execution took 5.234 seconds
```

### Relatório CSV

```csv
URL,H1,FirstParagraph,OutgoingLinks,ImageURLs
https://example.com,Welcome,This is an example...,15,3
https://example.com/about,About Us,We are a company...,8,2
```

## Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para:

1. Fazer fork do projeto
2. Criar uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abrir um Pull Request

## Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## Autor

**Luis Octavius**

- GitHub: [@luis-octavius](https://github.com/luis-octavius)

## Agradecimentos

- [goquery](https://github.com/PuerkitoBio/goquery) - Excelente biblioteca para parsing HTML
- [Colly](https://github.com/gocolly/colly) - Inspiração para algumas funcionalidades
- Comunidade Go - Pela linguagem incrível e ecossistema rico

---

Se este projeto foi útil para você, considere dar uma estrela!
