/*
 Navicat MySQL Data Transfer

 Source Server         : dcnews
 Source Server Type    : MySQL
 Source Server Version : 50636
 Source Host           : 8.8.8.8:3306
 Source Schema         : dcnews

 Target Server Type    : MySQL
 Target Server Version : 50636
 File Encoding         : 65001

 Date: 15/08/2023 15:04:24
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for dc_wx_association_table
-- ----------------------------
DROP TABLE IF EXISTS `dc_wx_association_table`;
CREATE TABLE `dc_wx_association_table`  (
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '添加时间',
  `dc_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'dc用户名称',
  `wx_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '微信用户名称',
  `wx_group` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '微信群组名称',
  `dc_channel_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'dc频道ID',
  `dc_channel_info` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'dc频道介绍',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

SET FOREIGN_KEY_CHECKS = 1;
