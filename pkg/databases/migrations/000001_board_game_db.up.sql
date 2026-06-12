BEGIN;

SET TIME ZONE 'Asia/Bangkok';

-- Install uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE sequence
CREATE SEQUENCE users_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE games_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE bookings_id_seq START WITH 1 INCREMENT BY 1;

-- Auto update
CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

-- Create enum
CREATE TYPE "table_status" AS ENUM (
    'available',
    'occupied',
    'maintenance',
    'hidden'
);

CREATE TYPE "game_status" AS ENUM (
    'available',
    'borrowed',
    'damaged',
    'lost',
    'maintenance'
);

CREATE TYPE "booking_status" AS ENUM (
  'booked',
  'active',
  'cancelled',
  'noshow_expired',
  'completed'
);

CREATE TYPE "payment_status" AS ENUM (
  'pending',
  'success',
  'failed'
);

CREATE TABLE "roles" (
  "id" SERIAL PRIMARY KEY,
  "title" VARCHAR NOT NULL UNIQUE
);

CREATE TABLE "users" (
  "id" VARCHAR PRIMARY KEY DEFAULT CONCAT('U', LPAD(NEXTVAL('users_id_seq')::TEXT, 6, '0')),
  "username" VARCHAR UNIQUE NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "role_id" INT NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "oauth" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" VARCHAR NOT NULL,
  "access_token" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "board_tables" (
  "id" SERIAL PRIMARY KEY,
  "table_number" VARCHAR UNIQUE NOT NULL,
  "seat_capacity" INT NOT NULL,
  "status" table_status NOT NULL DEFAULT 'available',
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT "chk_seat_capacity" CHECK ("seat_capacity" > 0)
);

-- ✨ แก้ไขเพิ่มคอลัมน์ในตาราง games
CREATE TABLE "games" (
  "id" VARCHAR PRIMARY KEY DEFAULT CONCAT('G', LPAD(NEXTVAL('games_id_seq')::TEXT, 6, '0')),
  "title" VARCHAR UNIQUE NOT NULL,
  "description" VARCHAR,
  "min_players" INT NOT NULL,          -- 🔥 จำนวนผู้เล่นขั้นต่ำ
  "max_players" INT NOT NULL,          -- 🔥 จำนวนผู้เล่นสูงสุด
  "playing_time" INT NOT NULL,         -- 🔥 ระยะเวลาเล่นเฉลี่ย (นาที)
  "status" game_status NOT NULL DEFAULT 'available',
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now(),
  -- 🔒 ดักข้อมูล: จำนวนผู้เล่นสูงสุดต้องไม่น้อยกว่าขั้นต่ำ และทุกค่าต้องมากกว่า 0
  CONSTRAINT "chk_players_range" CHECK ("max_players" >= "min_players"),
  CONSTRAINT "chk_positive_values" CHECK ("min_players" > 0 AND "playing_time" > 0)
);

CREATE TABLE "game_images" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "game_id" VARCHAR NOT NULL,
  "filename" VARCHAR UNIQUE NOT NULL,
  "url" VARCHAR UNIQUE NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "categories" (
  "id" SERIAL PRIMARY KEY,
  "title" VARCHAR UNIQUE NOT NULL
);

CREATE TABLE "games_categories" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "game_id" VARCHAR NOT NULL,
  "category_id" INT NOT NULL,
  CONSTRAINT "unique_game_category" UNIQUE ("game_id", "category_id")
);

CREATE TABLE "bookings" (
  "id" VARCHAR PRIMARY KEY DEFAULT CONCAT('B', LPAD(NEXTVAL('bookings_id_seq')::TEXT, 6, '0')),
  "user_id" VARCHAR NOT NULL,
  "table_id" INT NOT NULL,
  "total_players" INT NOT NULL,
  "start_time" TIMESTAMP NOT NULL,
  "end_time" TIMESTAMP NOT NULL,
  "actual_start_time" TIMESTAMP,
  "actual_end_time" TIMESTAMP,
  "status" booking_status NOT NULL DEFAULT 'booked',
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT "chk_booking_time" CHECK ("end_time" > "start_time"),
  CONSTRAINT "chk_total_players" CHECK ("total_players" > 0)
);

CREATE TABLE "booking_games" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "booking_id" VARCHAR NOT NULL,
  "game_id" VARCHAR NOT NULL,
  CONSTRAINT "unique_booking_game" UNIQUE ("booking_id", "game_id")
);

CREATE TABLE "booking_rates" (
  "id" SERIAL PRIMARY KEY,
  "min_hours" INT NOT NULL,
  "max_hours" INT NOT NULL,
  "price_per_hour" DECIMAL(10,2) NOT NULL,
  CONSTRAINT chk_booking_rate_hour CHECK (
  min_hours > 0
  AND max_hours >= min_hours )
);

CREATE TABLE "payments" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "booking_id" VARCHAR UNIQUE NOT NULL,
  "total_hours_price" DECIMAL(10,2) NOT NULL,
  "total_penalty_price" DECIMAL(10,2) NOT NULL DEFAULT 0,
  "grand_total" DECIMAL(10,2) NOT NULL,
  "rate_applied_per_hour" DECIMAL(10,2) NOT NULL,
  "status" payment_status NOT NULL DEFAULT 'pending',
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "oauth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "games_categories" ADD FOREIGN KEY ("game_id") REFERENCES "games" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "games_categories" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "bookings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "bookings" ADD FOREIGN KEY ("table_id") REFERENCES "board_tables" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "booking_games" ADD FOREIGN KEY ("booking_id") REFERENCES "bookings" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "booking_games" ADD FOREIGN KEY ("game_id") REFERENCES "games" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "payments" ADD FOREIGN KEY ("booking_id") REFERENCES "bookings" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "game_images" ADD FOREIGN KEY ("game_id") REFERENCES "games" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

CREATE TRIGGER set_updated_at_timestamp_users_table BEFORE UPDATE ON "users" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_oauth_table BEFORE UPDATE ON "oauth" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_board_tables_table BEFORE UPDATE ON "board_tables" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_games_table BEFORE UPDATE ON "games" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_bookings_table BEFORE UPDATE ON "bookings" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_payments_table BEFORE UPDATE ON "payments" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_game_images_table BEFORE UPDATE ON "game_images" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_table_id ON bookings(table_id);
CREATE INDEX idx_booking_games_booking_id ON booking_games(booking_id);

CREATE INDEX idx_games_status ON games(status);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_payments_status ON payments(status);

CREATE INDEX idx_bookings_start_time ON bookings(start_time);

CREATE INDEX idx_game_images_game_id ON game_images(game_id);

COMMIT;