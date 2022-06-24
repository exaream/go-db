DROP TABLE IF EXISTS `example_db`.`users`;

CREATE TABLE `example_db`.`users` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `status` int(11) UNSIGNED NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `example_db`.`users` (`id`, `name`, `email`, `status`, `created_at`, `updated_at`) VALUES
(1, 'Alice', 'example1@example.com', 0, '1885-09-01 00:00:00', '1885-09-01 00:00:00'),
(2, 'Billy', 'example2@example.com', 0, '1885-09-01 00:00:00', '1885-09-01 00:00:00'),
(3, 'Chris', 'example3@example.com', 0, '1885-09-01 00:00:00', '1885-09-01 00:00:00');
-- FYI: Doc Brown wrote a letter to Marty on September 1st, 1885 in the movie "Back to the Future 3".
