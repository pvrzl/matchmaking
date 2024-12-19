-- Remove comments
COMMENT ON COLUMN "subscriptions"."subscription_type" IS NULL;
COMMENT ON COLUMN "match"."status" IS NULL;

-- Drop indexes
DROP INDEX IF EXISTS "match_swiped_id_idx";
DROP INDEX IF EXISTS "match_swiper_id_idx";
DROP INDEX IF EXISTS "match_swiper_id_swiped_id_idx";
DROP INDEX IF EXISTS "users_email_idx";
DROP INDEX IF EXISTS "subscriptions_user_id_subscription_type_expires_at_idx";

-- Drop tables
DROP TABLE IF EXISTS "configurations";
DROP TABLE IF EXISTS "subscriptions";
DROP TABLE IF EXISTS "match";
DROP TABLE IF EXISTS "users";
