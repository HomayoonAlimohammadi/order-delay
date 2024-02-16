CREATE TABLE IF NOT EXISTS vendors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS agents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    vendor_id INT NOT NULL,
    delivery_time INTERVAL NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vendor_id) REFERENCES vendors(id)
);

CREATE TYPE trip_status AS ENUM ('assigned', 'at_vendor', 'picked', 'delivered');

CREATE TABLE IF NOT EXISTS trips (
    id BIGSERIAL PRIMARY KEY,
    status trip_status NOT NULL,
    order_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE TABLE IF NOT EXISTS delay_reports (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,   
    order_id BIGINT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE TABLE IF NOT EXISTS delay_checks (
    id BIGSERIAL PRIMARY KEY,
    agent_id INTEGER NOT NULL,
    report_id BIGINT NOT NULL,
    FOREIGN KEY (agent_id) REFERENCES agents(id),
    FOREIGN KEY (report_id) REFERENCES delay_reports(id)
);

-- Seed initial data
INSERT INTO vendors (name) VALUES 
('Vendor A'),
('Vendor B'),
('Vendor C');

INSERT INTO agents (name) VALUES 
('Agent X'),
('Agent Y'),
('Agent Z');

INSERT INTO orders (vendor_id, delivery_time) VALUES
(1, INTERVAL '30 minutes'),
(1, INTERVAL '45 minutes'),
(1, INTERVAL '50 minutes'),
(2, INTERVAL '30 minutes'),
(2, INTERVAL '45 minutes'),
(2, INTERVAL '50 minutes'),
(3, INTERVAL '45 minutes'),
(3, INTERVAL '50 minutes'),
(1, INTERVAL '60 minutes');

INSERT INTO trips (status, order_id) VALUES
  ('assigned', 1),
  ('at_vendor', 2),
  ('picked', 3),
  ('delivered', 4);
