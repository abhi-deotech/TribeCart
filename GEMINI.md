# TribeCart Monorepo

This document provides a detailed overview of the TribeCart monorepo structure, technologies used, and how different components interact.

## Project Overview

TribeCart is an e-commerce platform built as a monorepo, encompassing multiple frontend applications, shared UI components, API type definitions, and a suite of backend microservices. The project leverages modern web and backend technologies to provide a scalable and maintainable solution.

## Monorepo Structure

The monorepo is organized into the following top-level directories:

*   `apps/`: Contains the various frontend applications.
*   `packages/`: Houses shared libraries and configurations used across the monorepo.
*   `proto/`: Stores the Protocol Buffer (protobuf) definitions and generated Go code for gRPC services.
*   `services/`: Contains the backend microservices written in Go.

## Technologies Used

*   **Frontend:** Next.js, React, TypeScript, Tailwind CSS, Zustand, React Query, Final Form
*   **Backend:** Go, gRPC, Protocol Buffers, PostgreSQL
*   **Monorepo Management:** pnpm, Turbo
*   **Containerization:** Docker, Docker Compose

## Frontend Applications (`apps/`)

All frontend applications are built with Next.js and TypeScript.

*   **`admin-web/`**: The administrative web application for managing the e-commerce platform. It utilizes `@repo/ui` for shared UI components, `@tanstack/react-query` for data fetching, `final-form` for form management, and `zustand` for state management.
*   **`customer-web/`**: The customer-facing web application. (Note: Not included in the current `docker-compose.yml` setup, implying it's run separately or for local development.)
*   **`seller-web/`**: The seller-facing web application. (Note: Included in the current `docker-compose.yml` setup.)

## Shared Packages (`packages/`)

These packages provide reusable code and configurations across the monorepo.

*   **`api-types/`**: Contains TypeScript type definitions generated from the Protocol Buffer files, ensuring type safety between frontend and backend.
*   **`config/`**: Stores shared configurations, including ESLint presets and TypeScript base configurations, to maintain consistent coding standards.
*   **`ui/`**: A shared UI component library (e.g., `button.tsx`) used by the frontend applications. Currently, only `admin-web` explicitly lists this as a dependency, but it can be used by other frontend apps.

## Backend Services (`services/`)

The backend is composed of several microservices written in Go, communicating via gRPC and Protocol Buffers.

*   **`api-gateway/`**: The entry point for frontend applications, routing requests to the appropriate backend services. It handles CORS and depends on the `users` service.
*   **`users/`**: Manages user-related functionalities, including authentication and user profiles. It interacts with a PostgreSQL database.
*   **`products/`**: Handles product catalog management. It interacts with a PostgreSQL database.
*   **`orders/`**: Handles the creation and management of customer orders. (Note: Included in the current `docker-compose.yml` setup.)
*   **`payments/`**: Handles payment processing. (Note: Not included in the current `docker-compose.yml` setup, implying it's run separately or for local development.)

### Go Module Dependencies

All Go services depend on the shared `github.com/hepstore/tribecart/proto` module for protobuf definitions.
*   `orders`, `payments`, and `products` services use a `replace` directive in their `go.mod` files to point to the local `proto` directory, ensuring they use the monorepo's protobuf definitions during development.
*   `api-gateway` and `users` services now have this `replace` directive. For consistent local development within the monorepo, it is recommended to add `replace github.com/hepstore/tribecart/proto => ../../proto` to their `go.mod` files as well.

## Protobuf Definitions (`proto/`)

This directory contains the `.proto` files that define the gRPC service contracts and data structures for inter-service communication. It also holds the generated Go code for these definitions.

The `generate.sh` script (or `pnpm generate:proto-ts` from the root `package.json`) is responsible for generating the Go and TypeScript code from these `.proto` files.

## Docker Compose Setup

The `docker-compose.yml` file orchestrates the core services for local development:

*   **`postgres`**: A PostgreSQL database instance for data persistence.
*   **`api-gateway`**: The API Gateway service.
*   **`users-service`**: The Users microservice.
*   **`products`**: The Products microservice.
*   **`orders-service`**: The Orders microservice.
*   **`admin-web`**: The Admin web application.
*   **`seller-web`**: The Seller web application.

All services communicate within a shared `tribecart-network`. Environment variables for database connections are expected to be provided via a `.env` file.

## Build System (Turbo)

Turbo is used as the monorepo's build system, enabling efficient task execution and caching across packages and applications.

Key `turbo.json` configurations:

*   **`globalDependencies`**: Includes `**/.env.*local`, ensuring changes to local environment files invalidate the cache.
*   **`tasks`**: Defines `build`, `lint`, `test`, and `dev` scripts with appropriate dependencies and caching strategies.

## Getting Started

To set up and run the TribeCart project locally:

1.  **Clone the repository.**
2.  **Install pnpm:** If you don't have pnpm installed, follow the instructions on the [pnpm website](https://pnpm.io/installation).
3.  **Install dependencies:** From the root of the monorepo, run `pnpm install`.
4.  **Generate Protobuf types:** Run `pnpm generate:proto-ts` to generate TypeScript API types.
5.  **Set up environment variables:** Create a `.env` file in the root directory based on `.env.example` (if available) or define the necessary PostgreSQL environment variables (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`).
6.  **Start services with Docker Compose:** Run `docker compose up --build` to build and start the defined services.
7.  **Run frontend applications (optional):** For `customer-web` and `seller-web`, navigate to their respective directories and run `pnpm dev` if you wish to run them outside of Docker Compose.
