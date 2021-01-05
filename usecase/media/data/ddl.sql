--  媒体表[0-49]
CREATE TABLE "component_media" (
    "increment_id" int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    "deleted" tinyint unsigned DEFAULT 0 COMMENT '0:未删除 1:删除',
    "create_time" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    "update_time" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    "delete_time" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '删除时间',
    PRIMARY KEY ("increment_id")
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='组件——媒体表[0-49]';
