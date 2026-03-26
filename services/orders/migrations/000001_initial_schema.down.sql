-- Drop triggers first to avoid dependency issues
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;
DROP TRIGGER IF EXISTS update_order_items_updated_at ON order_items;
DROP TRIGGER IF EXISTS log_order_status_change ON orders;
DROP TRIGGER IF EXISTS update_order_totals_after_insert ON order_items;
DROP TRIGGER IF EXISTS update_order_totals_after_update ON order_items;
DROP TRIGGER IF EXISTS update_order_totals_after_delete ON order_items;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS log_order_status_change();
DROP FUNCTION IF EXISTS calculate_order_totals();
DROP FUNCTION IF EXISTS cancel_order(UUID, TEXT);

-- Drop tables
DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

-- Drop types
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS payment_method;
