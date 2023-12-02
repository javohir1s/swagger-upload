
CREATE TABLE "timezone" (
    "guid" UUID PRIMARY KEY,
    "title" VARCHAR(24),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);
