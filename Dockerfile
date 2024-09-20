# Usar uma imagem base do Go
FROM golang:1.22 AS builder

# Definir o diretório de trabalho
WORKDIR /app

# Copiar os arquivos go.mod e go.sum e instalar as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o restante do código fonte
COPY . .

# Compilar os dois executáveis e verificar a existência
RUN go build -o app cmd/app/main.go
RUN go build -o api cmd/api/api.go

# Etapa final
FROM alpine:latest

WORKDIR /app

# Copiar os binários do estágio de construção
COPY --from=builder /app/app .
COPY --from=builder /app/api .


# Garantir que os executáveis tenham permissões de execução
RUN chmod +x app api
RUN apk add --no-cache libc6-compat

# Expor as portas, se necessário (ajuste conforme sua necessidade)
EXPOSE 8081
EXPOSE 8080

# Para o docker-compose, não vamos definir um CMD aqui
