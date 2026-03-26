# TribeCart — Product Requirements Document (PRD)

> **Version:** 1.0.0  
> **Date:** 2026-03-27  
> **Status:** Living Document — Update with every major decision  
> **Authors:** TribeCart Engineering & Product Team  
> **Classification:** Internal — Blueprint for all development

---

## Table of Contents

1. [Vision & Mission](#1-vision--mission)
2. [Working Backwards — The Customer Letter](#2-working-backwards--the-customer-letter)
3. [Target Users & Personas](#3-target-users--personas)
4. [Market Opportunity](#4-market-opportunity)
5. [Product Philosophy & Design Principles](#5-product-philosophy--design-principles)
6. [Platform Architecture Overview](#6-platform-architecture-overview)
7. [Feature Specifications by Domain](#7-feature-specifications-by-domain)
   - 7.1 [User & Identity Service](#71-user--identity-service)
   - 7.2 [Product Catalog Service](#72-product-catalog-service)
   - 7.3 [Search & Discovery](#73-search--discovery)
   - 7.4 [Shopping Cart & Checkout](#74-shopping-cart--checkout)
   - 7.5 [Orders Service](#75-orders-service)
   - 7.6 [Payments Service](#76-payments-service)
   - 7.7 [Inventory & Fulfillment](#77-inventory--fulfillment)
   - 7.8 [Seller Portal](#78-seller-portal)
   - 7.9 [Admin Portal](#79-admin-portal)
   - 7.10 [Customer Portal](#710-customer-portal)
   - 7.11 [Reviews & Ratings](#711-reviews--ratings)
   - 7.12 [Notifications Service](#712-notifications-service)
   - 7.13 [Analytics & Reporting](#713-analytics--reporting)
   - 7.14 [Promotions & Discounts](#714-promotions--discounts)
8. [Non-Functional Requirements](#8-non-functional-requirements)
9. [API Gateway & Communication Contracts](#9-api-gateway--communication-contracts)
10. [Database Strategy](#10-database-strategy)
11. [Security & Compliance](#11-security--compliance)
12. [DevOps & Infrastructure](#12-devops--infrastructure)
13. [Phased Roadmap](#13-phased-roadmap)
14. [Success Metrics & KPIs](#14-success-metrics--kpis)
15. [Risks & Mitigations](#15-risks--mitigations)
16. [Glossary](#16-glossary)

---

## 1. Vision & Mission

### Vision
To build the most **trustworthy, efficient, and empowering** e-commerce marketplace in the world — one where any seller can reach any customer, and every customer gets exactly what they need, delivered on time, every time.

### Mission
TribeCart empowers sellers of all sizes — from solo artisans to enterprise brands — to build thriving businesses online, while giving customers a seamless, joyful, and safe shopping experience.

### Core Differentiators
| Differentiator | Description |
|---|---|
| **Tribe Community** | Social layers — buyers and sellers form communities ("Tribes") around shared interests |
| **Transparent Pricing** | No hidden fees; sellers see exactly what they earn before listing |
| **Built-in Analytics** | Every seller gets Amazon-level analytics without a third-party subscription |
| **Speed-first** | Sub-200ms API responses; next-day delivery network from day one |
| **Seller-First Policy** | Disputes favor the data, not the louder party |

---

## 2. Working Backwards — The Customer Letter

*(Amazon's "Working Backwards" model: Write the press release before the code.)*

---

**FOR IMMEDIATE RELEASE**

**TribeCart Launches: The Marketplace That Finally Works for Everyone**

*Sellers keep more. Buyers find better. Everyone wins.*

[City, Date] — TribeCart today announced the launch of the TribeCart Marketplace, a next-generation e-commerce platform designed from the ground up to serve both sellers and buyers with radical transparency, speed, and community.

"I've sold on three other platforms," said a beta seller. "On TribeCart, I set up my shop in 20 minutes, launched 50 products using the bulk importer, and I can see *exactly* which customers bought what and why. No other platform gives me this."

For customers, TribeCart delivers a curated, community-driven shopping experience. Product pages are rich with verified reviews, video demonstrations, and seller-backed guarantees. A purchase on TribeCart is a purchase backed by data.

"I found a handmade leather bag from a seller three cities away," said a beta customer. "I ordered Tuesday, it arrived Thursday. The seller even sent a video of it being made. That's not a purchase — that's a story."

TribeCart is available today at tribecart.com.

---

## 3. Target Users & Personas

### 3.1 Customer — "Priya, 28"
- **Background**: Urban professional, shops online 3–5x/month
- **Goals**: Find quality products fast, trust the seller, get clear delivery dates
- **Pain Points**: Fake reviews, hidden shipping fees, confusing return policies
- **TribeCart Promises**: Verified reviews, transparent pricing, one-click returns

### 3.2 Seller — "Ravi, 35"
- **Background**: Mid-size retail business, moving from offline to online
- **Goals**: Grow sales, understand customers, minimize platform friction
- **Pain Points**: High commission rates, opaque algorithms, poor customer data access
- **TribeCart Promises**: Lowest fees in class, full analytics access, dedicated onboarding

### 3.3 Power Seller — "TechGiant Corp"
- **Background**: Large brand with 10K+ SKUs, needs API access and bulk management
- **Goals**: Sync ERP, manage inventory programmatically, get SLA guarantees
- **Pain Points**: Rate limits, poor API docs, manual fulfillment
- **TribeCart Promises**: Documented gRPC/REST API, bulk operations, dedicated account manager

### 3.4 Admin — "Ananya, Internal"
- **Background**: Operations lead managing platform compliance and disputes
- **Goals**: Fast resolution, clear audit trails, fraud detection
- **TribeCart Promises**: Rich admin portal, real-time alerts, full audit log

---

## 4. Market Opportunity

| Market | Size (2025 Est.) |
|---|---|
| India E-commerce GMV | $150 Billion |
| Southeast Asia E-commerce | $200 Billion |
| Global Marketplace Software | $18 Billion |
| Target Addressable Market (TAM) | $25 Billion |
| Serviceable Market (SAM, Year 1) | $2 Billion |

### Competitive Landscape
| Platform | Strength | Weakness TribeCart Addresses |
|---|---|---|
| Amazon | Scale, logistics | Seller unfriendly, high fees |
| Flipkart | India reach | Opaque seller tools |
| Shopify | Seller empowerment | No built-in marketplace/buyers |
| Meesho | Social commerce | Low quality control |
| **TribeCart** | Community + data + transparency | — |

---

## 5. Product Philosophy & Design Principles

1. **Customer Obsession** — Every feature starts with "what does this do for the buyer or seller?"
2. **Data is a First-Class Citizen** — Every action creates a data point; every data point creates an insight
3. **Radical Transparency** — Fees, algorithms, and policies are always visible
4. **Speed As a Feature** — Slowness is a bug. Every page load under 2 seconds. Every API under 200ms P99
5. **Design for Scale from Day 1** — Architecture must handle 1M concurrent users without a rewrite
6. **Security by Default** — Zero plaintext secrets, zero unencrypted PII in transit or at rest
7. **Fail Loudly, Recover Gracefully** — Circuit breakers, retries, and dead-letter queues everywhere
8. **Mobile First** — 80% of users will access TribeCart via mobile; design for thumbs

---

## 6. Platform Architecture Overview

### 6.1 High-Level Architecture

```
[Customer Web]  [Seller Web]  [Admin Web]  [Mobile App]
      |               |             |             |
      └───────────────┴─────────────┴─────────────┘
                             │
                    ┌────────▼───────┐
                    │  API Gateway   │   HTTP/REST (public)
                    │  (Go/Gorilla)  │   JWT Auth, Rate Limiting
                    └────────┬───────┘
                             │ gRPC (internal)
         ┌───────────────────┼──────────────────────┐
         │                   │                      │
   ┌─────▼─────┐      ┌──────▼──────┐      ┌───────▼──────┐
   │  Users    │      │  Products   │      │   Orders     │
   │  Service  │      │  Service    │      │   Service    │
   └─────┬─────┘      └──────┬──────┘      └───────┬──────┘
         │                   │                      │
   ┌─────▼─────┐      ┌──────▼──────┐      ┌───────▼──────┐
   │ Payments  │      │  Search     │      │  Inventory   │
   │  Service  │      │  Service    │      │   Service    │
   └─────┬─────┘      └──────┬──────┘      └───────┬──────┘
         │                   │                      │
   ┌─────▼────────────────────────────────────────────────┐
   │              Shared Infrastructure                    │
   │  PostgreSQL | Redis | S3 | ElasticSearch | Kafka      │
   └───────────────────────────────────────────────────────┘
```

### 6.2 Technology Stack

| Layer | Technology | Rationale |
|---|---|---|
| **API Gateway** | Go + Gorilla Mux | Lightweight, high-performance HTTP routing |
| **Microservices** | Go 1.24 | Statically-typed, fast, concurrent, excellent gRPC support |
| **Service Communication** | gRPC + Protocol Buffers | Type-safe, efficient, auto-generates client code |
| **Frontend (Admin/Seller)** | Next.js 15 + React 19 | SSR, excellent DX, App Router |
| **Frontend (Customer)** | Next.js 15 | SEO-critical, SSR for product pages |
| **State Management** | Zustand | Lightweight, minimal boilerplate |
| **Server State / Fetching** | TanStack Query v5 | Caching, background sync, revalidation |
| **Primary Database** | PostgreSQL 16 (on Render Managed DB) | ACID-compliant, high-performance managed relational DB |
| **Caching Layer** | Redis (on Render Managed Redis) | session management, rate limiting, cart state |
| **Search Engine** | Elasticsearch (Elastic Cloud) | Full-text product search, faceted filtering |
| **Message Broker** | Apache Kafka (Upstash) | Async event streaming between services (serverless Kafka) |
| **Object Storage** | Cloudflare R2 / S3-compatible | Product images, documents, exports (platform agnostic) |
| **CDN** | Cloudflare | Global asset delivery, image resizing at edge |
| **Compute Platform** | **Render (Web & Private Services)** | Native Go/Next.js support, zero-downtime, auto-scaling |
| **Build System** | Turbo (pnpm) | Monorepo task orchestration |
| **IaC** | **Render Blueprints (`render.yaml`)** | Reproducible, infrastructure-as-code for Render environments |

---

## 7. Feature Specifications by Domain

---

### 7.1 User & Identity Service

**Service Name:** `users-service`  
**gRPC Port:** `8081`  
**Responsible Team:** Platform Core

#### 7.1.1 Entities

```protobuf
// Full specification — see proto/tribecart/v1/users.proto

User {
  id, first_name, last_name, email, phone_number,
  role (CUSTOMER | SELLER | ADMIN | SUPER_ADMIN),
  status (ACTIVE | PENDING | SUSPENDED | DELETED),
  email_verified, phone_verified,
  addresses[], metadata{},
  created_at, updated_at, last_login_at
}

Address {
  id, user_id, label, line1, line2,
  city, state, postal_code, country,
  is_default, created_at, updated_at
}
```

#### 7.1.2 Features

| Feature | Priority | Description |
|---|---|---|
| **Registration** | P0 | Email + password. Auto-assigns `CUSTOMER` role. Sends verification email |
| **Login** | P0 | Returns JWT access + refresh token pair. 15-min access token, 7-day refresh |
| **Refresh Token** | P0 | Silent token rotation; old refresh token invalidated on use |
| **Forgot / Reset Password** | P0 | Time-limited token (15 min) sent via email |
| **Email Verification** | P0 | Required before first purchase or product listing |
| **Phone Verification** | P1 | OTP via SMS, required for seller onboarding |
| **Address Book** | P1 | Multiple addresses, one default, full CRUD |
| **Seller Registration** | P1 | Upgrade `CUSTOMER` → `SELLER` after KYC verification |
| **OAuth** | P2 | Google, Facebook sign-in as alternatives |
| **2FA / TOTP** | P2 | Time-based OTP for admin accounts |
| **Session Management** | P1 | View and revoke active sessions |

#### 7.1.3 Security Requirements

- Passwords hashed with **argon2id** (never bcrypt for new implementations)
- JWT signed with **RS256** (asymmetric keys); public key exposed at `/.well-known/jwks.json`
- Refresh tokens stored as **hashed** values in PostgreSQL; raw token is single-use
- Rate limit: 5 failed login attempts → 15-minute account lockout
- All PII encrypted at rest using AES-256-GCM
- GDPR/DPDP compliant: data deletion endpoint, export endpoint

#### 7.1.4 Database Schema (PostgreSQL)

```sql
-- users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(20),
    password_hash TEXT NOT NULL,  -- argon2id
    role user_role NOT NULL DEFAULT 'CUSTOMER',
    status user_status NOT NULL DEFAULT 'PENDING',
    email_verified BOOLEAN NOT NULL DEFAULT false,
    phone_verified BOOLEAN NOT NULL DEFAULT false,
    last_login_at TIMESTAMPTZ,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- refresh_tokens table
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,  -- SHA-256 of the raw token
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- addresses table  
CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label VARCHAR(50),
    line1 TEXT NOT NULL,
    line2 TEXT,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL DEFAULT 'IN',
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

### 7.2 Product Catalog Service

**Service Name:** `products-service`  
**gRPC Port:** `50051`  
**Responsible Team:** Catalogue

#### 7.2.1 Core Product Model

```
Product
├── Identity: id, sku, barcode, seller_id
├── Content: name, description (HTML), seo_title, seo_description, seo_keywords
├── Media: images[], videos[]
├── Pricing: price, sale_price, cost_price, currency
├── Inventory: stock_quantity, track_inventory, min_stock_level
├── Dimensions: weight, length, width, height
├── Classification: type, status, is_featured, is_visible
├── Tax: is_taxable, tax_class_id
├── Shipping: requires_shipping, shipping_class_id
├── Relations: category_ids[], tags[], specifications{}, metadata{}
└── Variants → ProductVariant[]
       └── sku, barcode, price_adjustment, stock_quantity, attributes{}, images[]
```

#### 7.2.2 Features

| Feature | Priority | Description |
|---|---|---|
| **Create/Update/Delete Product** | P0 | Full CRUD for sellers |
| **Product Variants** | P0 | Size, color, material — each variant has independent stock |
| **Rich Media** | P0 | Up to 9 images, 2 videos per product; processed through S3 + CloudFront |
| **Category Taxonomy** | P0 | Hierarchical (up to 5 levels), admin-managed |
| **Bulk Import** | P1 | CSV/JSON upload, up to 10,000 SKUs per batch |
| **Bulk Export** | P1 | Export catalog for external use |
| **Soft Delete** | P0 | Products never hard-deleted; `deleted_at` timestamp set |
| **Draft Mode** | P0 | `DRAFT` status allows listing without going live |
| **SEO Fields** | P1 | Custom title, meta description, keywords per product |
| **Product Specifications** | P1 | Key-value pairs (e.g., "Material": "Cotton", "Origin": "India") |
| **Related Products** | P2 | ML-suggested or seller-defined related items |
| **Product Questions & Answers** | P2 | Community Q&A on product pages |

#### 7.2.3 Image Processing Pipeline

```
Seller Upload → S3 Raw Bucket → Lambda Trigger → Image Processing
  → Generate: thumbnail (100x100), small (400x400), large (800x800), original
  → Compress: WebP format, quality 85
  → CDN: CloudFront distribution across 3+ edge nodes
  → Store: public URLs in product.images[]
```

---

### 7.3 Search & Discovery

**Service Name:** `search-service`  
**Engine:** Elasticsearch 8.x  
**Responsible Team:** Discovery

#### 7.3.1 Search Features

| Feature | Priority | Description |
|---|---|---|
| **Full-Text Search** | P0 | Keyword search across name, description, brand, tags |
| **Faceted Filtering** | P0 | Filter by category, price range, rating, brand, attributes |
| **Autocomplete / Typeahead** | P0 | Prefix suggestions as user types, under 50ms |
| **Spell Correction** | P1 | "Nkie shoes" → "Nike shoes" |
| **Synonym Expansion** | P1 | "phone" matches "mobile", "smartphone", "handset" |
| **Personalized Ranking** | P2 | ML model re-ranks results based on user history |
| **Trending Searches** | P1 | Display top 10 trending searches per region |
| **Search Analytics** | P1 | Track queries, click-through rates, zero-result queries |
| **Category Browse** | P0 | Hierarchical category drill-down independent of text search |
| **Saved Searches** | P2 | Users can save a search query and receive alerts |

#### 7.3.2 Elasticsearch Index Schema

```json
{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "name": { "type": "text", "analyzer": "english", "boost": 3 },
      "description": { "type": "text", "analyzer": "english" },
      "brand": { "type": "keyword", "copy_to": "suggest" },
      "tags": { "type": "keyword" },
      "category_ids": { "type": "keyword" },
      "price": { "type": "double" },
      "sale_price": { "type": "double" },
      "rating": { "type": "float" },
      "review_count": { "type": "integer" },
      "stock_quantity": { "type": "integer" },
      "status": { "type": "keyword" },
      "seller_id": { "type": "keyword" },
      "is_featured": { "type": "boolean" },
      "images": { "type": "keyword", "index": false },
      "specifications": { "type": "object" },
      "suggest": { "type": "completion" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}
```

#### 7.3.3 Sync Strategy

- `products-service` publishes a `product.upserted` or `product.deleted` event to Kafka on every change
- `search-service` consumes these events and updates the Elasticsearch index asynchronously
- Maximum indexing lag: **5 seconds** in P99

---

### 7.4 Shopping Cart & Checkout

**Service Name:** `cart-service`  
**Storage:** Redis (ephemeral), PostgreSQL (persistent)  
**Responsible Team:** Commerce

#### 7.4.1 Cart Features

| Feature | Priority | Description |
|---|---|---|
| **Guest Cart** | P0 | Cart persists via cookie/session token before login |
| **Authenticated Cart** | P0 | Cart tied to user account; persists across devices |
| **Guest → Auth Cart Merge** | P0 | On login, guest cart items merged with user cart |
| **Add / Update / Remove Items** | P0 | Variant-aware; real-time stock validation on add |
| **Cart Expiry** | P1 | Abandoned cart expires after 30 days; email reminder at 1 hour |
| **Save for Later** | P1 | Move items out of cart to a wishlist-like saved area |
| **Coupon Application** | P1 | Apply promo codes; validate and show discount breakdown |
| **Multi-Seller Cart** | P0 | Single cart can hold items from multiple sellers; grouped at checkout |
| **Price Lock** | P0 | Price shown in cart is honored for 15 minutes after entering checkout |
| **Cart Summary** | P0 | Subtotal, discount, estimated shipping, estimated tax, total |

#### 7.4.2 Checkout Flow

```
1. Review Cart
2. Select / Enter Shipping Address
3. Choose Shipping Method (Standard / Express / Same-Day where available)
4. Select Payment Method (Card / UPI / Wallet / COD / BNPL)
5. Apply Coupon (optional)
6. Review Order Summary (itemized by seller for multi-seller cart)
7. Place Order → triggers payment intent
8. Payment Confirmation → order created, inventory reserved
9. Order Confirmation Page + Email
```

#### 7.4.3 Real-Time Inventory Check

At each step from "Add to Cart" through "Place Order", a stock check must occur. If an item goes out of stock during checkout, the user must be notified before payment is collected. Never charge a customer for an out-of-stock item.

---

### 7.5 Orders Service

**Service Name:** `orders-service`  
**gRPC Port:** `8083`  
**Responsible Team:** Commerce

#### 7.5.1 Order Lifecycle

```
PENDING_PAYMENT
      │
      ▼ (payment confirmed)
  PAYMENT_CONFIRMED
      │
      ▼ (seller accepts)
  PROCESSING
      │
      ▼ (shipped)
  SHIPPED
      │
      ▼ (delivered)
  DELIVERED
      │
      ▼ (return window closed)
  COMPLETED
      
  [CANCELLED] ← from PENDING_PAYMENT or PROCESSING
  [RETURN_REQUESTED] ← from DELIVERED (within 7 days)
  [RETURN_IN_TRANSIT] ← after return pickup
  [REFUNDED] ← after return verified
```

#### 7.5.2 Features

| Feature | Priority | Description |
|---|---|---|
| **Create Order** | P0 | Atomic: deduct stock + create payment intent in one transaction |
| **Multi-Seller Order Split** | P0 | A single customer order splits into sub-orders per seller |
| **Order Detail Page** | P0 | Full item list, statuses, tracking info, invoice |
| **Cancel Order** | P0 | Allowed before SHIPPED; triggers refund automatically |
| **Return Request** | P1 | Within 7 days of delivery; seller notified; courier pickup scheduled |
| **Real-time Tracking** | P1 | Webhook integration with fulfillment partners |
| **Order History** | P0 | Paginated list with filters (date, status, seller) |
| **Order Invoice** | P1 | Auto-generated PDF invoice per order |
| **Seller Order Dashboard** | P0 | View new orders, accept/reject, mark as shipped |

#### 7.5.3 Database Schema

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES users(id),
    status order_status NOT NULL DEFAULT 'PENDING_PAYMENT',
    currency CHAR(3) NOT NULL DEFAULT 'INR',
    subtotal BIGINT NOT NULL,        -- all amounts in paise (1 INR = 100 paise)
    discount_amount BIGINT NOT NULL DEFAULT 0,
    shipping_amount BIGINT NOT NULL DEFAULT 0,
    tax_amount BIGINT NOT NULL DEFAULT 0,
    total_amount BIGINT NOT NULL,
    shipping_address_id UUID REFERENCES addresses(id),
    payment_method VARCHAR(50),
    payment_id VARCHAR(255),
    coupon_code VARCHAR(100),
    notes TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES users(id),
    product_id UUID NOT NULL,
    variant_id UUID,
    product_name VARCHAR(500) NOT NULL,   -- snapshot at time of order
    variant_name VARCHAR(255),
    unit_price BIGINT NOT NULL,           -- in paise
    quantity INTEGER NOT NULL,
    discount_amount BIGINT NOT NULL DEFAULT 0,
    status order_item_status NOT NULL DEFAULT 'PENDING',
    tracking_number VARCHAR(255),
    carrier VARCHAR(100),
    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ
);
```

---

### 7.6 Payments Service

**Service Name:** `payments-service`  
**gRPC Port:** `8084`  
**Payment Gateways:** Razorpay (primary), Stripe (international)  
**Responsible Team:** Finance

#### 7.6.1 Supported Payment Methods

| Method | Status | Notes |
|---|---|---|
| **Credit / Debit Card** | P0 | Visa, Mastercard, Rupay via Razorpay |
| **UPI** | P0 | Real-time via Razorpay (GPay, PhonePe, BHIM, etc.) |
| **Net Banking** | P1 | 50+ banks via Razorpay |
| **Wallets** | P1 | Paytm, Amazon Pay, etc. via Razorpay |
| **Cash on Delivery** | P1 | Available on eligible orders (order total < ₹10,000) |
| **EMI** | P2 | No-cost EMI via Razorpay |
| **Buy Now Pay Later** | P2 | Integration with Simpl, LazyPay |
| **International Cards** | P2 | Stripe for non-INR orders |

#### 7.6.2 Payment Flow

```
Customer confirms order
  → payments-service creates PaymentIntent with Razorpay
  → Frontend uses Razorpay Checkout SDK
  → Customer completes payment
  → Razorpay sends webhook to payments-service
  → payments-service verifies webhook signature
  → payments-service publishes payment.succeeded event to Kafka
  → orders-service updates order status to PAYMENT_CONFIRMED
  → inventory-service confirms stock reservation
  → notifications-service sends confirmation email/SMS
```

#### 7.6.3 Seller Payouts

- Seller payouts occur **T+2 days** after delivery confirmation
- Payout held if return request is active
- Seller payout dashboard shows: Pending, Cleared, On Hold balances
- All payout records maintain a full ledger with reason codes

#### 7.6.4 Refunds

| Trigger | Timeline | Method |
|---|---|---|
| Cancelled order (before shipping) | T+0 (instant) | Original payment method |
| Cancelled order (after shipping) | T+2 (after return) | Original payment method |
| Return accepted | T+3 (after delivery to seller) | Original payment method or wallet credit |
| Seller rejection | T+1 (after customer dispute) | Original payment method |

---

### 7.7 Inventory & Fulfillment

**Service Name:** `inventory-service`  
**Responsible Team:** Fulfilment

#### 7.7.1 Inventory Model

```
Each ProductVariant has:
  - reserved_quantity  (items in active orders, not yet shipped)
  - available_quantity (sellable = total - reserved)
  - total_quantity     (physical stock)
```

#### 7.7.2 Features

| Feature | Priority | Description |
|---|---|---|
| **Real-time Stock** | P0 | Stock atomically decremented on payment confirmation |
| **Soft Reservation** | P0 | Stock reserved at checkout start, released if payment fails after 15 min |
| **Low Stock Alerts** | P1 | Email/dashboard alert when below `min_stock_level` |
| **Stock Movement Log** | P1 | Audit trail of every stock change with reason code |
| **Restock Request** | P2 | Sellers can raise restock requests with notes |
| **Multi-Warehouse** | P3 | Support for sellers with multiple fulfillment centers |

#### 7.7.3 Fulfillment Partners (Phase 2)

- Shiprocket (primary)
- Delhivery
- Blue Dart (premium)
- India Post (rural)
- Direct seller fulfillment (no third-party)

---

### 7.8 Seller Portal

**App:** `apps/seller-web` (Next.js)  
**Port:** `3001` in Docker  
**Responsible Team:** Seller Experience

#### 7.8.1 Sections

| Section | Features |
|---|---|
| **Dashboard** | GMV, orders today/week, traffic, conversion rate, revenue trend chart |
| **Products** | List, create, edit, bulk import/export, manage variants, manage images |
| **Orders** | View new orders, accept, mark as packed, mark as shipped, track, handle returns |
| **Inventory** | Current stock levels, movement history, set min stock alerts |
| **Analytics** | Sales by product, by category, by geography, customer lifetime value |
| **Promotions** | Create coupons, flash sales, bundle deals |
| **Payouts** | Balance overview, transaction history, bank account management |
| **Reviews** | Read and respond to customer reviews (no editing/deleting) |
| **Settings** | Shop profile, logo, banner, return policy, shipping zones |
| **Onboarding** | Step-by-step wizard: KYC → Bank details → First product → First sale |

#### 7.8.2 KYC Requirements (India)

- PAN Card (mandatory for all sellers)
- GST Number (mandatory if annual revenue > ₹20 Lakhs)
- Bank Account with IFSC (for payouts)
- Business registration documents (optional but increases trust badge)

---

### 7.9 Admin Portal

**App:** `apps/admin-web` (Next.js)  
**Port:** `3000` in Docker  
**Responsible Team:** Operations

#### 7.9.1 Sections

| Section | Features |
|---|---|
| **Dashboard** | Platform GMV, active users, active sellers, orders per minute, system health |
| **Users** | Search, view, edit, suspend, delete; view purchase history |
| **Sellers** | Approve/reject KYC, suspend sellers, view seller metrics |
| **Products** | Moderate listings, approve/reject products, manage categories |
| **Orders** | Override order statuses, initiate refunds, manage disputes |
| **Payments** | View all payment transactions, trigger manual refunds |
| **Promotions** | Create platform-wide discount campaigns |
| **Categories** | Manage the category tree (add, edit, reorder, icons) |
| **Analytics** | Full platform analytics: cohort analysis, funnel reports, revenue attribution |
| **System** | Feature flags, configuration values, audit log |
| **Content** | Manage homepage banners, featured products, curated lists |

#### 7.9.2 Role-Based Access Control (RBAC)

| Role | Permissions |
|---|---|
| `SUPER_ADMIN` | All permissions, can create/delete other admins |
| `ADMIN` | All except system config and admin management |
| `MODERATOR` | Products, reviews, user flags only |
| `FINANCE` | Payments, payouts, refunds only |
| `ANALYST` | Read-only access to all analytics |

---

### 7.10 Customer Portal

**App:** `apps/customer-web` (Next.js)  
**Port:** `3002` in Docker  
**Responsible Team:** Customer Experience

#### 7.10.1 Pages

| Page | Features |
|---|---|
| **Homepage** | Hero banners, featured categories, trending products, personalized recommendations |
| **Search Results** | Full-text search, facets (category, price, brand, rating), sort options |
| **Category Browse** | Grid/List view, filter panel, sort options |
| **Product Detail Page** | Multiple images, description, specifications, variants, reviews, Q&A, seller info |
| **Cart** | Live price, coupon field, shipping estimate, order summary |
| **Checkout** | Address selection, shipping method, payment selection, order review |
| **Order Confirmation** | Thank you page with order ID and estimated delivery |
| **My Orders** | List with status, detail view, tracking, cancel/return actions |
| **Profile** | Edit personal info, manage addresses, change password |
| **Wishlist** | Save products; share wishlist via link |
| **Tribe / Community** | Follow sellers, see their new listings, tribe-specific deals |

#### 7.10.2 Performance Requirements

| Metric | Target |
|---|---|
| Largest Contentful Paint (LCP) | < 2.5 seconds |
| First Input Delay (FID) | < 100ms |
| Cumulative Layout Shift (CLS) | < 0.1 |
| Time to First Byte (TTFB) | < 200ms |
| Core Web Vitals | All Green |
| Lighthouse Score | > 90 (all categories) |

---

### 7.11 Reviews & Ratings

**Service Name:** `reviews-service`  
**Responsible Team:** Trust & Safety

#### 7.11.1 Features

| Feature | Priority | Description |
|---|---|---|
| **Star Rating** | P0 | 1–5 stars, verified purchasers only |
| **Written Review** | P0 | 20–2000 characters, post-purchase only, 48-hour hold before publish |
| **Review Photos** | P1 | Up to 5 photos per review |
| **Helpful Votes** | P1 | "Was this helpful?" yes/no |
| **Seller Response** | P1 | Sellers can respond publicly (once per review) |
| **Review Moderation** | P0 | Profanity filter, spam detection, admin manual review |
| **Verified Purchase Badge** | P0 | Only shown for reviews from confirmed buyers |
| **Rating Breakdown** | P0 | 5-star histogram on product page |
| **Fake Review Detection** | P2 | ML model flags suspicious patterns |

---

### 7.12 Notifications Service

**Service Name:** `notifications-service`  
**Channels:** Email (SendGrid), SMS (Twilio/MSG91), Push (Firebase FCM), In-App  
**Responsible Team:** Platform Core

#### 7.12.1 Notification Events

| Event | Channel | Recipient |
|---|---|---|
| Account created | Email | Customer |
| Email verification | Email | Customer |
| Password reset | Email | Customer |
| Order placed | Email + SMS | Customer + Seller |
| Payment confirmed | Email + SMS + Push | Customer |
| Order shipped | Email + SMS + Push | Customer |
| Order delivered | Email + Push | Customer |
| Return approved | Email + SMS | Customer |
| Refund processed | Email | Customer |
| New order (seller) | Email + SMS + Push | Seller |
| Low stock alert | Email + In-App | Seller |
| Payout processed | Email | Seller |
| KYC approved/rejected | Email | Seller |
| New review on product | In-App | Seller |
| System alert | Email + Slack | Admin |

#### 7.12.2 Template Engine

- Templates stored in database, editable by admin without code deployment
- Variables: `{{user.first_name}}`, `{{order.id}}`, `{{order.total}}`, etc.
- A/B testing support for subject lines and call-to-action variants
- Unsubscribe links required in all marketing emails (CAN-SPAM / DPDP compliant)

---

### 7.13 Analytics & Reporting

**Service Name:** `analytics-service`  
**Storage:** ClickHouse (OLAP) or Redshift  
**Responsible Team:** Data

#### 7.13.1 Seller Analytics

| Report | Metrics | Granularity |
|---|---|---|
| **Sales Overview** | GMV, orders, units sold, ASP | Hourly, Daily, Weekly, Monthly |
| **Product Performance** | Views, add-to-cart, purchases, conversion rate | Per product |
| **Traffic Sources** | Direct, search, referral, social | By channel |
| **Customer Geography** | Orders by city, state | Heat map |
| **Review Analytics** | Average rating over time, response rate | Monthly |
| **Return Rate** | Returns / orders, top return reasons | Monthly |

#### 7.13.2 Admin Analytics

- Platform GMV by day/month/year
- Active buyers vs. churned buyers (cohort analysis)
- Top sellers by GMV, by growth rate
- Category performance
- Search query analytics (top searches, zero-result searches)
- Funnel: Homepage → Search → PDP → Cart → Checkout → Purchase
- Fraud signals: order velocity, device fingerprint anomalies

---

### 7.14 Promotions & Discounts

**Service Name:** `promotions-service`  
**Responsible Team:** Growth

#### 7.14.1 Coupon Types

| Type | Description | Example |
|---|---|---|
| **Flat Discount** | Fixed amount off | Save ₹100 |
| **Percentage Discount** | % off subtotal | 15% off |
| **Free Shipping** | Waive shipping fee | Free shipping on orders > ₹499 |
| **Buy X Get Y** | Auto-add free item | Buy 2 get 1 free |
| **Category Coupon** | Applies only to specific categories | 20% off Electronics |
| **First Order** | New users only | ₹200 off first order |
| **Seller Coupon** | Issued by seller, funded by seller | Seller absorbs discount |
| **Platform Coupon** | Issued by TribeCart, platform absorbs cost | Seasonal sale |

#### 7.14.2 Rules Engine

```
Coupon Entity {
  code: string (unique, alphanumeric)
  type: flat | percentage | free_shipping | bxgy
  value: number
  min_order_value: number
  max_discount_cap: number (for percentage type)
  applicable_user_ids: [] (empty = all users)
  applicable_category_ids: [] (empty = all categories)
  applicable_seller_ids: [] (empty = all sellers)
  usage_limit_total: number
  usage_limit_per_user: number
  starts_at: timestamp
  expires_at: timestamp
  is_stackable: boolean (can be combined with other coupons?)
}
```

---

## 8. Non-Functional Requirements

### 8.1 Performance

| Metric | Requirement |
|---|---|
| API Gateway P50 latency | < 50ms |
| API Gateway P99 latency | < 200ms |
| Product search P99 latency | < 100ms |
| Checkout page load | < 1.5 seconds |
| Order creation API | < 500ms end-to-end |
| Max concurrent users (Day 1) | 100,000 |
| Max concurrent users (Year 1 target) | 1,000,000 |

### 8.2 Availability

| Service | SLA Target |
|---|---|
| API Gateway | 99.99% (< 1hr downtime/year) |
| Orders Service | 99.99% |
| Payments Service | 99.999% (< 5min downtime/year) |
| Products Service | 99.9% |
| Search Service | 99.9% |
| Admin Portal | 99.5% |

### 8.3 Scalability

- All services must be stateless and horizontally scalable
- **Render Auto-scaling** configured for all public and private services based on CPU/Memory usage
- Database connection pooling managed via Render DB connection limits or `pg-bouncer` sidecar if needed
- Read replicas for PostgreSQL for all read-heavy queries (product browsing, order history)

### 8.4 Observability

| Component | Tool |
|---|---|
| Distributed Tracing | Jaeger / OpenTelemetry |
| Metrics | Prometheus + Grafana |
| Log Aggregation | ELK Stack (Elasticsearch, Logstash, Kibana) |
| Alerting | PagerDuty + Grafana Alerting |
| Uptime Monitoring | Pingdom / Uptime Robot |
| Error Tracking | Sentry |

---

## 9. API Gateway & Communication Contracts

### 9.1 REST API Conventions

All public facing APIs follow REST conventions:

```
Base URL: https://api.tribecart.com/v1

Authentication: Bearer JWT in Authorization header
Content-Type: application/json
Date Format: ISO 8601 (2026-03-27T00:00:00Z)
Currency: All monetary values in paise (integer), display layer converts to INR
Pagination: { data: [], page: 1, page_size: 20, total_count: 1000 }
Error Format: { error: { code: "PRODUCT_NOT_FOUND", message: "...", details: {} } }
```

### 9.2 Rate Limiting

| Consumer | Limit |
|---|---|
| Unauthenticated | 30 req/min |
| Authenticated Customer | 120 req/min |
| Authenticated Seller | 300 req/min |
| Admin | 600 req/min |
| Webhooks (inbound) | 1000 req/min |

### 9.3 Internal gRPC Communication

All inter-service calls use gRPC with mTLS (mutual TLS). Service discovery via **Render Private Networking** (e.g., `user-service:8081` within the private network).

---

## 10. Database Strategy

### 10.1 Per-Service Database Isolation

Each service owns its own database. No cross-service database queries. Data sharing happens through events (Kafka) or gRPC calls.

| Service | Database | Why |
|---|---|---|
| users-service | PostgreSQL `tribecart_users` | ACID, relational |
| products-service | PostgreSQL `tribecart_products` | ACID, complex queries |
| orders-service | PostgreSQL `tribecart_orders` | ACID, financial integrity |
| payments-service | PostgreSQL `tribecart_payments` | ACID, financial integrity |
| cart-service | Redis | Speed, ephemeral |
| search-service | Elasticsearch | Full-text, faceted search |
| analytics-service | ClickHouse | OLAP, columnar |
| sessions | Redis | Fast lookups, TTL |

### 10.2 Migration Strategy

- Use **golang-migrate** for all SQL migrations
- Migrations run at service startup (idempotent)
- Never drop columns in production; mark as deprecated, remove in next major version
- All migrations version-controlled in `services/{name}/migrations/`

---

## 11. Security & Compliance

### 11.1 Application Security

| Control | Implementation |
|---|---|
| **API Authentication** | JWT RS256 with 15-min expiry |
| **Password Hashing** | argon2id with time=1, memory=64MB, threads=4 |
| **Transport Security** | TLS 1.3 minimum; mTLS between internal services |
| **Input Validation** | Strict protobuf type validation; sanitize all HTML input |
| **SQL Injection** | Parameterized queries only; no string-concatenation SQL |
| **CORS** | Allowlist only; `*` forbidden in production |
| **Rate Limiting** | Redis-backed token bucket algorithm |
| **CSRF** | Double-submit cookie pattern for all state-changing requests |
| **Secrets Management** | **Render Secret Groups**; never in code or env files; synced across services |
| **Dependency Scanning** | Dependabot + OWASP Dependency Check in CI |
| **Container Scanning** | Trivy scan in CI pipeline |

### 11.2 Data Privacy (DPDP Act 2023 / GDPR)

- User data export endpoint: `GET /v1/me/data-export`
- User data deletion endpoint: `DELETE /v1/me` (soft delete + anonymization)
- Data retention: User PII deleted 2 years after account closure
- Cookie consent banner on first visit
- All analytics data anonymized (no PII in ClickHouse)
- Privacy Policy and Terms reviewed by legal counsel before launch

### 11.3 Payment Security

- PCI-DSS compliance: TribeCart never stores raw card data; Razorpay is PCI-DSS Level 1
- Webhook signatures verified using HMAC-SHA256
- All payment events idempotent (webhook replay protection)

---

## 12. DevOps & Infrastructure

### 12.1 Environments

| Environment | Purpose | Deployment Trigger |
|---|---|---|
| `local` | Developer machines | Manual (`docker compose up`) |
| `dev` | Shared integration testing | Push to `develop` branch |
| `staging` | Pre-production verification | Push to `staging` branch |
| `production` | Live traffic | Manual promote from staging |

### 12.2 CI/CD Pipeline (GitHub Actions)

```
PR Created → Lint → Test → Build → Docker Build → Push to ECR
              ↓
        Merge to develop → Deploy to dev → Smoke tests
              ↓
        Merge to staging → Deploy to staging → E2E tests → Notify QA
              ↓
        Manual approve → Deploy to production → Progressive rollout (10% → 50% → 100%)
```

### 12.3 Render Deployment Architecture

- **Blueprints:** All infrastructure defined in `render.yaml` (Blueprints)
- **Web Services:** Public-facing API Gateway and Frontend apps (Next.js)
- **Private Services:** Internal microservices (Go) not exposed to the public internet
- **Managed Services:** Managed PostgreSQL and Redis with automated backups and failover
- **Auto-Deploys:** CI/CD triggers on git push to specific branches
- **Health Checks:** Native Render liveness/readiness probes configured for all services

### 12.4 Data Backups

| Data | Backup Frequency | Retention | RTO | RPO |
|---|---|---|---|---|
| PostgreSQL | Continuous (Render Point-in-Time Recovery) | 30 days | < 1 hour | < 5 min |
| Redis | Managed by Render | 7 days | < 30 min | < 1 hour |
| S3/R2 (media) | Versioning enabled | Indefinite | < 1 min | 0 |
| Elasticsearch | Managed snapshots to S3 | 14 days | < 2 hours | < 24 hours |

---

## 13. Phased Roadmap

### Phase 1 — Foundation (Months 1–3)
> **Goal:** Core transactional engine working end-to-end

- [ ] Fix protocol drift (users & products service match proto definitions)
- [ ] Implement password hashing (argon2id) in users-service
- [ ] JWT authentication (RS256) in API Gateway
- [ ] Products CRUD (full proto compliance)
- [ ] Cart service (Redis)
- [ ] Orders service (basic lifecycle)
- [ ] Payments service (Razorpay integration)
- [ ] Customer Web: Homepage, PDP, Cart, Checkout, Order Confirmation
- [ ] Seller Web: Product management, Order management
- [ ] Admin Web: User management, Product moderation
- [ ] Email notifications (order placed, shipped, delivered)
- [ ] Docker Compose local dev working for all services

### Phase 2 — Growth (Months 4–6)
> **Goal:** Seller empowerment and customer retention

- [ ] Elasticsearch integration + product search
- [ ] Reviews & ratings system
- [ ] Promotions / coupons engine
- [ ] Seller analytics dashboard
- [ ] Return/refund flow
- [ ] Bulk product import/export
- [ ] Push notifications (FCM)
- [ ] Wishlist
- [ ] Multi-warehouse inventory

### Phase 3 — Scale (Months 7–12)
> **Goal:** Performance, personalization, trust

- [ ] **Render Production Scale Deployment** (Auto-scaling, mTLS)
- [ ] Personalized product recommendations (ML model)
- [ ] Fake review detection (ML model)
- [ ] Kafka event streaming implementation
- [ ] ClickHouse analytics pipeline
- [ ] International payments (Stripe)
- [ ] Mobile App (React Native)
- [ ] Seller mobile app
- [ ] Affiliate / referral program
- [ ] API program for Power Sellers

### Phase 4 — Ecosystem (Year 2+)
> **Goal:** Platform becomes an ecosystem

- [ ] TribeCart Marketplace API (public)
- [ ] Third-party logistics integrations (Shiprocket, Delhivery)
- [ ] Financial products for sellers (working capital loans)
- [ ] Live commerce (streaming + shoppable video)
- [ ] Social commerce features (Tribe communities)
- [ ] AI-powered seller assistant

---

## 14. Success Metrics & KPIs

### Business KPIs

| KPI | Month 3 | Month 6 | Month 12 |
|---|---|---|---|
| Registered Sellers | 500 | 2,000 | 10,000 |
| Active Customers | 5,000 | 25,000 | 200,000 |
| GMV / Month | ₹10L | ₹1Cr | ₹10Cr |
| Orders / Day | 100 | 1,000 | 10,000 |
| Average Order Value | ₹800 | ₹900 | ₹1,000 |
| Platform Take Rate | 8% | 7% | 6% |

### Product Health KPIs

| KPI | Target |
|---|---|
| Customer NPS | > 50 |
| Seller NPS | > 60 |
| Cart Abandonment Rate | < 65% |
| Checkout Conversion | > 60% |
| Order Defect Rate (ODR) | < 1% |
| Seller Onboarding Time | < 30 minutes |
| First Listing Time | < 1 hour |
| Support Ticket Resolution | < 24 hours |

### Technical KPIs

| KPI | Target |
|---|---|
| API Uptime | > 99.9% |
| P99 API Latency | < 200ms |
| Deploy Frequency | Daily |
| Change Failure Rate | < 5% |
| MTTR (Mean Time to Recovery) | < 30 minutes |
| Test Coverage | > 80% |

---

## 15. Risks & Mitigations

| Risk | Severity | Probability | Mitigation |
|---|---|---|---|
| **Protocol Drift** (services diverge from protos) | High | High (existing) | Proto-first development mandate; CI check that services compile against proto |
| **Payment Failure Cascade** | Critical | Low | Circuit breakers; idempotent payment intents; dead-letter queues |
| **Inventory Oversell** | High | Medium | Optimistic locking on stock decrement; distributed locks via Redis |
| **Data Breach** | Critical | Low | argon2id passwords; AES-256 PII; Vault for secrets; quarterly pen tests |
| **Seller Fraud** | High | Medium | KYC enforcement; payout hold period; ML anomaly detection |
| **Traffic Spikes** | High | High (sales events) | **Render Auto-scaling**; Redis caching; CDN for static assets |
| **Search Quality** | Medium | Medium | A/B test ranking algorithms; monitor zero-result rate weekly |
| **Scope Creep** | Medium | High | Strict Phase gates; no Phase 2 work until Phase 1 is stable |
| **Vendor Lock-in** | Medium | Low | Abstract payment gateway behind interface; multi-cloud-ready configs |

---

## 16. Glossary

| Term | Definition |
|---|---|
| **GMV** | Gross Merchandise Value — total value of all products sold on the platform |
| **ASP** | Average Selling Price — GMV / number of orders |
| **PDPAge** | Product Detail Page — the page showing a single product's full information |
| **SKU** | Stock Keeping Unit — a unique identifier for each distinct product variant |
| **ODR** | Order Defect Rate — % of orders resulting in a defect (negative review, refund, chargeback) |
| **Paise** | Smallest currency unit in India (100 paise = 1 INR). All internal amounts stored in paise |
| **Proto / Protobuf** | Protocol Buffers — Google's language-agnostic data serialization format used for gRPC |
| **gRPC** | Google Remote Procedure Call — the internal communication protocol between TribeCart microservices |
| **JWT** | JSON Web Token — the authentication token format used for API access |
| **Render Blueprints** | Infrastructure-as-Code for Render (`render.yaml`); defines the entire environment |
| **Private Service** | A Render service only accessible within the private network (for internal gRPC) |
| **KYC** | Know Your Customer — the identity verification process for sellers |
| **DPDP** | Digital Personal Data Protection Act 2023 — India's primary data privacy law |
| **SLA** | Service Level Agreement — the committed uptime/performance guarantee for a service |
| **Tribe** | TribeCart's term for a community of buyers/sellers organized around a shared interest or seller |
| **Take Rate** | The % of GMV that TribeCart keeps as platform revenue |
| **MTTR** | Mean Time to Recovery — average time to restore service after an incident |

---

*This document is the single source of truth for TribeCart's product development. Every feature, every API, every database table should trace back to a requirement in this document. If it's not here, build it here first.*

---

**Next Action for Engineering:** Fix Protocol Drift (Phase 1, Item 1). See `audit_report.md` for full details.
