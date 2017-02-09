CREATE TABLE `castro_accounts` (
  `id` INT(11) NOT NULL,
  `account_id` INT(11) NOT NULL,
  `points` INT(11) DEFAULT 0,
  `admin` TINYINT(1) DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;