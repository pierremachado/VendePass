# Dockerfile para a aplicação React com Vite
FROM node:18

# Definir o diretório de trabalho
WORKDIR /app

# Copiar os arquivos de configuração e instalar as dependências
COPY package.json package-lock.json ./
RUN npm install

# Copiar o restante do código fonte
COPY . .

# Executar o servidor de desenvolvimento do Vite
CMD ["npm", "run", "dev"]

# Expor a porta padrão do Vite
EXPOSE 5173
