CREATE TABLE IF NOT EXISTS unique_index_table (
    id INT NOT NULL AUTO_INCREMENT,
    metadata VARCHAR(100),
    value DECIMAL(10, 2),
    occurred_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS non_unique_index_table (
    id INT,
    metadata VARCHAR(100),
    value DECIMAL(10, 2),
    occurred_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX (id)
);
