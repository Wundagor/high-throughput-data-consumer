-- goapp/create_table.sql

CREATE TABLE IF NOT EXISTS destination_data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL
);
