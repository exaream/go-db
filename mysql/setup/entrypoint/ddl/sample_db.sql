DROP TABLE IF EXISTS `sample_db`.`users`;

CREATE TABLE `sample_db`.`users` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `status` int(11) UNSIGNED NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `sample_db`.`users` (`id`, `name`, `email`, `status`, `created_at`, `updated_at`) VALUES
(1, 'Alice', 'sample1@sample.com', 1, '2022-06-13 09:18:17', '2022-06-14 18:03:31'),
(2, 'Bob', 'sample2@sample.com', 1, '2022-06-13 09:18:17', '2022-06-14 20:19:53'),
(3, 'Chris', 'sample3@sample.com', 3, '2022-06-13 09:18:17', '2022-06-14 18:07:52');
