-- Create databases.
CREATE DATABASE example_db OWNER exampleuser TEMPLATE template0 ENCODING 'UTF-8' LC_COLLATE 'ja_JP.UTF-8' LC_CTYPE 'ja_JP.UTF-8';
CREATE DATABASE test_dbutil_db OWNER exampleuser TEMPLATE template0 ENCODING 'UTF-8' LC_COLLATE 'ja_JP.UTF-8' LC_CTYPE 'ja_JP.UTF-8';
CREATE DATABASE test_example_db OWNER exampleuser TEMPLATE template0 ENCODING 'UTF-8' LC_COLLATE 'ja_JP.UTF-8' LC_CTYPE 'ja_JP.UTF-8';

-- Create tables.
\connect example_db;
\i /docker-entrypoint-initdb.d/ddl/example.sql;

\connect test_dbutil_db;
\i /docker-entrypoint-initdb.d/ddl/example.sql;

\connect test_example_db;
\i /docker-entrypoint-initdb.d/ddl/example.sql;
