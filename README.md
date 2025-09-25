# Multithread CEP Consult

Este projeto realiza consultas de CEP utilizando múltiplas APIs de forma concorrente em Go.

## Pré-requisitos
- Go 1.21 ou superior

## Como executar

1. Clone o repositório:
   ```sh
   git clone <url-do-repositorio>
   cd multithread
   ```

2. Execute o programa principal:
   ```sh
   go run cmd/main.go
   ```

3. O resultado será exibido no terminal, mostrando a resposta das APIs consultadas.

## Estrutura do projeto
- `cmd/main.go`: Ponto de entrada do programa.
- `internal/types/`: Tipos e interfaces utilizadas.
- `internal/workers/`: Implementações das APIs consultadas.

## Personalização
- Para consultar outro CEP, altere a variável `cep` no arquivo `cmd/main.go`.
- Para adicionar novas APIs, implemente a interface `ConsultCep` em `internal/types/types.go`.

## Licença
MIT
