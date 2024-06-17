-- name: InsertUser :exec
INSERT INTO users (id,
                   auth_id,
                   email,
                   name,
                   role,
                   image_url,
                   created_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW());

-- name: GetAdmin :many
SELECT *
FROM users
WHERE users.email = $1
  AND users.password = $2
  AND users.role = 'admin' LIMIT 1;

-- name: GetAllUsers :many
SELECT *
from users
where deleted_at is null;

-- name: GetUserByAuthID :one
SELECT *
FROM users
WHERE users.auth_id = $1 LIMIT 1;
-- name: UpdateUser :exec
UPDATE users
SET name       = $1,
    image_url  = $2,
    updated_at = NOW()
WHERE id = $3;
-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1;

-- name: InsertCategory :exec
INSERT INTO categories (id, name, created_at)
VALUES ($1, $2, NOW());
-- name: UpdateCategory :exec
UPDATE categories
SET name       = $1,
    updated_at = NOW()
WHERE id = $2;
-- name: DeleteCategory :exec
UPDATE categories
SET deleted_at = NOW()
WHERE id = $1;

-- name: GetAllCategories :many
SELECT *
from categories
where deleted_at is null;

-- name: InsertNews :exec
INSERT INTO news (id, author, title, description, content, url, image_url, publish_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW());

-- name: UpdateNews :exec
UPDATE news
SET title       = $1,
    description = $2,
    content     = $3,
    author      = $4,
    url         = $5,
    image_url   = $6,
    publish_at  = $7,
    updated_at  = NOW()
WHERE id = $8;
-- name: DeleteNews :exec
UPDATE news
SET deleted_at = NOW()
WHERE id = $1;

-- name: GetAllNews :many
SELECT n.id                           AS id,
       n.author,
       n.title,
       n.description,
       n.content,
       n.url,
       n.image_url,
       n.publish_at,
       n.created_at                   AS created_at,
       n.updated_at                   AS updated_at,
       n.deleted_at                   AS deleted_at,
       json_agg(hc.category_id::uuid) AS category_ids,
       COALESCE(COUNT(Distinct v.user_id), 0)  AS view_count
FROM news n
         Left JOIN has_categories hc
                   ON n.id = hc.news_id
         LEFT JOIN views v ON n.id = v.news_id
WHERE n.deleted_at is null
GROUP BY n.id,
         n.author,
         n.title,
         n.description,
         n.content,
         n.url,
         n.image_url,
         n.publish_at,
         n.created_at,
         n.updated_at,
         n.deleted_at;

-- name: InsertLike :exec
INSERT INTO likes (news_id, user_id)
VALUES ($1, $2);

-- name: DeleteLike :exec
DELETE
from likes
Where news_id = $1
  and user_id = $2;

-- name: InsertDisLike :exec
INSERT INTO dislikes (news_id, user_id)
VALUES ($1, $2);

-- name: DeleteDisLike :exec
DELETE
from dislikes
Where news_id = $1
  and user_id = $2;

-- name: InsertHasCategory :exec
INSERT INTO has_categories (news_id, category_id)
VALUES ($1, $2);

-- name: InsertSave :exec
Insert into saves (news_id, user_id)
values ($1, $2);

-- name: DeleteSave :exec
DELETE
from saves
Where news_id = $1
  and user_id = $2;

-- name: GetSaves :many
SELECT n.id                           AS id,
       n.author,
       n.title,
       n.description,
       n.content,
       n.url,
       n.image_url,
       n.publish_at,
       n.created_at                   AS created_at,
       n.updated_at                   AS updated_at,
       n.deleted_at                   AS deleted_at,
       json_agg(hc.category_id::uuid) AS category_ids
FROM news n
         Left JOIN has_categories hc ON n.id = hc.news_id
         join saves s on s.news_id = n.id
where s.user_id = $1
  and n.deleted_at is null
GROUP BY n.id,
         n.author,
         n.title,
         n.description,
         n.content,
         n.url,
         n.image_url,
         n.publish_at,
         n.created_at,
         n.updated_at,
         n.deleted_at;

-- name: GetLike :one
SELECT *
from likes
Where news_id = $1
  and user_id = $2;

-- name: GetDislike :one
SELECT *
from dislikes
Where news_id = $1
  and user_id = $2;

-- name: GetNews :one
SELECT n.id                           AS id,
       n.author,
       n.title,
       n.description,
       n.content,
       n.url,
       n.image_url,
       n.publish_at,
       n.created_at                   AS created_at,
       n.updated_at                   AS updated_at,
       n.deleted_at                   AS deleted_at,
       json_agg(hc.category_id::uuid) AS category_ids,
       COALESCE(COUNT(Distinct v.user_id), 0)  AS view_count
FROM news n
         Left JOIN has_categories hc
                   ON n.id = hc.news_id
         LEFT JOIN views v ON n.id = v.news_id
where id = $1
  and deleted_at is null
GROUP BY n.id,
         n.author,
         n.title,
         n.description,
         n.content,
         n.url,
         n.image_url,
         n.publish_at,
         n.created_at,
         n.updated_at,
         n.deleted_at;

-- name: DeleteHasCategory :exec
DELETE
from has_categories
where news_id = $1;

-- name: SearchNews :many
SELECT n.id                           AS id,
       n.author                       AS author,
       n.title                        AS title,
       n.description                  AS description,
       n.content,
       n.url,
       n.image_url,
       n.publish_at,
       n.created_at                   AS created_at,
       n.updated_at                   AS updated_at,
       n.deleted_at                   AS deleted_at,
       json_agg(hc.category_id::uuid) AS category_ids,
       COALESCE(COUNT(Distinct v.user_id), 0)  AS view_count
FROM news n
         Left JOIN has_categories hc
                   ON n.id = hc.news_id
         LEFT JOIN views v ON n.id = v.news_id
WHERE deleted_at is null
  and (
    author LIKE '%' || $1 || '%'
        OR description LIKE '%' || $1 || '%'
        OR title LIKE '%' || $1 || '%'
    )

GROUP BY n.id,
         n.author,
         n.title,
         n.description,
         n.content,
         n.url,
         n.image_url,
         n.publish_at,
         n.created_at,
         n.updated_at,
         n.deleted_at;

-- name: SearchUsers :many
Select *
from users
where email LIKE '%' || $1 || '%'
   or name LIKE '%' || $1 || '%' and deleted_at is null;

-- name: SearchCategories :many
Select *
from Categories
where name LIKE '%' || $1 || '%'
  and deleted_at is null;

-- name: GetNewsByIds :many
SELECT n.id                          AS id,
       n.author,
       n.title,
       n.description,
       n.content,
       n.url,
       n.image_url,
       n.publish_at,
       n.created_at                  AS created_at,
       n.updated_at                  AS updated_at,
       n.deleted_at                  AS deleted_at,
       COALESCE(COUNT(Distinct v.user_id), 0) AS view_count
FROM news n
         LEFT JOIN views v ON n.id = v.news_id
WHERE n.id = ANY ($1::uuid[])
  AND n.deleted_at IS NULL
GROUP BY n.id,
         n.author,
         n.title,
         n.description,
         n.content,
         n.url,
         n.image_url,
         n.publish_at,
         n.created_at,
         n.updated_at,
         n.deleted_at;


-- name: InsertComment :exec
Insert INTO comments (news_id, user_id, text, published_at)
values ($1, $2, $3, NOW());

-- name: QueryCommentByNews :many
select c.text, c.published_at, u.name, image_url
from comments c
         JOIN users u
              on c.user_id = u.id
where news_id = $1
order BY c.published_at DESC;

-- name: InsertView :exec
Insert into views (news_id, user_id)
values ($1, $2);

