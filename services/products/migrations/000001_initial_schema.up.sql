-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(36) PRIMARY KEY,
    seller_id VARCHAR(36),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sku VARCHAR(100),
    barcode VARCHAR(100),
    price DECIMAL(10, 2) NOT NULL,
    sale_price DECIMAL(10, 2),
    cost_price DECIMAL(10, 2),
    stock_quantity INTEGER NOT NULL DEFAULT 0,
    track_inventory BOOLEAN NOT NULL DEFAULT true,
    min_stock_level INTEGER,
    weight DECIMAL(10, 2),
    length DECIMAL(10, 2),
    width DECIMAL(10, 2),
    height DECIMAL(10, 2),
    type INTEGER NOT NULL DEFAULT 1, -- 1: physical, 2: digital, 3: service
    status INTEGER NOT NULL DEFAULT 1, -- 1: draft, 2: active, 3: archived, 4: out of stock, 5: discontinued
    is_featured BOOLEAN NOT NULL DEFAULT false,
    is_visible BOOLEAN NOT NULL DEFAULT true,
    requires_shipping BOOLEAN NOT NULL DEFAULT true,
    is_taxable BOOLEAN NOT NULL DEFAULT true,
    tax_class_id VARCHAR(36),
    seo_title VARCHAR(255),
    seo_description TEXT,
    seo_keywords VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for products table
CREATE INDEX IF NOT EXISTS idx_products_seller_id ON products(seller_id);
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_barcode ON products(barcode);
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);
CREATE INDEX IF NOT EXISTS idx_products_updated_at ON products(updated_at);

-- Create product_categories table for many-to-many relationship
CREATE TABLE IF NOT EXISTS product_categories (
    product_id VARCHAR(36) NOT NULL,
    category_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (product_id, category_id),
    CONSTRAINT fk_product
        FOREIGN KEY (product_id) 
        REFERENCES products(id)
        ON DELETE CASCADE
);

-- Create indexes for product_categories table
CREATE INDEX IF NOT EXISTS idx_product_categories_product_id ON product_categories(product_id);
CREATE INDEX IF NOT EXISTS idx_product_categories_category_id ON product_categories(category_id);

-- Create stock_movements table
CREATE TABLE IF NOT EXISTS stock_movements (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    variant_id VARCHAR(36),
    quantity INTEGER NOT NULL,
    reference_id VARCHAR(100),
    reason VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_stock_movement_product
        FOREIGN KEY (product_id) 
        REFERENCES products(id)
        ON DELETE CASCADE
);

-- Create indexes for stock_movements table
CREATE INDEX IF NOT EXISTS idx_stock_movements_product_id ON stock_movements(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_variant_id ON stock_movements(variant_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_created_at ON stock_movements(created_at);

-- Create product_variants table
CREATE TABLE IF NOT EXISTS product_variants (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    sku VARCHAR(100),
    barcode VARCHAR(100),
    price DECIMAL(10, 2) NOT NULL,
    sale_price DECIMAL(10, 2),
    cost_price DECIMAL(10, 2),
    stock_quantity INTEGER NOT NULL DEFAULT 0,
    track_inventory BOOLEAN NOT NULL DEFAULT true,
    min_stock_level INTEGER,
    weight DECIMAL(10, 2),
    length DECIMAL(10, 2),
    width DECIMAL(10, 2),
    height DECIMAL(10, 2),
    is_default BOOLEAN NOT NULL DEFAULT false,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_product_variant_product
        FOREIGN KEY (product_id) 
        REFERENCES products(id)
        ON DELETE CASCADE
);

-- Create indexes for product_variants table
CREATE INDEX IF NOT EXISTS idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX IF NOT EXISTS idx_product_variants_sku ON product_variants(sku);
CREATE INDEX IF NOT EXISTS idx_product_variants_barcode ON product_variants(barcode);

-- Create product_variant_attributes table
CREATE TABLE IF NOT EXISTS product_variant_attributes (
    variant_id VARCHAR(36) NOT NULL,
    attribute_name VARCHAR(100) NOT NULL,
    attribute_value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (variant_id, attribute_name),
    CONSTRAINT fk_variant_attribute_variant
        FOREIGN KEY (variant_id) 
        REFERENCES product_variants(id)
        ON DELETE CASCADE
);

-- Create indexes for product_variant_attributes table
CREATE INDEX IF NOT EXISTS idx_product_variant_attributes_variant_id ON product_variant_attributes(variant_id);
CREATE INDEX IF NOT EXISTS idx_product_variant_attributes_name_value ON product_variant_attributes(attribute_name, attribute_value);

-- Create product_images table
CREATE TABLE IF NOT EXISTS product_images (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    variant_id VARCHAR(36),
    url VARCHAR(512) NOT NULL,
    alt_text VARCHAR(255),
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_product_image_product
        FOREIGN KEY (product_id) 
        REFERENCES products(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_product_image_variant
        FOREIGN KEY (variant_id) 
        REFERENCES product_variants(id)
        ON DELETE SET NULL
);

-- Create indexes for product_images table
CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images(product_id);
CREATE INDEX IF NOT EXISTS idx_product_images_variant_id ON product_images(variant_id);
CREATE INDEX IF NOT EXISTS idx_product_images_sort_order ON product_images(sort_order);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to update updated_at columns
CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_product_variants_updated_at
BEFORE UPDATE ON product_variants
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to update product stock based on variants
CREATE OR REPLACE FUNCTION update_product_stock_from_variants()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        UPDATE products
        SET stock_quantity = (
            SELECT COALESCE(SUM(stock_quantity), 0)
            FROM product_variants
            WHERE product_id = OLD.product_id
              AND deleted_at IS NULL
        )
        WHERE id = OLD.product_id;
    ELSE
        UPDATE products
        SET stock_quantity = (
            SELECT COALESCE(SUM(stock_quantity), 0)
            FROM product_variants
            WHERE product_id = NEW.product_id
              AND deleted_at IS NULL
        )
        WHERE id = NEW.product_id;
    END IF;
    
    RETURN NULL;
END;
$$ language 'plpgsql';

-- Create triggers to update product stock when variants change
CREATE TRIGGER update_product_stock_on_variant_insert
AFTER INSERT ON product_variants
FOR EACH ROW EXECUTE FUNCTION update_product_stock_from_variants();

CREATE TRIGGER update_product_stock_on_variant_update
AFTER UPDATE OF stock_quantity, deleted_at ON product_variants
FOR EACH ROW EXECUTE FUNCTION update_product_stock_from_variants();

CREATE TRIGGER update_product_stock_on_variant_delete
AFTER DELETE ON product_variants
FOR EACH ROW EXECUTE FUNCTION update_product_stock_from_variants();
