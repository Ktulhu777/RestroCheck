-- Таблица категорий блюд (например, "Завтраки", "Обеды")
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE 
);

-- Таблица блюд с привязкой к категории
CREATE TABLE menu (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    photo_url TEXT NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL
);

CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_menu_name ON menu(name);
CREATE INDEX idx_menu_category_id ON menu(category_id);