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
RUN pnpm --filter admin-web build

# final stage
FROM node:20-alpine

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy only the built admin-web application
COPY --from=build /app/apps/admin-web/.next ./apps/admin-web/.next
COPY --from=build /app/apps/admin-web/public ./apps/admin-web/public
COPY --from=build /app/apps/admin-web/package.json ./apps/admin-web/package.json


# Copy monorepo root files for pnpm workspace context
COPY --from=build /app/package.json ./package.json
COPY --from=build /app/pnpm-lock.yaml ./pnpm-lock.yaml
COPY --from=build /app/pnpm-workspace.yaml ./pnpm-workspace.yaml

# Install production dependencies for admin-web
RUN pnpm install --prod --filter=admin-web

EXPOSE 3000

CMD ["sh", "-c", "cd apps/admin-web && pnpm start"]
