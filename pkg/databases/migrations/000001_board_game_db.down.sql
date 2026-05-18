BEGIN;

-- Drop triggers
DROP TRIGGER IF EXISTS set_updated_at_timestamp_users_table ON users;
DROP TRIGGER IF EXISTS set_updated_at_timestamp_oauth_table ON oauth;
DROP TRIGGER IF EXISTS set_updated_at_timestamp_board_tables_table ON board_tables;
DROP TRIGGER IF EXISTS set_updated_at_timestamp_games_table ON games;
DROP TRIGGER IF EXISTS set_updated_at_timestamp_bookings_table ON bookings;
DROP TRIGGER IF EXISTS set_updated_at_timestamp_payments_table ON payments;
DROP TRIGGER IF EXISTS set_updated_at_timestamp_game_images_table ON game_images;

-- Drop indexes
DROP INDEX IF EXISTS idx_bookings_user_id;
DROP INDEX IF EXISTS idx_bookings_table_id;
DROP INDEX IF EXISTS idx_booking_games_booking_id;
DROP INDEX IF EXISTS idx_games_status;
DROP INDEX IF EXISTS idx_bookings_status;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_bookings_start_time;
DROP INDEX IF EXISTS idx_game_images_game_id;

-- Drop tables
DROP TABLE IF EXISTS penalties;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS booking_rates;
DROP TABLE IF EXISTS booking_games;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS games_categories;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS game_images;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS board_tables;
DROP TABLE IF EXISTS oauth;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;

-- Drop enums
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS booking_status;
DROP TYPE IF EXISTS game_status;

-- Drop function
DROP FUNCTION IF EXISTS set_updated_at_column;

-- Drop sequences
DROP SEQUENCE IF EXISTS users_id_seq;
DROP SEQUENCE IF EXISTS games_id_seq;
DROP SEQUENCE IF EXISTS bookings_id_seq;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;