CREATE TABLE `castro_fortumo_payments` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `account` VARCHAR(45) NULL,
  `points` INT NULL,
  `price` INT NULL,
  `currency` VARCHAR(45) NULL,
  `sender` VARCHAR(75) NULL,
  `operator` VARCHAR(45) NULL,
  `payment_id` VARCHAR(45) NULL,
  `created_at` INT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
