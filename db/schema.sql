-- gifs
CREATE TABLE IF NOT EXISTS "gifs" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, "basename" varchar, "directory" varchar, "size" integer, "md5" varchar, "shared_link_id" varchar, "created_at" datetime NOT NULL, "updated_at" datetime NOT NULL);
CREATE INDEX IF NOT EXISTS "index_gifs_on_shared_link_id" ON "gifs" ("shared_link_id");
CREATE INDEX IF NOT EXISTS "index_gifs_on_md5" ON "gifs" ("md5");
CREATE INDEX IF NOT EXISTS "index_gifs_on_basename_and_directory" ON "gifs" ("basename", "directory");

-- shared_links
CREATE TABLE IF NOT EXISTS "shared_links" ("id" varchar NOT NULL PRIMARY KEY, "gif_id" integer, "remote_path" varchar, "count" integer DEFAULT 0, "created_at" datetime NOT NULL, "updated_at" datetime NOT NULL, CONSTRAINT "fk_golang_35031788c2"
FOREIGN KEY ("gif_id")
  REFERENCES "gifs" ("id")
);
CREATE INDEX IF NOT EXISTS "index_shared_links_on_id" ON "shared_links" ("id");
CREATE INDEX IF NOT EXISTS "index_shared_links_on_gif_id" ON "shared_links" ("gif_id");
CREATE INDEX IF NOT EXISTS "index_shared_links_on_count" ON "shared_links" ("count");
