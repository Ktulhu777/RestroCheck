-- Таблица заказов с комментариями и временем выполнения
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    waiter_id INTEGER REFERENCES waiters(id) ON DELETE SET NULL, -- Кто принял заказ
    created_at TIMESTAMP, -- Время создания заказа
    completed_at TIMESTAMP, -- Время выполнения заказа (заполняется бэкендом)
    actual_completed_at TIMESTAMP, -- Время фактического завершения заказа (заполняется в момент завершения)
    comment TEXT -- Комментарий официанта
);

-- Таблица позиций заказа (какие блюда в заказе)
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id INTEGER NOT NULL REFERENCES menu(id) ON DELETE CASCADE,
    category TEXT NOT NULL, -- Категория блюда (например, маленький, большой)
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price NUMERIC NOT NULL CHECK (price >= 0), -- Цена фиксируется на момент заказа
    UNIQUE (order_id, menu_item_id, category) -- Уникальное сочетание заказа, блюда и категории
);

-- Индексы для улучшения производительности
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_orders_waiter_id ON orders(waiter_id);
