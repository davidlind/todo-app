-- Write your migrate up statements here
create table webb (id serial, content text not null);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
