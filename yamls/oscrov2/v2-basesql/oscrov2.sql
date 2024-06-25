/*
 Navicat Premium Data Transfer

 Source Server         : 192.168.1.99
 Source Server Type    : MySQL
 Source Server Version : 50720
 Source Host           : 192.168.1.99:3306
 Source Schema         : oscro

 Target Server Type    : MySQL
 Target Server Version : 50720
 File Encoding         : 65001

 Date: 17/07/2023 10:17:06
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app
-- ----------------------------
DROP TABLE IF EXISTS `app`;
CREATE TABLE `app` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(40) DEFAULT NULL,
  `request_id` longtext,
  `namespace_id` bigint(20) unsigned DEFAULT NULL,
  `app_base_id` bigint(20) unsigned DEFAULT NULL,
  `create_by` longtext,
  `edit_by` longtext,
  `create_at` datetime(3) DEFAULT NULL,
  `edit_at` datetime(3) DEFAULT NULL,
  `status` varchar(191) DEFAULT 'stop',
  `desc` longtext,
  `modify` varchar(191) DEFAULT 'false',
  `import_image` varchar(191) DEFAULT 'false',
  `app_kind` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_namespace_id` (`name`,`namespace_id`),
  KEY `fk_app_namespace` (`namespace_id`),
  KEY `fk_app_app_base` (`app_base_id`),
  CONSTRAINT `fk_app_app_base` FOREIGN KEY (`app_base_id`) REFERENCES `app_base` (`id`),
  CONSTRAINT `fk_app_namespace` FOREIGN KEY (`namespace_id`) REFERENCES `namespace` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=85 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for app_base
-- ----------------------------
DROP TABLE IF EXISTS `app_base`;
CREATE TABLE `app_base` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `space` bigint(20) unsigned DEFAULT '0',
  `create_by` longtext,
  `edit_by` longtext,
  `request_id` longtext,
  `create_at` datetime(3) DEFAULT NULL,
  `edit_at` datetime(3) DEFAULT NULL,
  `desc` longtext,
  `app_kind` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_space` (`name`,`space`) USING BTREE,
  KEY `idx_app_base_space` (`space`)
) ENGINE=InnoDB AUTO_INCREMENT=111 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for app_base_yaml
-- ----------------------------
DROP TABLE IF EXISTS `app_base_yaml`;
CREATE TABLE `app_base_yaml` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_base_id` bigint(20) unsigned DEFAULT NULL,
  `name` varchar(64) NOT NULL,
  `kind` varchar(64) NOT NULL,
  `create_by` longtext,
  `edit_by` longtext,
  `lock` bigint(20) unsigned DEFAULT '0',
  `content` longtext,
  `request_id` longtext,
  `create_at` datetime(3) DEFAULT NULL,
  `edit_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_app_base_yaml_app_base_id` (`app_base_id`),
  CONSTRAINT `fk_app_base_yaml_app_base` FOREIGN KEY (`app_base_id`) REFERENCES `app_base` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=271 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for app_occupy_port
-- ----------------------------
DROP TABLE IF EXISTS `app_occupy_port`;
CREATE TABLE `app_occupy_port` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned DEFAULT NULL,
  `port_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_app_occupy_port_app` (`app_id`),
  KEY `fk_app_occupy_port_pool_port_id` (`port_id`) USING BTREE,
  CONSTRAINT `fk_app_occupy_port_app` FOREIGN KEY (`app_id`) REFERENCES `app` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_app_occupy_port_pool_port_id` FOREIGN KEY (`port_id`) REFERENCES `pool_port` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=61 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for app_yaml
-- ----------------------------
DROP TABLE IF EXISTS `app_yaml`;
CREATE TABLE `app_yaml` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) DEFAULT NULL,
  `content` longtext,
  `running_content` longtext,
  `kind` varchar(32) NOT NULL,
  `app_id` bigint(20) unsigned DEFAULT NULL,
  `create_at` datetime(3) DEFAULT NULL,
  `edit_at` datetime(3) DEFAULT NULL,
  `request_id` longtext,
  `modify` varchar(191) DEFAULT 'false',
  PRIMARY KEY (`id`),
  KEY `fk_app_yaml_app` (`app_id`),
  CONSTRAINT `fk_app_yaml_app` FOREIGN KEY (`app_id`) REFERENCES `app` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=301 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for audit_log
-- ----------------------------
DROP TABLE IF EXISTS `audit_log`;
CREATE TABLE `audit_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `request_id` varchar(36) DEFAULT NULL,
  `workspace` varchar(32) DEFAULT NULL,
  `namespace` varchar(32) DEFAULT NULL,
  `server` varchar(32) DEFAULT NULL,
  `url` longtext,
  `method` varchar(32) DEFAULT NULL,
  `result` varchar(32) DEFAULT NULL,
  `message` longtext,
  `error` longtext,
  `created_at` bigint(20) DEFAULT NULL,
  `create_by` varchar(191) DEFAULT NULL,
  `is_nsq` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_audit_log_request_id` (`request_id`),
  KEY `idx_audit_log_workspace` (`workspace`),
  KEY `idx_audit_log_namespace` (`namespace`),
  KEY `idx_audit_log_server` (`server`),
  KEY `idx_audit_log_create_by` (`create_by`)
) ENGINE=InnoDB AUTO_INCREMENT=17643 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for base_yaml
-- ----------------------------
DROP TABLE IF EXISTS `base_yaml`;
CREATE TABLE `base_yaml` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `kind` varchar(64) NOT NULL,
  `content` longtext,
  `create_by` varchar(16) DEFAULT NULL,
  `edit_at` datetime(3) DEFAULT NULL,
  `edit_by` longtext,
  `space` bigint(20) unsigned DEFAULT '0',
  `request_id` varchar(191) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_kind` (`name`,`kind`),
  KEY `idx_base_yaml_space` (`space`),
  KEY `idx_base_yaml_request_id` (`request_id`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for cluster_node
-- ----------------------------
DROP TABLE IF EXISTS `cluster_node`;
CREATE TABLE `cluster_node` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `server` varchar(64) NOT NULL,
  `request_id` varchar(64) DEFAULT NULL,
  `status` varchar(16) DEFAULT 'creating',
  `harbor` bigint(20) unsigned DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_server` (`name`,`server`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for cluster_node_quota
-- ----------------------------
DROP TABLE IF EXISTS `cluster_node_quota`;
CREATE TABLE `cluster_node_quota` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `cluster_node_id` bigint(20) unsigned DEFAULT NULL,
  `cpu_occupy` varchar(16) DEFAULT NULL,
  `cpu_set` varchar(16) DEFAULT NULL,
  `mem_occupy` varchar(16) DEFAULT NULL,
  `mem_set` varchar(16) DEFAULT NULL,
  `storage_occupy` varchar(16) DEFAULT NULL,
  `storage_set` varchar(16) DEFAULT '20Ti',
  `rsa_key` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_cluster_node_quota_cluster_node_id` (`cluster_node_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for cluster_node_security
-- ----------------------------
DROP TABLE IF EXISTS `cluster_node_security`;
CREATE TABLE `cluster_node_security` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `cluster_node_id` bigint(20) unsigned DEFAULT NULL,
  `token` varchar(32) DEFAULT NULL,
  `cert` longtext,
  `key` longtext,
  `method` varchar(32) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cluster_node_security_cluster_node_id` (`cluster_node_id`),
  CONSTRAINT `fk_cluster_node_security_cluster_node` FOREIGN KEY (`cluster_node_id`) REFERENCES `cluster_node` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for dev_scene
-- ----------------------------
DROP TABLE IF EXISTS `dev_scene`;
CREATE TABLE `dev_scene` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `namespace_id` bigint(20) unsigned DEFAULT NULL,
  `pod_name` varchar(32) DEFAULT NULL,
  `service_ip` varchar(32) DEFAULT NULL,
  `service_port` varchar(32) DEFAULT NULL,
  `ssh_port` varchar(32) DEFAULT NULL,
  `status` varchar(32) DEFAULT NULL,
  `language` varchar(32) DEFAULT NULL,
  `repo_url` varchar(120) DEFAULT NULL,
  `container_id` varchar(32) DEFAULT NULL,
  `node_name` varchar(32) DEFAULT NULL,
  `desc` varchar(128) DEFAULT NULL,
  `create_by` varchar(191) DEFAULT NULL,
  `create_at` datetime(3) DEFAULT NULL,
  `request_id` longtext,
  `image_tag` longtext,
  `suffix_name` varchar(8) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_dev_scene_workspace_id` (`workspace_id`),
  KEY `idx_dev_scene_create_by` (`create_by`),
  KEY `fk_dev_scene_namespace` (`namespace_id`),
  CONSTRAINT `fk_dev_scene_namespace` FOREIGN KEY (`namespace_id`) REFERENCES `namespace` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=147 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for image_sync
-- ----------------------------
DROP TABLE IF EXISTS `image_sync`;
CREATE TABLE `image_sync` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `image` varchar(256) DEFAULT NULL,
  `tag` varchar(256) DEFAULT NULL,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `namespace_id` bigint(20) unsigned DEFAULT NULL,
  `app_base_id` bigint(20) unsigned DEFAULT NULL,
  `request_id` varchar(64) DEFAULT NULL,
  `create_by` varchar(64) DEFAULT NULL,
  `status` varchar(64) DEFAULT NULL,
  `status_text` longtext,
  `exec_id` bigint(20) DEFAULT NULL,
  `sync_type` longtext,
  `policy_id` bigint(20) DEFAULT NULL,
  `app_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for namespace
-- ----------------------------
DROP TABLE IF EXISTS `namespace`;
CREATE TABLE `namespace` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `request_id` longtext,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `cluster_node_id` bigint(20) unsigned DEFAULT NULL,
  `status` varchar(191) DEFAULT 'creating',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `name_workspace_id_cluster_node_id` (`name`),
  KEY `fk_namespace_cluster_node` (`cluster_node_id`),
  KEY `fk_namespace_workspace` (`workspace_id`),
  CONSTRAINT `fk_namespace_cluster_node` FOREIGN KEY (`cluster_node_id`) REFERENCES `cluster_node` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `fk_namespace_workspace` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=167 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for namespace_app_constraint
-- ----------------------------
DROP TABLE IF EXISTS `namespace_app_constraint`;
CREATE TABLE `namespace_app_constraint` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `namespace_id` bigint(20) unsigned DEFAULT NULL,
  `app_set` bigint(20) unsigned DEFAULT '5',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=83 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for namespace_constraint
-- ----------------------------
DROP TABLE IF EXISTS `namespace_constraint`;
CREATE TABLE `namespace_constraint` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `namespace_id` bigint(20) unsigned DEFAULT NULL,
  `quota_cpu` varchar(32) DEFAULT '2',
  `quota_mem` varchar(32) DEFAULT '2Gi',
  `quota_storage` varchar(32) DEFAULT '2Gi',
  `limits_pod_cpu` varchar(32) DEFAULT '8',
  `limits_pod_mem` varchar(32) DEFAULT '8Gi',
  `limits_pvc` varchar(32) DEFAULT '200Gi',
  `limits_container_cpu` varchar(32) DEFAULT '1',
  `limits_container_mem` varchar(32) DEFAULT '200Mi',
  `limits_container_cpu_default_request` varchar(32) DEFAULT '100m',
  `limits_container_mem_default_request` varchar(32) DEFAULT '100Mi',
  PRIMARY KEY (`id`),
  KEY `fk_namespace_namespace_constraint` (`namespace_id`),
  CONSTRAINT `fk_namespace_constraint_namespace` FOREIGN KEY (`namespace_id`) REFERENCES `namespace` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_namespace_namespace_constraint` FOREIGN KEY (`namespace_id`) REFERENCES `namespace` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=157 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for namespace_user
-- ----------------------------
DROP TABLE IF EXISTS `namespace_user`;
CREATE TABLE `namespace_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `namespace_id` bigint(20) unsigned DEFAULT NULL,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `user_name` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_namespace_user_namespace_id` (`namespace_id`),
  KEY `idx_namespace_user_workspace_id` (`workspace_id`),
  KEY `idx_namespace_user_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=241 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for pool_address
-- ----------------------------
DROP TABLE IF EXISTS `pool_address`;
CREATE TABLE `pool_address` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `cluster_node_id` bigint(20) unsigned DEFAULT NULL,
  `inside` varchar(32) DEFAULT 'true',
  `extranet` varchar(32) DEFAULT 'false',
  `traefik` varchar(32) DEFAULT 'false',
  `proxy_ip` longtext,
  `distribution` varchar(32) DEFAULT 'false',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name_cluster_node_id` (`name`,`cluster_node_id`) USING BTREE,
  KEY `fk_pool_address_cluster_node` (`cluster_node_id`) USING BTREE,
  CONSTRAINT `fk_pool_address_cluster_node` FOREIGN KEY (`cluster_node_id`) REFERENCES `cluster_node` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for pool_domain
-- ----------------------------
DROP TABLE IF EXISTS `pool_domain`;
CREATE TABLE `pool_domain` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `schema` varchar(32) DEFAULT 'http',
  `domain` varchar(128) DEFAULT NULL,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `pool_address_id` bigint(20) unsigned DEFAULT NULL,
  `app_id` bigint(20) unsigned DEFAULT NULL,
  `service_name` varchar(64) DEFAULT NULL,
  `crt` longtext,
  `key` longtext,
  `cluster_node_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `domain_pool_address_id` (`domain`,`pool_address_id`) USING BTREE,
  KEY `fk_pool_domain_pool_address` (`pool_address_id`) USING BTREE,
  KEY `fk_pool_domain_workspace` (`workspace_id`) USING BTREE,
  KEY `fk_pool_domain_cluster_node` (`cluster_node_id`) USING BTREE,
  CONSTRAINT `fk_pool_domain_cluster_node` FOREIGN KEY (`cluster_node_id`) REFERENCES `cluster_node` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `fk_pool_domain_pool_address` FOREIGN KEY (`pool_address_id`) REFERENCES `pool_address` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `fk_pool_domain_workspace` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for pool_harbor
-- ----------------------------
DROP TABLE IF EXISTS `pool_harbor`;
CREATE TABLE `pool_harbor` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `schema` varchar(32) NOT NULL,
  `domain` varchar(32) NOT NULL,
  `address` varchar(32) NOT NULL,
  `port` varchar(8) DEFAULT '5000',
  `cluster_node_id` bigint(20) unsigned DEFAULT NULL,
  `register_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name_cluster_node_id` (`name`,`cluster_node_id`) USING BTREE,
  KEY `fk_pool_harbor_cluster_node` (`cluster_node_id`) USING BTREE,
  CONSTRAINT `fk_pool_harbor_cluster_node` FOREIGN KEY (`cluster_node_id`) REFERENCES `cluster_node` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for pool_harbor_security
-- ----------------------------
DROP TABLE IF EXISTS `pool_harbor_security`;
CREATE TABLE `pool_harbor_security` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `pool_harbor_id` bigint(20) unsigned DEFAULT NULL,
  `username` varchar(32) NOT NULL,
  `password` varchar(32) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_pool_harbor_security_pool_harbor_id` (`pool_harbor_id`) USING BTREE,
  CONSTRAINT `fk_pool_harbor_security_pool_harbor` FOREIGN KEY (`pool_harbor_id`) REFERENCES `pool_harbor` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for pool_md_harbor
-- ----------------------------
DROP TABLE IF EXISTS `pool_md_harbor`;
CREATE TABLE `pool_md_harbor` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `schema` varchar(32) NOT NULL,
  `domain` varchar(32) NOT NULL,
  `address` varchar(32) NOT NULL,
  `port` varchar(8) DEFAULT '5000',
  `username` varchar(32) NOT NULL,
  `password` varchar(32) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for pool_port
-- ----------------------------
DROP TABLE IF EXISTS `pool_port`;
CREATE TABLE `pool_port` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `cluster_node_id` bigint(20) unsigned DEFAULT NULL,
  `pool_address_id` bigint(20) unsigned DEFAULT NULL,
  `status` varchar(32) DEFAULT 'free',
  `app_id` bigint(20) unsigned DEFAULT NULL,
  `app_name` longtext,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name_pool_address_id` (`name`,`pool_address_id`) USING BTREE,
  KEY `fk_pool_port_workspace` (`workspace_id`) USING BTREE,
  KEY `fk_pool_port_pool_address` (`pool_address_id`) USING BTREE,
  CONSTRAINT `fk_pool_port_pool_address` FOREIGN KEY (`pool_address_id`) REFERENCES `pool_address` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `fk_pool_port_workspace` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=115 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `role_kind` bigint(20) DEFAULT NULL,
  `title` varchar(32) NOT NULL,
  `as_name` varchar(32) NOT NULL,
  `role_desc` varchar(300) NOT NULL,
  `en_desc` varchar(300) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `request_id` longtext,
  `role_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_user_role` (`role_id`),
  CONSTRAINT `fk_user_role` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=97 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workspace
-- ----------------------------
DROP TABLE IF EXISTS `workspace`;
CREATE TABLE `workspace` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `status` varchar(32) DEFAULT 'creating',
  `request_id` longtext,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workspace_app_constraint
-- ----------------------------
DROP TABLE IF EXISTS `workspace_app_constraint`;
CREATE TABLE `workspace_app_constraint` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `set` bigint(20) unsigned DEFAULT '20',
  `occupy` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_workspace_workspace_app_constraint` (`workspace_id`),
  CONSTRAINT `fk_workspace_workspace_app_constraint` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=97 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workspace_calculate_constraint
-- ----------------------------
DROP TABLE IF EXISTS `workspace_calculate_constraint`;
CREATE TABLE `workspace_calculate_constraint` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `cpu_set` varchar(32) DEFAULT '8',
  `cpu_occupy` varchar(32) DEFAULT NULL,
  `mem_set` varchar(32) DEFAULT '16Gi',
  `mem_occupy` varchar(32) DEFAULT NULL,
  `storage_set` varchar(32) DEFAULT '100Gi',
  `storage_occupy` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_workspace_workspace_calculate` (`workspace_id`),
  CONSTRAINT `fk_workspace_workspace_calculate` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=95 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workspace_clusternode_constraint
-- ----------------------------
DROP TABLE IF EXISTS `workspace_clusternode_constraint`;
CREATE TABLE `workspace_clusternode_constraint` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `cluster_node_id` bigint(20) unsigned DEFAULT NULL,
  `cpu_set` varchar(32) DEFAULT '2',
  `cpu_occupy` varchar(32) DEFAULT '0',
  `mem_set` varchar(32) DEFAULT '2Gi',
  `mem_occupy` varchar(32) DEFAULT '0',
  `storage_set` varchar(32) DEFAULT '10Gi',
  `storage_occupy` varchar(32) DEFAULT '0',
  `status` varchar(191) DEFAULT 'creating',
  `request_id` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `workspace_id_cluster_node_id` (`workspace_id`,`cluster_node_id`),
  KEY `fk_workspace_clusterNode_constraint_cluster_node` (`cluster_node_id`),
  CONSTRAINT `fk_workspace_clusterNode_constraint_cluster_node` FOREIGN KEY (`cluster_node_id`) REFERENCES `cluster_node` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_workspace_clusterNode_constraint_workspace` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=191 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workspace_dev_scene_constraint
-- ----------------------------
DROP TABLE IF EXISTS `workspace_dev_scene_constraint`;
CREATE TABLE `workspace_dev_scene_constraint` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `set` bigint(20) unsigned DEFAULT '20',
  `occupy` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_workspace_workspace_dev_scene_constraint` (`workspace_id`),
  CONSTRAINT `fk_workspace_workspace_dev_scene_constraint` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=97 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workspace_git
-- ----------------------------
DROP TABLE IF EXISTS `workspace_git`;
CREATE TABLE `workspace_git` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint(20) unsigned NOT NULL,
  `url` longtext NOT NULL,
  `create_by` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workspace_ipaddress
-- ----------------------------
DROP TABLE IF EXISTS `workspace_ipaddress`;
CREATE TABLE `workspace_ipaddress` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `pool_address_id` bigint(20) unsigned DEFAULT NULL,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_workspace_ipaddress_pool_address` (`pool_address_id`) USING BTREE,
  KEY `fk_workspace_ipaddress_workspace` (`workspace_id`) USING BTREE,
  CONSTRAINT `fk_workspace_ipaddress_pool_address` FOREIGN KEY (`pool_address_id`) REFERENCES `pool_address` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_workspace_ipaddress_workspace` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for workspace_pool_port_constraint
-- ----------------------------
DROP TABLE IF EXISTS `workspace_pool_port_constraint`;
CREATE TABLE `workspace_pool_port_constraint` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `set` bigint(20) unsigned DEFAULT '20',
  `occupy` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_workspace_workspace_pool_port_constraint` (`workspace_id`),
  CONSTRAINT `fk_workspace_workspace_pool_port_constraint` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=97 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workspace_user
-- ----------------------------
DROP TABLE IF EXISTS `workspace_user`;
CREATE TABLE `workspace_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `workspace_id` bigint(20) unsigned DEFAULT NULL,
  `user_name` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_workspace_user_user_id` (`user_id`),
  KEY `idx_workspace_user_workspace_id` (`workspace_id`),
  CONSTRAINT `fk_workspace_user_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_workspace_user_workspace` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=95 DEFAULT CHARSET=utf8;


INSERT INTO `role`(`id`, `name`, `role_kind`, `title`, `as_name`, `role_desc`, `en_desc`) VALUES (1, 'SuperAdmin', 1, '超级管理员', 'admin', '资源/工作空间管理', 'Resource/workspace management');
INSERT INTO `role`(`id`, `name`, `role_kind`, `title`, `as_name`, `role_desc`, `en_desc`) VALUES (3, 'WorkspaceAdmin', 2, '管理员', 'admin', '具有工作空间下所有权限', 'Have all permissions in the workspace');
INSERT INTO `role`(`id`, `name`, `role_kind`, `title`, `as_name`, `role_desc`, `en_desc`) VALUES (5, 'Visitor', 3, '游客用户', 'ops', '游客对工作空间下仅查看权限', 'Tourists only have permission to view the workspace.');
INSERT INTO `role`(`id`, `name`, `role_kind`, `title`, `as_name`, `role_desc`, `en_desc`) VALUES (7, 'Unauthorized', 4, '未授权用户', 'None', '未绑定任何工作空间', 'No workspace is bound.');
INSERT INTO `role`(`id`, `name`, `role_kind`, `title`, `as_name`, `role_desc`, `en_desc`) VALUES (9, 'OpsRole', 5, '运维工程师', 'ops', '具备所有名称空间下的权限', 'Have permissions under all namespaces');
INSERT INTO `role`(`id`, `name`, `role_kind`, `title`, `as_name`, `role_desc`, `en_desc`) VALUES (11, 'TestRole', 6, '测试工程师', 'test', '具备所属名称空间下的权限，不具备冒烟环境权限', 'Have the permission under the namespace, and do not have the permission of the smoking environment.');
INSERT INTO `role`(`id`, `name`, `role_kind`, `title`, `as_name`, `role_desc`, `en_desc`) VALUES (13, 'DevelopRole', 7, '开发工程师', 'dev', '具备所属名称空间下的权限，不具备容器管理删除及修改权限以及应用管理运行APP的权限', 'It has the permissions under the namespace, the permissions of container management deletion and modification, and the permissions of application management to run the APP.');



INSERT INTO `base_yaml`(`id`, `name`, `kind`, `content`, `create_by`, `edit_at`, `edit_by`, `space`, `request_id`) VALUES (1, 'application', 'Application', '{\"kind\": \"Application\", \"spec\": {\"selector\": {\"matchLabels\": {\"app.oscro.io/name\": \"Value\"}}, \"addOwnerRef\": true, \"componentKinds\": []}, \"metadata\": {\"name\": \"Value\", \"namespace\": \"Value\", \"finalizers\": [\"foregroundDeletion\"]}, \"apiVersion\": \"app.oscro.io/v1beta1\"}', NULL, '2022-11-30 01:58:38.000', NULL, 0, NULL);
INSERT INTO `base_yaml`(`id`, `name`, `kind`, `content`, `create_by`, `edit_at`, `edit_by`, `space`, `request_id`) VALUES (23, '无状态服务', 'Deployment', '{\"kind\": \"Deployment\", \"spec\": {\"replicas\": 1, \"selector\": {\"matchLabels\": {\"app.oscro.io/name\": \"VALUE\"}}, \"template\": {\"spec\": {\"volumes\": [{\"name\": \"VALUEN\", \"configMap\": {\"name\": \"VALUECM\"}}], \"containers\": [{\"name\": \"VALUE\", \"image\": \"octahub.8lab.cn:5000/PROJECT/IMAGE:TAG\", \"resources\": {\"limits\": {\"cpu\": \"100m\", \"memory\": \"500Mi\"}, \"requests\": {\"cpu\": \"100m\", \"memory\": \"500Mi\"}}, \"volumeMounts\": [{\"name\": \"VALUEN\", \"subPath\": \"FILEVALUE\", \"mountPath\": \"PATHVALUE\"}]}], \"imagePullSecrets\": [{\"name\": \"registry-secret\"}]}, \"metadata\": {\"labels\": {\"app.oscro.io/name\": \"VALUE\"}}}}, \"metadata\": {\"name\": \"VALUE\", \"labels\": {\"app.oscro.io/name\": \"VALUE\"}}, \"apiVersion\": \"apps/v1\"}', NULL, '2022-11-30 01:57:24.000', NULL, 0, NULL);
INSERT INTO `base_yaml`(`id`, `name`, `kind`, `content`, `create_by`, `edit_at`, `edit_by`, `space`, `request_id`) VALUES (25, '服务', 'Service', '{\"kind\": \"Service\", \"spec\": {\"type\": \"ClusterIP\", \"ports\": [{\"name\": \"8080-8080\", \"port\": 8080, \"protocol\": \"TCP\", \"targetPort\": 8080}], \"selector\": {\"app.oscro.io/name\": \"VALUE\"}, \"externalIPs\": [\"10.0.0.10\"]}, \"metadata\": {\"name\": \"VALUE\", \"labels\": {\"app.oscro.io/name\": \"VALUE\"}, \"namespace\": \"VLAUE\"}, \"apiVersion\": \"v1\"}', NULL, '2022-11-30 01:57:01.000', NULL, 0, NULL);
INSERT INTO `base_yaml`(`id`, `name`, `kind`, `content`, `create_by`, `edit_at`, `edit_by`, `space`, `request_id`) VALUES (27, '配置', 'Configmap', '{\"data\": {\"KEY1\": \"VALUE111\", \"KEY2\": \"VALUE2\", \"file.name\": \"{\\n  \\\"DB\\\": {\\n    \\\"dbHost\\\": \\\"XX\\\",\\n    \\\"dbPort\\\": 3306,\\n    \\\"dbName\\\": \\\"XX\\\",\\n    \\\"dbUser\\\": \\\"XXX\\\",\\n    \\\"dbPasswd\\\": \\\"XXX\\\"\\n  }\\n}\\n\"}, \"kind\": \"ConfigMap\", \"metadata\": {\"name\": \"VALUE\", \"labels\": {\"app.oscro.io/name\": \"VALUE\"}, \"namespace\": \"VALUE\"}, \"apiVersion\": \"v1\"}', NULL, '2022-11-30 01:56:47.000', NULL, 0, NULL);
INSERT INTO `base_yaml`(`id`, `name`, `kind`, `content`, `create_by`, `edit_at`, `edit_by`, `space`, `request_id`) VALUES (29, '存储卷', 'PersistentVolumeClaim', '{\"kind\": \"PersistentVolumeClaim\", \"spec\": {\"resources\": {\"requests\": {\"storage\": \"VALUEGi\"}}, \"volumeMode\": \"Filesystem\", \"accessModes\": [\"ReadWriteOnce\"], \"storageClassName\": \"csi-rbd-sc\"}, \"metadata\": {\"name\": \"VALUE\"}, \"apiVersion\": \"v1\"}', NULL, '2022-12-05 02:26:20.000', NULL, 0, NULL);
INSERT INTO `base_yaml`(`id`, `name`, `kind`, `content`, `create_by`, `edit_at`, `edit_by`, `space`, `request_id`) VALUES (33, '定时任务', 'CronJob', '{\"apiVersion\":\"batch/v1\",\"kind\":\"CronJob\",\"metadata\":{\"labels\":{\"app.oscro.io/name\":\"ethanim-iris-cron\"},\"name\":\"ethanim-iris-cron-job\",\"namespace\":\"l4-ethanim-test\"},\"spec\":{\"failedJobsHistoryLimit\":1,\"jobTemplate\":{\"spec\":{\"completions\":1,\"template\":{\"spec\":{\"containers\":[{\"args\":[\"python3 tasks/cron_statistics_user_daily_online_duration.py\"],\"image\":\"octahub.8lab.cn:5000/triathon/ethanim-iris-cron:v2212301710\",\"name\":\"ethanim-iris-cron\",\"resources\":{\"limits\":{\"cpu\":\"150m\",\"memory\":\"512Mi\"},\"requests\":{\"cpu\":\"50m\",\"memory\":\"50Mi\"}},\"volumeMounts\":[{\"mountPath\":\"/usr/local/cron/configs/production.py\",\"name\":\"conf\",\"subPath\":\"production.py\"}]}],\"imagePullSecrets\":[{\"name\":\"registry-secret\"}],\"restartPolicy\":\"Never\",\"volumes\":[{\"configMap\":{\"name\":\"ethanim-iris-cron-cm\"},\"name\":\"conf\"}]}}}},\"schedule\":\"0 3 * * *\",\"successfulJobsHistoryLimit\":1}}', 'admin', '2023-07-11 02:52:36.015', 'admin', 0, 'c9473092-fccb-4ea7-8b46-c747b36f4414');
INSERT INTO `base_yaml`(`id`, `name`, `kind`, `content`, `create_by`, `edit_at`, `edit_by`, `space`, `request_id`) VALUES (35, '密钥', 'Secret', '{\"data\": {\"key1\": \"VALUE\", \"key2\": \"VALUE\"}, \"kind\": \"Secret\", \"type\": \"Opaque\", \"metadata\": {\"name\": \"VALUE\", \"namespace\": \"VALUE\"}, \"apiVersion\": \"v1\"}', NULL, '2022-12-19 06:42:41.000', NULL, 0, NULL);
INSERT INTO `base_yaml`(`id`, `name`, `kind`, `content`, `create_by`, `edit_at`, `edit_by`, `space`, `request_id`) VALUES (43, '有状态服务', 'StatefulSet', '{\"kind\": \"StatefulSet\", \"spec\": {\"replicas\": 1, \"selector\": {\"matchLabels\": {\"app.oscro.io/name\": \"VALUE\"}}, \"template\": {\"spec\": {\"volumes\": [{\"name\": \"VALUEN\", \"configMap\": {\"name\": \"VALUECM\"}}], \"containers\": [{\"name\": \"VALUE\", \"image\": \"octahub.8lab.cn:5000/PROJECT/IMAGE:TAG\", \"resources\": {\"limits\": {\"cpu\": \"100m\", \"memory\": \"500Mi\"}, \"requests\": {\"cpu\": \"100m\", \"memory\": \"500Mi\"}}, \"volumeMounts\": [{\"name\": \"VALUEN\", \"subPath\": \"FILEVALUE\", \"mountPath\": \"PATHVALUE\"}]}], \"imagePullSecrets\": [{\"name\": \"registry-secret\"}]}, \"metadata\": {\"labels\": {\"app.oscro.io/name\": \"VALUE\"}}}, \"serviceName\": \"VALUESN\"}, \"metadata\": {\"name\": \"VALUE\", \"labels\": {\"app.oscro.io/name\": \"VALUE\"}}, \"apiVersion\": \"apps/v1\"}', NULL, '2023-02-21 02:23:50.000', NULL, 0, NULL);
INSERT INTO `base_yaml`(`id`, `name`, `kind`, `content`, `create_by`, `edit_at`, `edit_by`, `space`, `request_id`) VALUES (61, 'Job', 'Job', '{\"apiVersion\":\"batch/v1\",\"kind\":\"Job\",\"metadata\":{\"labels\":{\"app.oscro.io/name\":\"LABEL\"},\"name\":\"my-job\",\"namespace\":\"defalut\"},\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"command\":[\"date\"],\"image\":\"octahub.8lab.cn:5000/PROJECT/IMAGE:TAG\",\"name\":\"my-job\",\"resources\":{\"limits\":{\"cpu\":\"100m\",\"memory\":\"500Mi\"},\"requests\":{\"cpu\":\"100m\",\"memory\":\"500Mi\"}}}],\"imagePullSecrets\":[{\"name\":\"registry-secret\"}],\"restartPolicy\":\"Never\"}}}}', 'admin', '2023-07-26 07:00:27.604', 'admin', 0, 'b12f57e0-751a-451b-83f2-f14ebea9f9bf');

SET FOREIGN_KEY_CHECKS = 1;
