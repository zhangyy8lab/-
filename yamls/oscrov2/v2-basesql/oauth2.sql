/*
 Navicat Premium Data Transfer

 Source Server         : 192.168.1.99
 Source Server Type    : MySQL
 Source Server Version : 50720
 Source Host           : 192.168.1.99:3306
 Source Schema         : oauth2

 Target Server Type    : MySQL
 Target Server Version : 50720
 File Encoding         : 65001

 Date: 17/07/2023 10:19:04
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for oauth2_client
-- ----------------------------
DROP TABLE IF EXISTS `oauth2_client`;
CREATE TABLE `oauth2_client` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `client_id` longtext,
  `domain` longtext,
  `secret` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for oauth2_code
-- ----------------------------
DROP TABLE IF EXISTS `oauth2_code`;
CREATE TABLE `oauth2_code` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `expired_at` bigint(20) DEFAULT NULL,
  `code` varchar(128) DEFAULT NULL,
  `client_id` longtext,
  `username` longtext,
  `redirect_uri` longtext,
  `scope` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=835 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for oauth2_users
-- ----------------------------
DROP TABLE IF EXISTS `oauth2_users`;
CREATE TABLE `oauth2_users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `pass` varchar(128) DEFAULT NULL,
  `role` varchar(60) DEFAULT NULL,
  `e_mail` varchar(128) DEFAULT NULL,
  `phone_number` varchar(60) DEFAULT NULL,
  `wallet_id` varchar(256) DEFAULT NULL,
  `last_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=163 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
