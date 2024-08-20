/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80022
 Source Host           : localhost:3306
 Source Schema         : im

 Target Server Type    : MySQL
 Target Server Version : 80022
 File Encoding         : 65001

 Date: 20/12/2020 20:31:24
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for communities
-- ----------------------------
DROP TABLE IF EXISTS `communities`;
CREATE TABLE `communities` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '群名称',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '群图片',
  `owner_id` int NOT NULL DEFAULT '0' COMMENT '群创建人 id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `owner_id_index` (`owner_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='群表';

-- ----------------------------
-- Table structure for community_users
-- ----------------------------
DROP TABLE IF EXISTS `community_users`;
CREATE TABLE `community_users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `community_id` int NOT NULL DEFAULT '0',
  `user_id` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `community_id_index` (`community_id`),
  KEY `user_id_index` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='群组用户';

-- ----------------------------
-- Table structure for contacts
-- ----------------------------
DROP TABLE IF EXISTS `contacts`;
CREATE TABLE `contacts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户 ID',
  `friend_id` int NOT NULL DEFAULT '0' COMMENT '好友 ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id_index` (`user_id`),
  KEY `friend_id_index` (`friend_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户好友表';

-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL DEFAULT '0' COMMENT '发送消息的用户 ID',
  `cmd` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1:群聊, 2:私聊',
  `to_id` int NOT NULL DEFAULT '0' COMMENT '对应的群 ID 或者私聊用户的 ID',
  `media` tinyint(1) NOT NULL DEFAULT '1' COMMENT '消息样式',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '消息的内容',
  `pic` varchar(255) NOT NULL DEFAULT '' COMMENT '图片预览',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '服务的地址',
  `memo` varchar(255) NOT NULL DEFAULT '' COMMENT '简单描述',
  `amount` int NOT NULL DEFAULT '0' COMMENT '和数字相关',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id_index` (`user_id`),
  KEY `to_id_index` (`to_id`),
  KEY `uid_tid_index` (`user_id`,`to_id`),
  KEY `tid_uid_index` (`to_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=97 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='消息表';

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `token` varchar(50) NOT NULL DEFAULT '' COMMENT 'token',
  `avatar` varchar(160) NOT NULL DEFAULT '' COMMENT '用户头像',
  `sex` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0:未知, 1:男, 2:女',
  `nick_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `salt` varchar(10) NOT NULL DEFAULT '' COMMENT '盐值',
  `online` int NOT NULL DEFAULT '0',
  `stat` tinyint(1) NOT NULL DEFAULT '1' COMMENT '账户状态(0:冻结, 1:可用)',
  `mobile` varchar(11) NOT NULL DEFAULT '',
  `password` varchar(40) NOT NULL DEFAULT '',
  `memo` varchar(255) NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `mobile_unique_index` (`mobile`),
  KEY `token_index` (`token`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;
