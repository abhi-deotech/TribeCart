# TribeCart: A Modern E-commerce Platform (Monorepo Edition!)

Welcome to TribeCart! This project is a great example of how modern, scalable e-commerce applications are built. If you're used to a MERN (MongoDB, Express.js, React, Node.js) stack, you'll find some familiar concepts here, but also some exciting new ones!

## What is TribeCart?

TribeCart is an e-commerce platform designed as a **monorepo**. Think of a monorepo as a single, large code repository that holds multiple, distinct projects. In our case, it contains:

*   **Multiple Frontend Applications:** Like the admin dashboard, customer-facing store, and a seller portal.
*   **Shared Libraries:** Reusable code for UI components, API definitions, and configurations.
*   **Backend Microservices:** Smaller, independent services that handle specific parts of the application (e.g., users, products, orders).

## Why a Monorepo? (And why it's different from a typical MERN setup)

In a traditional MERN project, you might have one `client` folder for your React app and one `server` folder for your Node.js/Express API. Everything is tightly coupled.

A monorepo, especially with tools like `pnpm` and `Turbo`, helps us:

*   **Manage Code Easily:** All related projects are in one place, making it easier to find and share code.
*   **Share Code Effectively:** We can create shared UI components or API types that all our frontend apps and backend services can use, ensuring consistency.
*   **Optimize Builds:** Tools like Turbo can intelligently build only the parts of the code that have changed, making development faster.

## Project Architecture: Beyond MERN

While a MERN stack typically uses a single Node.js/Express backend, TribeCart uses a **microservices architecture** for its backend. This means instead of one big server, we have several smaller, specialized servers (services) that communicate with each other.

### Frontend (Next.js)

Our frontend applications are built with **Next.js**, a React framework. Next.js is fantastic for building modern web applications because it supports:

*   **Server-Side Rendering (SSR):** Helps with performance and SEO.
*   **File-System Routing:** Makes it easy to organize your pages.
*   **API Routes:** You can even build simple APIs directly within your Next.js app if needed.

You'll find our frontend apps in the `apps/` directory:

*   `admin-web/`: The dashboard for administrators.
*   `customer-web/`: The online store for customers.
*   `seller-web/`: The portal for sellers to manage their products and orders.

### Backend (Go Microservices)

Instead of Node.js/Express, our backend services are written in **Go**. Go is a powerful, efficient language often used for building high-performance backend systems. Each service (like `users`, `products`, `orders`, `payments`) is a separate Go application.

You'll find these services in the `services/` directory.

### Database (PostgreSQL)

For data storage, we use **PostgreSQL**, a robust and popular relational database. Each microservice might interact with the database to manage its specific data (e.g., the `users` service manages user data in PostgreSQL).

## The Power of Protobufs and gRPC (Why not REST/JSON like MERN?)

This is one of the biggest differences from a typical MERN stack!

### What is Protocol Buffers (Protobuf)?

Imagine you want to send data between your frontend and backend, or between two backend services. In a MERN stack, you usually send data as **JSON** (JavaScript Object Notation). JSON is human-readable and flexible.

**Protobuf** is a language-neutral, platform-neutral, extensible mechanism for serializing structured data. Think of it as a highly efficient, binary format for data. You define your data structures once in a `.proto` file, and then you can generate code in many different programming languages (like Go, JavaScript, Python, Java) that can easily read and write that data.

**Why use Protobuf?**

*   **Efficiency:** Protobuf messages are much smaller than JSON, leading to faster data transfer.
*   **Speed:** Encoding and decoding Protobuf messages is faster than JSON.
*   **Strong Typing:** Because you define your data structure in a `.proto` file, the generated code provides strong type checking, which helps prevent errors and makes your code more reliable. This is a big win for large, complex systems.
*   **Language Agnostic:** Since you can generate code for many languages, it's perfect for a microservices architecture where different services might be written in different languages.

You'll find our `.proto` files in the `proto/` directory.

### What is gRPC?

Just as REST is a common architectural style for APIs that often uses JSON over HTTP, **gRPC** is a modern, high-performance Remote Procedure Call (RPC) framework that uses Protobuf as its interface definition language and HTTP/2 for transport.

**Why use gRPC?**

*   **Performance:** Built on HTTP/2, gRPC offers features like multiplexing (sending multiple requests over a single connection) and header compression, making it very fast.
*   **Strongly Typed APIs:** Because it uses Protobuf, gRPC APIs are strongly typed, meaning you know exactly what data to send and receive, reducing bugs.
*   **Streaming:** gRPC supports different types of streaming (client-side, server-side, and bi-directional), which is great for real-time applications.

**Comparison to REST/JSON (MERN vs. TribeCart)**

| Feature           | Typical MERN (REST/JSON)                               | TribeCart (gRPC/Protobuf)                                  |
| :---------------- | :----------------------------------------------------- | :--------------------------------------------------------- |
| **Architecture**  | Monolithic (one big backend)                           | Microservices (many small, independent backends)           |
| **Backend Lang.** | Node.js (JavaScript)                                   | Go                                                         |
| **API Protocol**  | REST over HTTP/1.1                                     | gRPC over HTTP/2                                           |
| **Data Format**   | JSON (human-readable, text-based)                      | Protobuf (binary, highly efficient)                        |
| **API Definition**| Often less formal, relies on documentation/examples    | `.proto` files (formal, strongly typed, code-generatable)  |
| **Performance**   | Good for many use cases, but can be less efficient     | Very high performance, lower latency, efficient bandwidth  |
| **Type Safety**   | Less inherent type safety (runtime checks often needed)| Strong type safety from generated code                     |
| **Use Case**      | Web apps, mobile apps, general-service APIs            | Microservices, high-performance APIs, inter-service comm.  |

## Key Components Explained

### Frontend Applications (`apps/`)

*   `admin-web/`: The administrative interface. Uses `@repo/ui` for shared components, `@tanstack/react-query` for data fetching, `final-form` for forms, and `zustand` for state management.
*   `customer-web/`: The public-facing e-commerce store.
*   `seller-web/`: The portal for sellers to manage their listings.

### Shared Packages (`packages/`)

*   `api-types/`: Contains TypeScript types generated from our `.proto` files. This ensures your frontend knows exactly what data to expect from the backend.
*   `config/`: Holds shared configurations like ESLint rules and TypeScript settings, ensuring consistent code style across the monorepo.
*   `ui/`: Our shared UI component library. Components defined here (like buttons) can be reused across all frontend applications.

### Backend Services (`services/`)

*   `api-gateway/`: This is the "front door" for our frontend applications. It receives requests from the web apps and forwards them to the correct microservice. It also handles things like CORS (Cross-Origin Resource Sharing).
*   `users/`: Manages all user-related data and operations (e.g., user registration, login, profiles). It connects to the PostgreSQL database.
*   `products/`: Manages all product-related data and operations (e.g., listing products, product details). It also connects to the PostgreSQL database.
*   `orders/`: Handles the creation and management of customer orders.
*   `payments/`: Manages payment processing.

## Getting Started: Your First Steps!

Follow these steps to get TribeCart up and running on your local machine.

### Prerequisites

Before you start, make sure you have these installed:

*   **Git:** For cloning the repository.
*   **Node.js & npm/yarn:** (If you plan to run frontend apps outside Docker)
*   **pnpm:** Our preferred package manager for the monorepo. If you don't have it, install it globally:
    ```bash
    npm install -g pnpm
    ```
*   **Docker & Docker Compose:** Essential for running our backend services and database.

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-repo/TribeCart.git # Replace with actual repo URL
    cd TribeCart
    ```
2.  **Install monorepo dependencies:**
    ```bash
    pnpm install
    ```

### Generating API Types

Our frontend applications need to understand the data structures defined in our `.proto` files. We generate TypeScript types from these:

```bash
pnpm generate:proto-ts
```

### Environment Variables

The `docker-compose.yml` file expects some environment variables, especially for the PostgreSQL database. Create a file named `.env` in the root of the `TribeCart` directory and add the following (you can choose your own values):

```
POSTGRES_USER=tribecart_user
POSTGRES_PASSWORD=tribecart_password
POSTGRES_DB=tribecart_db
```

### Running Services with Docker Compose

This will build the Docker images for our backend services and the `admin-web` frontend, and then start them along with the PostgreSQL database.

```bash
docker compose up --build
```

*   **`--build`**: This flag ensures that Docker images are rebuilt from their `Dockerfile`s. It's good practice to use this when you make changes to service code or Dockerfiles.
*   You should see logs from `postgres`, `users-service`, `products`, `api-gateway`, and `admin-web`.

### Running Frontend Applications Locally (Optional)

While `admin-web` is included in `docker-compose.yml`, you might want to run `customer-web` or `seller-web` locally for development, or run `admin-web` outside of Docker.

1.  **Open a new terminal.**
2.  **Navigate to the app directory:**
    ```bash
    cd apps/customer-web # or apps/seller-web, or apps/admin-web
    ```
3.  **Start the development server:**
    ```bash
    pnpm dev
    ```
    *   This will usually start the app on `http://localhost:3000` (or another port if 3000 is taken).

## Contributing

We welcome contributions! Please read our `CONTRIBUTING.md` (if available) for guidelines on how to contribute.

---

Happy coding!
