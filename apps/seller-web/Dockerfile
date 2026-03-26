FROM node:20-alpine AS base

# Install pnpm
RUN npm install -g pnpm

# Install dependencies in the workspace
FROM base AS deps
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm fetch

# Build the UI package
FROM deps AS ui-build
WORKDIR /app
COPY packages/ui packages/ui
COPY packages/config packages/config
RUN pnpm --filter @repo/ui add -D @types/node
RUN pnpm --filter @repo/ui build

# Build the Next.js application
FROM deps AS builder
WORKDIR /app
COPY --from=ui-build /app/packages/ui packages/ui
COPY packages/tsconfig packages/tsconfig
COPY packages/config packages/config
COPY apps/seller-web apps/seller-web
COPY turbo.json turbo.json
COPY pnpm-workspace.yaml pnpm-workspace.yaml
RUN pnpm add lucide-react --filter seller-web
RUN pnpm --filter seller-web build

# Production image
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV production

# Copy built Next.js app
COPY --from=builder /app/apps/seller-web/.next apps/seller-web/.next
COPY --from=builder /app/apps/seller-web/public apps/seller-web/public
COPY --from=builder /app/apps/seller-web/package.json apps/seller-web/package.json

# Copy built UI package (dist folder)
COPY --from=ui-build /app/packages/ui/dist packages/ui/dist
COPY --from=ui-build /app/packages/ui/package.json packages/ui/package.json
COPY --from=ui-build /app/packages/ui/index.tsx packages/ui/index.tsx
COPY --from=ui-build /app/packages/ui/src packages/ui/src

# Set up pnpm for production
RUN npm install -g pnpm
RUN pnpm install --prod

EXPOSE 3000

CMD ["pnpm", "start", "--filter", "seller-web"]
