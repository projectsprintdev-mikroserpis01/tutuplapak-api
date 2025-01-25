ALTER TABLE users
ADD CONSTRAINT fk_file_id
FOREIGN KEY (file_id) REFERENCES files(id);
