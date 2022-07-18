CREATE DATABASE IF NOT EXISTS example CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
CREATE DATABASE IF NOT EXISTS example_test CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
GRANT ALL PRIVILEGES ON example.* TO exampleuser@'%';
GRANT ALL PRIVILEGES ON example_test.* TO exampleuser@'%';
FLUSH PRIVILEGES;

USE example;
SOURCE /docker-entrypoint-initdb.d/ddl/example.sql;
USE example_test;
SOURCE /docker-entrypoint-initdb.d/ddl/example.sql;
