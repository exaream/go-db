-- Create databases.
CREATE DATABASE IF NOT EXISTS example_db CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
CREATE DATABASE IF NOT EXISTS example_db_dbutil_pkg_test CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
CREATE DATABASE IF NOT EXISTS example_db_example_pkg_test CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
GRANT ALL PRIVILEGES ON example_db.* TO exampleuser@'%';
GRANT ALL PRIVILEGES ON example_db_dbutil_pkg_test.* TO exampleuser@'%';
GRANT ALL PRIVILEGES ON example_db_example_pkg_test.* TO exampleuser@'%';
FLUSH PRIVILEGES;

-- Create tables.
USE example_db;
SOURCE /docker-entrypoint-initdb.d/ddl/example.sql;
USE example_db_dbutil_pkg_test;
SOURCE /docker-entrypoint-initdb.d/ddl/example.sql;
USE example_db_example_pkg_test;
SOURCE /docker-entrypoint-initdb.d/ddl/example.sql;
