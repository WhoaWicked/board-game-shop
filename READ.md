migrate -source file://C:/Users/User/Desktop/LearningGO/kawaii-shop/pkg/databases/migrations -database 'postgres://kawaii:123456@localhost:4444/kawaii_db_test?sslmode=disable' -verbose up 

migrate -source file://C:/Users/User/Desktop/board-game-shop/pkg/databases/migrations -database 'postgres://paster:123456@localhost:5555/board_game_db?sslmode=disable' -verbose up 

migrate -source file://C:/Users/User/Desktop/board-game-shop/pkg/databases/migrations -database 'postgres://paster:123456@localhost:5555/board_game_db?sslmode=disable' -verbose down

