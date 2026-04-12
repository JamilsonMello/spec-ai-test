CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(60) NOT NULL,
    description VARCHAR(600) NOT NULL,
    price BIGINT NOT NULL,
    seller_id VARCHAR(36) NOT NULL,
    stock INT NOT NULL,
    category VARCHAR(50) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    community_id VARCHAR(36),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
