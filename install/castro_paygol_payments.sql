CREATE TABLE `castro_paygol_payments` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `transaction_id` VARCHAR(45) NULL,
  `custom` VARCHAR(45) NULL,
  `price` INT NULL,
  `points` INT NULL,
  `created_at` INT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;