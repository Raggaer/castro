CREATE TABLE `cloaka`.`castro_shop_checkout` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `offer` INT NULL,
  `amount` INT NULL,
  `player` VARCHAR(70) DEFAULT "",
  `given` INT NULL,
PRIMARY KEY (`id`));