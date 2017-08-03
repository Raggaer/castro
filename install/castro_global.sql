CREATE TABLE `castro_global` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `key` varchar(75) DEFAULT NULL,
  `value` blob,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;