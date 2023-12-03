CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "fio" VARCHAR(255),
    "phone_or_email" VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NOT NULL
);