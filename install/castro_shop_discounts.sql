CREATE TABLE `castro_shop_discounts` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(45) NULL,
  `created_at` INT NULL,
  `valid_till` INT NULL,
  `discount` INT NULL,
  `uses` INT NULL,
  `unlimited` INT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
