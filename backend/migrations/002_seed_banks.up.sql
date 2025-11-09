-- 002_seed_banks.up.sql
-- Seed initial banks data

INSERT INTO banks (id, name, api_base_url, deposit_rate) VALUES
('vbank', 'Virtual Bank', 'https://vbank.open.bankingapi.ru', 8.0),
('abank', 'Awesome Bank', 'https://abank.open.bankingapi.ru', 7.5),
('sbank', 'Smart Bank', 'https://sbank.open.bankingapi.ru', 9.0)
ON CONFLICT (id) DO NOTHING;
