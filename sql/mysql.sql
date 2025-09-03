DROP TABLE IF EXISTS `t_system_users`;
create table t_system_users
(
    id         varchar(32) comment '用户ID'
        primary key,
    username   varchar(32)                            not null comment '用户账号',
    password   varchar(128) default ''                not null comment '密码',
    nickname   varchar(32)                            not null comment '用户昵称',
    remark     varchar(512)                           null comment '备注',
    dept_id    varchar(32)                                 null comment '部门ID',
    post_ids   varchar(256)                           null comment '岗位编号数组',
    email      varchar(64)  default ''                null comment '用户邮箱',
    mobile     varchar(16)  default ''                null comment '手机号码',
    sex        tinyint      default 0                 null comment '用户性别',
    avatar     varchar(512) default ''                null comment '头像地址',
    status     tinyint      default 0                 not null comment '帐号状态（0停用 1正常）',
    login_ip   varchar(64)  default ''                null comment '最后登录IP',
    login_date datetime                               null comment '最后登录时间',
    create_by  varchar(64)  default ''                null comment '创建者',
    create_at  datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    update_by  varchar(64)  default ''                null comment '更新者',
    update_at  datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    delete_at  datetime                               null comment '删除时间',
    tenant_id  varchar(32)       default ''                 not null comment '租户编号',
    constraint idx_username
        unique (username, update_at, tenant_id)
)
    comment '用户信息表' collate = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `t_system_role`;
create table t_system_role
(
    id                  varchar(32) comment '角色ID'
        primary key,
    name                varchar(32)                            not null comment '角色名称',
    code                varchar(128)                           not null comment '角色权限字符串',
    sort                int                                    not null comment '显示顺序',
    data_scope          tinyint      default 1                 not null comment '数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）',
    data_scope_dept_ids varchar(512) default ''                not null comment '数据范围(指定部门数组)',
    status              tinyint                                not null comment '角色状态（0停用 1正常）',
    type                tinyint                                not null comment '角色类型',
    remark              varchar(512)                           null comment '备注',
    create_by           varchar(64)  default ''                null comment '创建者',
    create_at           datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    update_by           varchar(64)  default ''                null comment '更新者',
    update_at           datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    delete_at           bit                         null comment '删除时间',
    tenant_id           varchar(32)       default 0                 not null comment '租户编号'
)
    comment '角色信息表' collate = utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for t_system_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `t_system_role_menu`;
create table t_system_role_menu
(
    id        varchar(32) comment '自增编号'
        primary key,
    role_id   varchar(32)                                not null comment '角色ID',
    menu_id   varchar(32)                                not null comment '菜单ID',
    create_by varchar(64) default ''                null comment '创建者',
    create_at datetime    default CURRENT_TIMESTAMP not null comment '创建时间',
    update_by varchar(64) default ''                null comment '更新者',
    update_at datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    delete_at datetime                               comment '删除时间',
    tenant_id varchar(32)      default ''                 not null comment '租户编号'
)
    comment '角色和菜单关联表' collate = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `t_system_menu`;
create table t_system_menu
(
    id             varchar(32) comment '菜单ID'
        primary key,
    name           varchar(64)                            not null comment '菜单名称',
    permission     varchar(128) default ''                not null comment '权限标识',
    type           tinyint                                not null comment '菜单类型',
    sort           int          default 0                 not null comment '显示顺序',
    parent_id      varchar(32)       default 0                 not null comment '父菜单ID',
    path           varchar(200) default ''                null comment '路由地址',
    icon           varchar(128) default '#'               null comment '菜单图标',
    component      varchar(256)                           null comment '组件路径',
    component_name varchar(256)                           null comment '组件名',
    status         tinyint      default 0                 not null comment '菜单状态',
    visible        bit          default b'1'              not null comment '是否可见',
    keep_alive     bit          default b'1'              not null comment '是否缓存',
    always_show    bit          default b'1'              not null comment '是否总是显示',
    create_by      varchar(64)  default ''                null comment '创建者',
    create_at      datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    update_by      varchar(64)  default ''                null comment '更新者',
    update_at      datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    delete_at      datetime                         null comment '删除时间'
)
    comment '菜单权限表' collate = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `t_system_user_role`;
create table t_system_user_role
(
    id        varchar(32) comment '自增编号'
        primary key,
    user_id   varchar(32)                                not null comment '用户ID',
    role_id   varchar(32)                                not null comment '角色ID',
    create_by varchar(64) default ''                null comment '创建者',
    create_at datetime    default CURRENT_TIMESTAMP null comment '创建时间',
    update_by varchar(64) default ''                null comment '更新者',
    update_at datetime    default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    delete_at datetime                  null comment '删除时间',
    tenant_id varchar(32)      default ''                 not null comment '租户编号'
)
    comment '用户和角色关联表' collate = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `t_system_post`;
create table t_system_post
(
    id        varchar(32) comment '岗位ID'
        primary key,
    code      varchar(64)                           not null comment '岗位编码',
    name      varchar(64)                           not null comment '岗位名称',
    sort      int                                   not null comment '显示顺序',
    status    tinyint                               not null comment '状态（0停用 1正常）',
    remark    varchar(512)                          null comment '备注',
    create_by varchar(64) default ''                null comment '创建者',
    create_at datetime    default CURRENT_TIMESTAMP not null comment '创建时间',
    update_by varchar(64) default ''                null comment '更新者',
    update_at datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    delete_at datetime                              null comment '删除时间',
    tenant_id varchar(32)      default ''                 not null comment '租户编号'
)
    comment '岗位信息表' collate = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `t_system_user_post`;
create table t_system_user_post
(
    id        varchar(32) comment 'id'
        primary key,
    user_id   varchar(32)                       not null comment '用户ID',
    post_id   varchar(32)                       not null comment '岗位ID',
    create_by varchar(64) default ''                null comment '创建者',
    create_at datetime    default CURRENT_TIMESTAMP not null comment '创建时间',
    update_by varchar(64) default ''                null comment '更新者',
    update_at datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    delete_at datetime                              null comment '删除时间',
    tenant_id varchar(32)      default ''                 not null comment '租户编号'
)
    comment '用户岗位表' collate = utf8mb4_unicode_ci;


create table t_system_dept
(
    id             varchar(32) comment '部门id'
        primary key,
    name           varchar(32) default ''                not null comment '部门名称',
    parent_id      varchar(32)                       not null comment '父部门id',
    sort           int         default 0                 not null comment '显示顺序',
    leader_user_id varchar(32)                                null comment '负责人',
    phone          varchar(16)                           null comment '联系电话',
    email          varchar(64)                           null comment '邮箱',
    status         tinyint                               not null comment '部门状态（0停用 1正常）',
    create_by      varchar(64) default ''                null comment '创建者',
    create_at      datetime    default CURRENT_TIMESTAMP not null comment '创建时间',
    update_by      varchar(64) default ''                null comment '更新者',
    update_at      datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    delete_at      datetime                              null comment '删除时间',
    tenant_id      varchar(32)      default ''                 not null comment '租户编号'
)
    comment '部门表' collate = utf8mb4_unicode_ci;


create table t_system_tenant
(
    id              varchar(32) comment '租户编号'
        primary key,
    name            varchar(32)                            not null comment '租户名',
    contact_user_id varchar(32)                                 null comment '联系人的用户编号',
    contact_name    varchar(32)                            not null comment '联系人',
    contact_mobile  varchar(512)                           null comment '联系手机',
    status          tinyint      default 1                 not null comment '租户状态（0停用 1正常）',
    website         varchar(256) default ''                null comment '绑定域名',
    package_id      varchar(32)                                 not null comment '租户套餐编号',
    expire_time     datetime                               not null comment '过期时间',
    account_count   int                                    not null comment '账号数量',
    create_by       varchar(64)  default ''                not null comment '创建者',
    create_at       datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    update_by       varchar(64)  default ''                null comment '更新者',
    update_at       datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    delete_at       datetime                               null comment '删除时间'
)
    comment '租户表' collate = utf8mb4_unicode_ci;

create table t_system_tenant_package
(
    id        varchar(32) comment '套餐编号'
        primary key,
    name      varchar(32)                            not null comment '套餐名',
    status    tinyint      default 1                 not null comment '租户状态（0停用 1正常）',
    remark    varchar(256) default ''                null comment '备注',
    menu_ids  varchar(4096)                          not null comment '关联的菜单编号',
    create_by varchar(64)  default ''                not null comment '创建者',
    create_at datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    update_by varchar(64)  default ''                null comment '更新者',
    update_at datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    delete_at datetime                               null comment '删除时间'
)
    comment '租户套餐表' collate = utf8mb4_unicode_ci;