CREATE TABLE `castro_forum_message` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `post_id` INT NULL,
  `author` INT NULL,
  `message` MEDIUMTEXT NULL,
  `created_at` INT NULL,
  `updated_at` INT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;