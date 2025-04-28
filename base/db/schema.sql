-- Connect to the database
\c argodb;

-- Create the servers table
CREATE TABLE servers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    ip VARCHAR(15) NOT NULL,
    contact VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the records table
CREATE TABLE records (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    domain VARCHAR(255) NOT NULL,
    ip VARCHAR(15) NOT NULL,
    server INTEGER NOT NULL,
    FOREIGN KEY (server) REFERENCES servers(id) ON DELETE CASCADE
);
