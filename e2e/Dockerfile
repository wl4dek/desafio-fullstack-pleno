FROM mcr.microsoft.com/playwright:v1.60.0-noble

WORKDIR /app

# Habilita pnpm
RUN corepack enable

# Copia arquivos de dependências
COPY package.json package-lock.json pnpm-lock.yaml* ./

# Instala dependências
RUN pnpm install --frozen-lockfile

# Copia o restante do projeto
COPY . .

# Instala browsers do Playwright (normalmente já vêm na imagem)
RUN pnpm exec playwright install

CMD ["pnpm", "test"]