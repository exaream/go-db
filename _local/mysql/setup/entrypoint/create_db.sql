CREATE DATABASE IF NOT EXISTS example_db CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
GRANT ALL PRIVILEGES ON example_db.* TO exampleuser@'%';
FLUSH PRIVILEGES;
SOURCE /docker-entrypoint-initdb.d/ddl/example_db.sql;
