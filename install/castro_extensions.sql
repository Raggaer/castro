CREATE TABLE `castro_extensions` (
  `uid` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) DEFAULT NULL,
  `id` VARCHAR(45) DEFAULT NULL,
  `author` VARCHAR(45) DEFAULT NULL,
  `type` VARCHAR(45) DEFAULT NULL,
  `version` VARCHAR(45) DEFAULT NULL,
  `description` longtext,
  `installed` BIT NOT NULL DEFAULT 1,
  `created_at` BIGINT(20) NOT NULL,
  `updated_at` BIGINT(20) NOT NULL,
  UNIQUE KEY (`id`),
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `castro_extension_hooks` (
  `uid` INT NOT NULL AUTO_INCREMENT,
  `extension_id` VARCHAR(45) DEFAULT NULL,
  `type` VARCHAR(45) DEFAULT NULL,
  `script` VARCHAR(45) DEFAULT NULL,
  `enabled` BIT NOT NULL DEFAULT 1,
  PRIMARY KEY (`uid`),
  FOREIGN KEY (`extension_id`) REFERENCES `castro_extensions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `castro_extension_pages` (
  `uid` INT NOT NULL AUTO_INCREMENT,
  `extension_id` VARCHAR(45) DEFAULT NULL,
  `enabled` BIT NOT NULL DEFAULT 1,
  PRIMARY KEY (`uid`),
  FOREIGN KEY (`extension_id`) REFERENCES `castro_extensions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `castro_extension_widgets` (
  `uid` INT NOT NULL AUTO_INCREMENT,
  `extension_id` VARCHAR(45) DEFAULT NULL,
  `enabled` BIT NOT NULL DEFAULT 1,
  PRIMARY KEY (`uid`),
  FOREIGN KEY (`extension_id`) REFERENCES `castro_extensions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `castro_extension_templatehooks` (
  `uid` INT NOT NULL AUTO_INCREMENT,
  `extension_id` VARCHAR(45) DEFAULT NULL,
  `type` VARCHAR(45) DEFAULT NULL,
  `template` VARCHAR(45) DEFAULT NULL,
  `enabled` BIT NOT NULL DEFAULT b'1',
  PRIMARY KEY (`uid`),
  FOREIGN KEY (`extension_id`) REFERENCES `castro_extensions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;