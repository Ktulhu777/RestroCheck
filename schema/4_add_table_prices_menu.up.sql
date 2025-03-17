CREATE TABLE prices (
    id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu(id) ON DELETE CASCADE,
    size VARCHAR(50) NOT NULL CHECK (size IN ('Маленькая', 'Средняя', 'Большая')),
    price NUMERIC(10,2) CHECK (price >= 0) NOT NULL,
    UNIQUE (menu_item_id, size)
);

CREATE INDEX idx_prices_menu_item_id ON prices(menu_item_id);
CREATE INDEX idx_prices_size ON prices(size);
