-- Добавляем внешние ключи
ALTER TABLE posts
ADD CONSTRAINT fk_posts_author
FOREIGN KEY (author_id)
REFERENCES users(id)
ON DELETE CASCADE;

ALTER TABLE comments
ADD CONSTRAINT fk_comments_post
FOREIGN KEY (post_id)
REFERENCES posts(id)
ON DELETE CASCADE;

ALTER TABLE comments
ADD CONSTRAINT fk_comments_author
FOREIGN KEY (author_id)
REFERENCES users(id)
ON DELETE CASCADE;
