CREATE TABLE `castro_shop_checkout` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `offer` VARCHAR(255) NULL,
  `amount` VARCHAR(255) NULL,
  `player` VARCHAR(70) DEFAULT "",
  `given` INT NULL,
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;