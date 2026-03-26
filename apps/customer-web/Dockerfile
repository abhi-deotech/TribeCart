# build stage
FROM node:20-alpine AS build

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy the entire monorepo workspace
COPY . .

# Install dependencies for the entire workspace
RUN pnpm install --no-frozen-lockfile

# Explicitly build @repo/ui
RUN pnpm --filter @repo/ui build

# Build the application
RUN pnpm --filter customer-web build

# final stage
FROM node:20-alpine

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy only the built customer-web application
COPY --from=build /app/apps/customer-web/.next ./apps/customer-web/.next
COPY --from=build /app/apps/customer-web/public ./apps/customer-web/public
COPY --from=build /app/apps/customer-web/package.json ./apps/customer-web/package.json


# Copy monorepo root files for pnpm workspace context
COPY --from=build /app/package.json ./package.json
COPY --from=build /app/pnpm-lock.yaml ./pnpm-lock.yaml
COPY --from=build /app/pnpm-workspace.yaml ./pnpm-workspace.yaml

# Install production dependencies for customer-web
RUN pnpm install --prod --filter=customer-web

EXPOSE 3000

CMD ["sh", "-c", "cd apps/customer-web && pnpm start"]
