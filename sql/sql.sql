create table friend
(
    id          bigint auto_increment comment '''自增主键'''
        primary key,
    user_id     bigint                             not null comment '''用户id''',
    friend_id   bigint                             not null comment '''好友id''',
    create_time datetime default CURRENT_TIMESTAMP not null comment '''创建时间''',
    update_time datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '''更新时间'''
);

create table `group`
(
    id          bigint auto_increment comment '''自增主键'''
        primary key,
    name        longtext                           not null comment '''群组名称''',
    owner_id    bigint                             not null comment '''群主id''',
    create_time datetime default CURRENT_TIMESTAMP not null comment '''创建时间''',
    update_time datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '''更新时间'''
);

create table group_user
(
    id          bigint auto_increment comment '''自增主键'''
        primary key,
    group_id    bigint                             not null comment '''组id''',
    user_id     bigint                             not null comment '''用户id''',
    create_time datetime default CURRENT_TIMESTAMP not null comment '''创建时间''',
    update_time datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '''更新时间'''
);

create table message
(
    id           bigint auto_increment comment '''自增主键'''
        primary key,
    user_id      bigint                             not null comment '''用户id，指接受者用户id''',
    sender_id    bigint                             not null comment '''发送者用户id''',
    session_type tinyint                            not null comment '''聊天类型，群聊/单聊''',
    receiver_id  bigint                             not null comment '''接收者id，群聊id/用户id''',
    message_type tinyint                            not null comment '''消息类型,语言、文字、图片''',
    content      longblob                           not null comment '''消息内容''',
    seq          bigint                             not null comment '''消息序列号''',
    send_time    datetime default CURRENT_TIMESTAMP not null comment '''消息发送时间''',
    create_time  datetime default CURRENT_TIMESTAMP not null comment '''创建时间''',
    update_time  datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '''更新时间'''
);

create table user
(
    id           bigint auto_increment comment '''自增主键'''
        primary key,
    phone_number varchar(191)                       not null comment '''手机号''',
    nickname     longtext                           not null comment '''昵称''',
    password     longtext                           not null comment '''密码''',
    create_time  datetime default CURRENT_TIMESTAMP not null comment '''创建时间''',
    update_time  datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '''更新时间''',
    constraint uni_user_phone_number
        unique (phone_number)
);

create table users
(
    id           bigint auto_increment comment '''自增主键'''
        primary key,
    phone_number varchar(191)                       not null comment '''手机号''',
    nickname     longtext                           not null comment '''昵称''',
    password     longtext                           not null comment '''密码''',
    create_time  datetime default CURRENT_TIMESTAMP not null comment '''创建时间''',
    update_time  datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '''更新时间''',
    constraint uni_users_phone_number
        unique (phone_number)
);

