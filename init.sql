CREATE DATABASE IF NOT EXISTS questions_db
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE questions_db;

CREATE TABLE IF NOT EXISTS questions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    template_for_go TEXT NOT NULL,
    template_for_python TEXT NOT NULL,
    parameters JSON NOT NULL
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;