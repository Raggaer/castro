CREATE TABLE `castro_map` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(75) DEFAULT NULL,
  `data` blob,
  `created_at` TIMESTAMP DEFAULT NULL,
  `updated_at` TIMESTAMP DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;