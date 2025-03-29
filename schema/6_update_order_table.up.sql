-- Добавляем внешний ключ для столбца menu_item_id в таблице order_items, который ссылается на таблицу menu
ALTER TABLE order_items
ADD CONSTRAINT fk_menu_item_id FOREIGN KEY (menu_item_id)
REFERENCES menu(id) ON DELETE CASCADE;

-- Добавляем внешний ключ для столбца (menu_item_id, category) в таблице order_items, который ссылается на таблицу prices
ALTER TABLE order_items
ADD CONSTRAINT fk_prices FOREIGN KEY (menu_item_id, category)
REFERENCES prices(menu_item_id, size) ON DELETE CASCADE;