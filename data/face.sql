/*
 Navicat Premium Data Transfer

 Source Server         : 140.143.234.132
 Source Server Type    : MySQL
 Source Server Version : 50639
 Source Host           : 140.143.234.132:3306
 Source Schema         : face

 Target Server Type    : MySQL
 Target Server Version : 50639
 File Encoding         : 65001

 Date: 03/04/2018 14:09:56
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for vivi_user
-- ----------------------------
DROP TABLE IF EXISTS `vivi_user`;
CREATE TABLE `vivi_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(60) NOT NULL COMMENT '用户名',
  `password` char(40) NOT NULL COMMENT '密码',
  `face_token` varchar(200) DEFAULT NULL COMMENT '用户的头像tolen',
  `face_url` varchar(200) DEFAULT NULL COMMENT '脸部图片地址',
  `faceset_token` varchar(200) DEFAULT NULL COMMENT 'faceset_token',
  `create_time` int(11) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb4 COMMENT='系统登录用户信息表';

SET FOREIGN_KEY_CHECKS = 1;
