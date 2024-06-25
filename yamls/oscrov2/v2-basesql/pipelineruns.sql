/*
 Navicat Premium Data Transfer

 Source Server         : new-hk-mysql
 Source Server Type    : MySQL
 Source Server Version : 50720
 Source Host           : 192.168.2.253:3306
 Source Schema         : pipelineruns

 Target Server Type    : MySQL
 Target Server Version : 50720
 File Encoding         : 65001

 Date: 17/07/2023 10:23:21
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for pipelineruns
-- ----------------------------
DROP TABLE IF EXISTS `pipelineruns`;
CREATE TABLE `pipelineruns` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `namespace` varchar(128) DEFAULT NULL,
  `reruns` varchar(128) DEFAULT NULL,
  `data` longblob,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=897 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for taskruns
-- ----------------------------
DROP TABLE IF EXISTS `taskruns`;
CREATE TABLE `taskruns` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `namespace` varchar(128) DEFAULT NULL,
  `pipeline_run` varchar(64) DEFAULT NULL,
  `data` longblob,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2807 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;
