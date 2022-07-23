-- Create databases.
CREATE DATABASE example OWNER exampleuser TEMPLATE template0 ENCODING 'UTF-8' LC_COLLATE 'ja_JP.UTF-8' LC_CTYPE 'ja_JP.UTF-8';
CREATE DATABASE example_test OWNER exampleuser TEMPLATE template0 ENCODING 'UTF-8' LC_COLLATE 'ja_JP.UTF-8' LC_CTYPE 'ja_JP.UTF-8';

-- Create tables.
\connect example;
\i /docker-entrypoint-initdb.d/ddl/example.sql;

\connect example_test;
\i /docker-entrypoint-initdb.d/ddl/example.sql;
