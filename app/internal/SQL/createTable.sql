-- Создание таблицы заказов
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,              
    name VARCHAR(255) NOT NULL,    
    count INT NOT NULL,            
    status VARCHAR(50) NOT NULL,      
    created_at TIMESTAMP NOT NULL,    
    price NUMERIC(10, 2) NOT NULL     
);