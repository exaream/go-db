-- Create databases.
CREATE DATABASE IF NOT EXISTS example_db CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
CREATE DATABASE IF NOT EXISTS test_dbutil_db CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
CREATE DATABASE IF NOT EXISTS test_example_db CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
GRANT ALL PRIVILEGES ON example_db.* TO exampleuser@'%';
GRANT ALL PRIVILEGES ON test_dbutil_db.* TO exampleuser@'%';
GRANT ALL PRIVILEGES ON test_example_db.* TO exampleuser@'%';
FLUSH PRIVILEGES;

-- Create tables.
USE example_db;
SOURCE /docker-entrypoint-initdb.d/ddl/example.sql;
USE test_dbutil_db;
SOURCE /docker-entrypoint-initdb.d/ddl/example.sql;
USE test_example_db;
SOURCE /docker-entrypoint-initdb.d/ddl/example.sql;
