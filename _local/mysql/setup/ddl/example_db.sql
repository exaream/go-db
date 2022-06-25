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
(1, 'Alice', 'example1@example.com', 0, NOW(), NOW()),
(2, 'Billy', 'example2@example.com', 0, NOW(), NOW()),
(3, 'Chris', 'example3@example.com', 0, NOW(), NOW()),
(4, 'Daisy', 'example4@example.com', 0, NOW(), NOW()),
(5, 'Elise', 'example5@example.com', 0, NOW(), NOW());
