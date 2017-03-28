CREATE TABLE `castro_forum_post` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `category_id` INT NULL,
  `title` VARCHAR(45) NULL,
  `author` INT NULL,
  `created_at` INT NULL,
  `updated_at` INT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
