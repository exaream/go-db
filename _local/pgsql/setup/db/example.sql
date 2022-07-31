-- Create databases.
CREATE DATABASE example_db OWNER exampleuser TEMPLATE template0 ENCODING 'UTF-8' LC_COLLATE 'ja_JP.UTF-8' LC_CTYPE 'ja_JP.UTF-8';
CREATE DATABASE example_db_dbutil_pkg_test OWNER exampleuser TEMPLATE template0 ENCODING 'UTF-8' LC_COLLATE 'ja_JP.UTF-8' LC_CTYPE 'ja_JP.UTF-8';
CREATE DATABASE example_db_example_pkg_test OWNER exampleuser TEMPLATE template0 ENCODING 'UTF-8' LC_COLLATE 'ja_JP.UTF-8' LC_CTYPE 'ja_JP.UTF-8';

-- Create tables.
\connect example_db;
\i /docker-entrypoint-initdb.d/ddl/example.sql;

\connect example_db_dbutil_pkg_test;
\i /docker-entrypoint-initdb.d/ddl/example.sql;

\connect example_db_example_pkg_test;
\i /docker-entrypoint-initdb.d/ddl/example.sql;
