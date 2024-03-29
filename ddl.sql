create extension "uuid-ossp";

create table users
(
    id                 uuid        default uuid_generate_v4() not null
        constraint users_pk
            primary key,
    username           integer
        constraint users_username_uk
            unique,
    following_count    bigint default 0,
    followers_count    bigint default 0,
    bio                text,
    posts_count        bigint default 0,
    highlights_count   bigint default 0,
    name               varchar,
    profile_image_link text,
    is_business        bool,
    is_verified        bool,
    country            varchar,
    region             varchar,
    created_at         timestamptz default current_timestamp  not null,
    updated_at         timestamptz,
    deleted_at         timestamptz
);

create table businesses
(
    id             uuid        default uuid_generate_v4() not null
        constraint businesses_pk
            primary key,
    user_id        uuid                                   not null
        constraint businesses_users_id_fk
            references users,
    city_name      varchar,
    latitude       float8,
    longitude      float8,
    street_address text,
    zip_code       integer,
    created_at     timestamptz default current_timestamp  not null,
    updated_at     timestamptz,
    deleted_at     timestamptz
);

create table followers
(
    follower_id  uuid                                  not null
        constraint followers_users_id_fk
            references users,
    following_id uuid                                  not null
        constraint followers_users_id_fk2
            references users,
    followed_at  timestamptz default current_timestamp not null,
    constraint followers_pk
        primary key (following_id, follower_id)
);

create table followers_activity
(
    id           uuid        default uuid_generate_v4() not null
        constraint followers_activity_pk
            primary key,
    follower_id  uuid                                   not null
        constraint followers_activity_users_id_fk
            references users,
    following_id uuid                                   not null
        constraint followers_activity_users_id_fk2
            references users,
    is_unfollow  bool        default false              not null,
    created_at   timestamptz default current_timestamp  not null
);

create table locations
(
    id              bigint          not null
        constraint locations_pk
            primary key,
    has_public_page bool default false not null,
    name            varchar            not null
        constraint locations_name_uk
            unique,
    slug            varchar            not null
        constraint locations_slug_uk
            unique
);

create table posts
(
    id                bigint                                not null
        constraint posts_pk
            primary key,
    user_id           uuid                                  not null
        constraint posts_users_id_fk
            references users,
    caption           text,
    likes_count       bigint      default 0                 not null,
    comments_count    bigint      default 0                 not null,
    video_view_count  bigint      default 0                 not null,
    primary_image_url text,
    primary_video_url text,
    location_id       bigint
        constraint posts_locations_id_fk
            references locations,
    is_sponsored      bool        default false             not null,
    sponsor_id        uuid
        constraint posts_sponsor_id_fk
            references users,
    url               text                                  not null,
    created_at        timestamptz default current_timestamp not null,
    updated_at        timestamptz,
    deleted_at        timestamptz
);

create table post_images
(
    id         bigserial         not null
        constraint post_images_pk
            primary key,
    post_id    bigint            not null
        constraint post_images_posts_id_fk
            references posts,
    image_url  text              not null,
    post_order integer default 1 not null,
    created_at        timestamptz default current_timestamp not null,
    updated_at        timestamptz,
    deleted_at        timestamptz
);

create table highlights
(
    id      bigserial not null
        constraint highlights_pk
            primary key,
    user_id uuid      not null
        constraint highlights_users_id_fk
            references users,
    title   varchar,
    image   text      not null,
    created_at        timestamptz default current_timestamp not null,
    updated_at        timestamptz,
    deleted_at        timestamptz
);

create table stories
(
    id         uuid                                  not null
        constraint stories_pk
            primary key,
    user_id    uuid                                  not null
        constraint stories_fk
            references users,
    media_url  text                                  not null,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz
);

create table story_views
(
    story_id  uuid                                  not null
        constraint story_views_stories_id_fk
            references stories,
    viewer_id uuid                                  not null
        constraint story_views_users_id_fk
            references users,
    is_liked  bool        default false             not null,
    viewed_at timestamptz default current_timestamp not null,
    constraint story_views_pk
        primary key (story_id, viewer_id)
);

create table highlights_stories
(
    highlight_id bigint
        constraint highlights_stories_highlights_id_fk
            references highlights,
    story_id     uuid
        constraint highlights_stories_stories_id_fk
            references stories,
    created_at   timestamptz default current_timestamp not null,
    updated_at   timestamptz,
    deleted_at   timestamptz,
    constraint highlights_stories_pk
        primary key (highlight_id, story_id)
);

create table highlights_story_activity
(
    highlight_id bigint
        constraint highlights_story_activity_highlights_id_fk
            references highlights,
    story_id     uuid
        constraint highlights_story_activity_stories_id_fk
            references stories,
    is_removed   bool        default false             not null,
    created_at   timestamptz default current_timestamp not null,
    constraint highlights_story_activity_pk
        primary key (story_id, highlight_id)
);

create table hash_tags
(
    id         bigint
        constraint hash_tags_pk
            primary key,
    name       text                                  not null,
    created_by uuid constraint  hash_tags_users_id_fk references users,
    is_blocked bool        default false             not null,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz
);

create table post_tags
(
    post_id    bigint
        constraint post_tags_posts_id_fk
            references posts,
    tag_id     bigint
        constraint post_tags_hash_tags_id_fk
            references hash_tags,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz,
    constraint post_tags_pk
        primary key (post_id, tag_id)
);

create table story_tags
(
    story_id   uuid
        constraint story_tags_stories_id_fk
            references stories,
    tag_id     bigint
        constraint story_tags_hash_tags_id_fk
            references hash_tags,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz,
    constraint story_tags_pk
        primary key (story_id, tag_id)
);

create table block
(
    user_id    uuid
        constraint block_users_id_fk
            references users,
    blocked_id uuid
        constraint block_users_id_fk2
            references users,
    blocked_at timestamptz,
    constraint block_pk
        primary key (user_id, blocked_id)
);

create table block_activity
(
    user_id    uuid
        constraint block_activity_users_id_fk
            references users,
    blocked_id uuid
        constraint block_activity_users_id_fk2
            references users,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz,
    is_block   bool        default true              not null
);

create table restrict
(
    user_id          uuid                                  not null
        constraint restrict_users_id_fk
            references users,
    restrict_user_id uuid                                  not null
        constraint restrict_user_id_user_id_fk
            references users,
    created_at       timestamptz default current_timestamp not null,
    updated_at       timestamptz,
    deleted_at       timestamptz,
    constraint restrict_pk
        primary key (user_id, restrict_user_id)
);

create table restrict_activity
(
    user_id          uuid                                  not null
        constraint restrict_activity_users_id_fk
            references users,
    restrict_user_id uuid                                  not null
        constraint restrict_activity_users_id_fk2
            references users,
    is_restrict      bool        default true              not null,
    created_at       timestamptz default current_timestamp not null
);

create table comments
(
    id                bigint
        constraint comments_pk
            primary key,
    post_id           bigint
        constraint comments_posts_id_fk
            references posts,
    user_id           uuid
        constraint comments_users_id_fk
            references users,
    parent_comment_id bigint,
    comment_text      text                                  not null,
    created_at        timestamptz default current_timestamp not null,
    updated_at        timestamptz,
    deleted_at        timestamptz
);

create table comment_likes
(
    comment_id bigint                                not null
        constraint comment_likes_comments_id_fk
            references comments,
    liked_by   uuid                                  not null
        constraint comment_likes_users_id_fk
            references users,
    liked_at   timestamptz default current_timestamp not null,
    constraint comment_likes_pk
        primary key (comment_id, liked_by)
);

create table comment_activity
(
    comment_id bigint                                not null
        constraint comment_activity_comments_id_fk
            references comments,
    action_by  uuid                                  not null
        constraint comment_activity_users_id_fk
            references users,
    is_like    bool        default true              not null,
    created_at timestamptz default current_timestamp not null
);

