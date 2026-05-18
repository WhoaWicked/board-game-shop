CREATE TABLE "roles" (
  "id" int PRIMARY KEY,
  "title" varchar
);

CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "username" varchar UNIQUE,
  "email" varchar UNIQUE,
  "password" varchar,
  "role_id" int,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "oauth" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar,
  "access_token" varchar,
  "refresh_token" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "board_tables" (
  "id" varchar PRIMARY KEY,
  "table_number" varchar UNIQUE,
  "seat_capacity" int,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "games" (
  "id" varchar PRIMARY KEY,
  "title" varchar,
  "description" varchar,
  "status" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "game_images" (
  "id" varchar PRIMARY KEY,
  "filename" varchar
  "url" varchar
  "created_at" timestamp,
  "updated_at" timestamp
)

CREATE TABLE "categories" (
  "id" int PRIMARY KEY,
  "title" varchar UNIQUE
);

CREATE TABLE "games_categories" (
  "id" varchar PRIMARY KEY,
  "game_id" varchar,
  "category_id" int
);

CREATE TABLE "bookings" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar,
  "table_id" varchar,
  "total_players" int,
  "start_time" timestamp,
  "end_time" timestamp,
  "actual_start_time" timestamp,
  "actual_end_time" timestamp,
  "status" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "booking_games" (
  "id" varchar PRIMARY KEY,
  "booking_id" varchar,
  "game_id" varchar
);

CREATE TABLE "booking_rates" (
  "id" int PRIMARY KEY,
  "min_hours" int,
  "max_hours" int,
  "price_per_hour" decimal
);

CREATE TABLE "payments" (
  "id" varchar PRIMARY KEY,
  "booking_id" varchar UNIQUE,
  "total_hours_price" decimal,
  "total_penalty_price" decimal,
  "grand_total" decimal,
  "rate_applied_per_hour" decimal,
  "status" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "oauth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "games_categories" ADD FOREIGN KEY ("game_id") REFERENCES "games" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "games_categories" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "bookings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "bookings" ADD FOREIGN KEY ("table_id") REFERENCES "board_tables" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "booking_games" ADD FOREIGN KEY ("booking_id") REFERENCES "bookings" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "booking_games" ADD FOREIGN KEY ("game_id") REFERENCES "games" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "payments" ADD FOREIGN KEY ("booking_id") REFERENCES "bookings" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "game_images" ADD FOREIGN KEY ("game_id") REFERENCES "games" ("id") DEFERRABLE INITIALLY IMMEDIATE;
