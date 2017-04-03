CREATE TABLE `castro_paypal_payments` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `payment_id` VARCHAR(45) NULL,
  `payer_id` VARCHAR(45) NULL,
  `custom` VARCHAR(45) NULL,
  `package_name` VARCHAR(45) NULL,
  `state` VARCHAR(45) NULL,
  `created_at` INT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;