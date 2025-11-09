-- 002_seed_banks.down.sql
DELETE FROM banks WHERE id IN ('vbank', 'abank', 'sbank');
