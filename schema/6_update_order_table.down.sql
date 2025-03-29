-- Удаляем внешний ключ fk_menu_item_id
ALTER TABLE order_items
DROP CONSTRAINT fk_menu_item_id;

-- Удаляем внешний ключ fk_prices
ALTER TABLE order_items
DROP CONSTRAINT fk_prices;