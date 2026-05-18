BEGIN;

-- ==========================================
-- 1. INSERT ROLES (ตารางเดียวที่ล็อกเลข ID สำหรับ Bitwise 3 สิทธิ์)
-- ==========================================
INSERT INTO "roles" ("id", "title") VALUES 
(1, 'Customer'),
(2, 'Staff'),
(4, 'Admin');

-- หมุนเข็มไมล์เฉพาะตาราง roles ให้ไปรอที่เลข 4 (เพราะเป็น SERIAL)
SELECT setval('roles_id_seq', 4, true);


-- ==========================================
-- 2. INSERT USERS (ปล่อยรัน ID ออโต้เป็น U000001, U000002, ...)
-- ==========================================
INSERT INTO "users" ("username", "email", "password", "role_id") VALUES 
('admin001', 'admin@game.com', '$2a$10$4iyHX4jPPbIK5pQJHrlZK.1os0Vh2W/m7njD0QkGs57UDbi90PcKm', 4), -- ระบบรันให้เป็น U000001
('staff001', 'staff@game.com', '$2a$10$OlfwTN20Bv3lnTxbZBh3wu5KLaAKGd.wdIJ2lRdpyjOb6Dfz.zZ.e', 2),   -- ระบบรันให้เป็น U000002
('customer001', 'customer@gmail.com', '$2a$10$4iyHX4jPPbIK5pQJHrlZK.1os0Vh2W/m7njD0QkGs57UDbi90PcKm', 1);     -- ระบบรันให้เป็น U000003


-- ==========================================
-- 3. INSERT BOARD TABLES (ปล่อยรัน ID ออโต้เป็น 1, 2, 3, 4)
-- ==========================================
INSERT INTO "board_tables" ("table_number", "seat_capacity") VALUES 
('T01', 4),    -- id = 1
('T02', 4),    -- id = 2
('T03', 6),    -- id = 3
('T04', 8);    -- id = 4


-- ==========================================
-- 4. INSERT CATEGORIES (ปล่อยรัน ID ออโต้เป็น 1, 2, 3, 4)
-- ==========================================
INSERT INTO "categories" ("title") VALUES 
('Strategy'),  -- id = 1
('Party'),     -- id = 2
('Family'),    -- id = 3
('Bluffing');  -- id = 4


-- ==========================================
-- 5. INSERT GAMES (ปล่อยรัน ID ออโต้เป็น G000001, G000002, ...)
-- ==========================================
INSERT INTO "games" ("title", "description", "status") VALUES 
('Catan', 'Settlers of Catan - ค้าขาย วางแผน สร้างเมือง', 'available'), -- ระบบรันให้เป็น G000001
('Avalon', 'The Resistance: Avalon - เกมบลัฟฟ์จับกลุ่มคนทรยศ', 'available'), -- ระบบรันให้เป็น G000002
('Dixit', 'เกมการ์ดคำใบ้จากภาพวาดสุดแฟนตาซี', 'available'); -- ระบบรันให้เป็น G000003


-- ==========================================
-- 6. INSERT GAME IMAGES (ผูกตาม ID เกมที่ได้จากข้อ 5)
-- ==========================================
INSERT INTO "game_images" ("game_id", "filename", "url") VALUES 
('G000001', 'catan_main.jpg', 'https://images.unsplash.com/photo-1549056572-75914d5d5fd4?q=80&w=764&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
('G000002', 'avalon_main.jpg', 'https://images.unsplash.com/photo-1570303345338-e1f0eddf4946?q=80&w=1071&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
('G000003', 'dixit_main.jpg', 'https://images.unsplash.com/photo-1677188010559-0667a1ed33a0?q=80&w=1113&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D');


-- ==========================================
-- 7. INSERT GAMES_CATEGORIES (จับคู่ ID เกม กับ ID หมวดหมู่)
-- ==========================================
INSERT INTO "games_categories" ("game_id", "category_id") VALUES 
('G000001', 1), -- Catan คู่กับ Strategy (1)
('G000001', 3), -- Catan คู่กับ Family (3)
('G000002', 2), -- Avalon คู่กับ Party (2)
('G000002', 4), -- Avalon คู่กับ Bluffing (4)
('G000003', 2), -- Dixit คู่กับ Party (2)
('G000003', 3); -- Dixit คู่กับ Family (3)


-- ==========================================
-- 8. INSERT BOOKING RATES (ปล่อยรัน ID ออโต้เป็น 1, 2, 3)
-- ==========================================
INSERT INTO "booking_rates" ("min_hours", "max_hours", "price_per_hour") VALUES 
(1, 2, 50.00),  -- id = 1
(3, 5, 40.00),  -- id = 2
(6, 12, 30.00); -- id = 3

-- ==========================================
-- 9. INSERT BOOKINGS (จำลองลูกค้า customer_ton จองโต๊ะ VIP01)
-- ==========================================
-- ระบบรัน ID ออโต้เป็น B000001 โดยอิงตาม Foreign Key:
-- user_id = 'U000003' (customer_ton)
-- table_id = 4 (โต๊ะ VIP01)
INSERT INTO "bookings"
("user_id", "table_id", "total_players", "start_time", "end_time", "status")
VALUES
('U000003', 4, 4, NOW() + INTERVAL '2 hours', NOW() + INTERVAL '5 hours', 'booked');


-- ==========================================
-- 10. INSERT BOOKING_GAMES (จำลองว่าการจอง B000001 มีการยืมเกม Avalon ไปเล่น)
-- ==========================================
-- booking_id = 'B000001'
-- game_id = 'G000002' (Avalon)
INSERT INTO "booking_games" ("booking_id", "game_id") VALUES 
('B000001', 'G000002');

-- ==========================================
-- 11. INSERT PAYMENTS (บิลจ่ายเงิน อิงตามคอลัมน์เวอร์ชันใหม่)
-- ==========================================
-- ผูกกับบิลจอง 'B000001': เล่น 3 ชั่วโมง (เรทชั่วโมงละ 40.00) = 120.00 บาท ยังไม่มีค่าปรับ
INSERT INTO "payments" ("booking_id", "total_hours_price", "total_penalty_price", "grand_total", "rate_applied_per_hour", "status") VALUES 
('B000001', 120.00, 0.00, 120.00, 40.00, 'success');

COMMIT;