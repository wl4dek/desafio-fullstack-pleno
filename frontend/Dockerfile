FROM node:lts-alpine AS deps
WORKDIR /app

RUN apk add --no-cache libc6-compat
RUN npm i -g pnpm

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

FROM node:lts-alpine AS builder
WORKDIR /app

RUN npm i -g pnpm

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG API_URL
ENV API_URL=$API_URL
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

RUN pnpm run build

FROM node:lts-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000

RUN addgroup -S nextjs && adduser -S nextjs -G nextjs

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/node_modules ./node_modules

USER nextjs

EXPOSE 3000

CMD ["node_modules/.bin/next", "start"]
