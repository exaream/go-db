CREATE DATABASE IF NOT EXISTS sample_db CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
GRANT ALL PRIVILEGES ON sample_db.* TO opsuser@'%';
FLUSH PRIVILEGES;
SOURCE /docker-entrypoint-initdb.d/ddl/sample_db.sql;
