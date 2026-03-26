-- Create enum types
CREATE TYPE order_status AS ENUM (
    'PENDING',
    'PROCESSING',
    'SHIPPED',
    'DELIVERED',
    'CANCELLED',
    'REFUNDED'
);

CREATE TYPE payment_method AS ENUM (
    'CREDIT_CARD',
    'DEBIT_CARD',
    'PAYPAL',
    'STRIPE',
    'BANK_TRANSFER',
    'CRYPTO',
    'OTHER'
);

-- Create tables
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    status order_status NOT NULL DEFAULT 'PENDING',
    subtotal DECIMAL(10, 2) NOT NULL,
    tax_amount DECIMAL(10, 2) NOT NULL,
    shipping_cost DECIMAL(10, 2) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    payment_method payment_method NOT NULL,
    payment_id VARCHAR(255),
    shipping_address_id UUID NOT NULL,
    tracking_number VARCHAR(255),
    carrier VARCHAR(100),
    cancelled_at TIMESTAMP WITH TIME ZONE,
    cancelled_reason TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    tax_amount DECIMAL(10, 2) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE TABLE order_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    status order_status NOT NULL,
    notes TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_order_status FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);

-- Create triggers and functions
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_order_items_updated_at
BEFORE UPDATE ON order_items
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to log status changes
CREATE OR REPLACE FUNCTION log_order_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        INSERT INTO order_status_history (order_id, status, notes)
        VALUES (NEW.id, NEW.status, 'Status changed from ' || OLD.status || ' to ' || NEW.status);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER log_order_status_change
AFTER UPDATE OF status ON orders
FOR EACH ROW EXECUTE FUNCTION log_order_status_change();

-- Function to calculate order totals
CREATE OR REPLACE FUNCTION calculate_order_totals()
RETURNS TRIGGER AS $$
BEGIN
    -- Calculate subtotal, tax, and total for the order
    UPDATE orders o
    SET 
        subtotal = COALESCE((
            SELECT SUM(subtotal) 
            FROM order_items 
            WHERE order_id = NEW.order_id
        ), 0),
        tax_amount = COALESCE((
            SELECT SUM(tax_amount) 
            FROM order_items 
            WHERE order_id = NEW.order_id
        ), 0),
        total_amount = COALESCE((
            SELECT SUM(total_amount) 
            FROM order_items 
            WHERE order_id = NEW.order_id
        ), 0) + o.shipping_cost
    WHERE id = NEW.order_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for order items
CREATE TRIGGER update_order_totals_after_insert
AFTER INSERT ON order_items
FOR EACH ROW EXECUTE FUNCTION calculate_order_totals();

CREATE TRIGGER update_order_totals_after_update
AFTER UPDATE ON order_items
FOR EACH ROW EXECUTE FUNCTION calculate_order_totals();

CREATE TRIGGER update_order_totals_after_delete
AFTER DELETE ON order_items
FOR EACH ROW EXECUTE FUNCTION calculate_order_totals();

-- Create function to cancel an order
CREATE OR REPLACE FUNCTION cancel_order(
    p_order_id UUID,
    p_reason TEXT DEFAULT NULL
) 
RETURNS BOOLEAN AS $$
DECLARE
    v_status order_status;
BEGIN
    -- Get current status
    SELECT status INTO v_status
    FROM orders
    WHERE id = p_order_id
    FOR UPDATE;
    
    -- Validate status
    IF v_status = 'CANCELLED' OR v_status = 'REFUNDED' THEN
        RAISE EXCEPTION 'Cannot cancel order with status %', v_status;
    END IF;
    
    -- Update order status
    UPDATE orders
    SET 
        status = 'CANCELLED',
        cancelled_at = NOW(),
        cancelled_reason = p_reason,
        updated_at = NOW()
    WHERE id = p_order_id;
    
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;
