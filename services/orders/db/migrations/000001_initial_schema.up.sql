-- Create enum types
CREATE TYPE order_status AS ENUM (
    'ORDER_STATUS_PENDING',
    'ORDER_STATUS_PROCESSING',
    'ORDER_STATUS_SHIPPED',
    'ORDER_STATUS_DELIVERED',
    'ORDER_STATUS_CANCELLED',
    'ORDER_STATUS_REFUNDED',
    'ORDER_STATUS_FAILED'
);

CREATE TYPE payment_method AS ENUM (
    'PAYMENT_METHOD_CREDIT_CARD',
    'PAYMENT_METHOD_DEBIT_CARD',
    'PAYMENT_METHOD_PAYPAL',
    'PAYMENT_METHOD_BANK_TRANSFER',
    'PAYMENT_METHOD_CRYPTO'
);

-- Create tables
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    status order_status NOT NULL DEFAULT 'ORDER_STATUS_PENDING',
    subtotal DECIMAL(10, 2) NOT NULL,
    tax_amount DECIMAL(10, 2) NOT NULL,
    shipping_cost DECIMAL(10, 2) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    payment_method payment_method,
    payment_id VARCHAR(255),
    shipping_address_id UUID NOT NULL,
    tracking_number VARCHAR(100),
    carrier VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    cancelled_at TIMESTAMP WITH TIME ZONE,
    cancelled_reason TEXT,
    metadata JSONB
);

-- Add indexes
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- Order items table
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    tax_amount DECIMAL(10, 2) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    metadata JSONB
);

-- Add indexes for order items
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

-- Order status history table
CREATE TABLE order_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status order_status NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    metadata JSONB
);

-- Add indexes for status history
CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX idx_order_status_history_created_at ON order_status_history(created_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at
CREATE TRIGGER update_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_order_items_updated_at
BEFORE UPDATE ON order_items
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to log status changes
CREATE OR REPLACE FUNCTION log_order_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status != OLD.status THEN
        INSERT INTO order_status_history (order_id, status, notes, metadata)
        VALUES (NEW.id, NEW.status, 
                CASE 
                    WHEN NEW.status = 'ORDER_STATUS_CANCELLED' THEN 
                        COALESCE(NEW.cancelled_reason, 'Order was cancelled')
                    ELSE 
                        'Status changed from ' || OLD.status::TEXT || ' to ' || NEW.status::TEXT
                END,
                jsonb_build_object(
                    'previous_status', OLD.status,
                    'new_status', NEW.status,
                    'updated_at', NOW()
                )
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for status changes
CREATE TRIGGER log_order_status_change_trigger
AFTER UPDATE OF status ON orders
FOR EACH ROW EXECUTE FUNCTION log_order_status_change();

-- Create function to calculate order totals
CREATE OR REPLACE FUNCTION calculate_order_totals()
RETURNS TRIGGER AS $$
DECLARE
    order_subtotal DECIMAL(10, 2);
    order_tax DECIMAL(10, 2);
    order_total DECIMAL(10, 2);
BEGIN
    -- Calculate subtotal, tax, and total for the order
    SELECT 
        COALESCE(SUM(total_amount), 0),
        COALESCE(SUM(tax_amount), 0)
    INTO 
        order_subtotal,
        order_tax
    FROM order_items
    WHERE order_id = NEW.id;
    
    -- Calculate total (subtotal + tax + shipping)
    order_total := order_subtotal + order_tax + COALESCE(NEW.shipping_cost, 0);
    
    -- Update the order with calculated values
    UPDATE orders
    SET 
        subtotal = order_subtotal,
        tax_amount = order_tax,
        total_amount = order_total,
        updated_at = NOW()
    WHERE id = NEW.id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for order item changes
CREATE TRIGGER update_order_totals_trigger
AFTER INSERT OR UPDATE OR DELETE ON order_items
FOR EACH ROW EXECUTE FUNCTION calculate_order_totals();
