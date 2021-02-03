CREATE TABLE IF NOT EXISTS `tasks` (
    `id` SERIAL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT DEFAULT NULL,
    `state` VARCHAR(128) NOT NULL,
    `priority` VARCHAR(128) NOT NULL,
    `value` VARCHAR(128) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY (`name`, `description`)
)