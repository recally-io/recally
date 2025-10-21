// Authentication and Authorization Schema
// This file defines user management, OAuth integration, API keys, and token revocation tables


table "users" {
  schema = schema.public

  column "id" {
    null = false
    type = serial
  }

  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "username" {
    null = true
    type = varchar(255)
  }

  column "password_hash" {
    null = true
    type = text
  }

  column "email" {
    null = true
    type = varchar(255)
  }

  column "phone" {
    null = true
    type = varchar(50)
  }

  column "activate_assistant_id" {
    null = true
    type = uuid
  }

  column "activate_thread_id" {
    null = true
    type = uuid
  }

  column "status" {
    null    = false
    type    = varchar(255)
    default = sql("'pending'")
  }

  column "settings" {
    null    = true
    type    = jsonb
    default = sql("'{}'::JSONB")
  }

  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  index "users_uuid_key" {
    unique  = true
    columns = [column.uuid]
  }

  index "users_email_key" {
    unique  = true
    columns = [column.email]
  }

  index "idx_users_email" {
    unique = true
    on {
      expr = "LOWER(email)"
    }
    where = "email IS NOT NULL"
  }

  index "idx_users_phone" {
    unique  = true
    columns = [column.phone]
    where   = "phone IS NOT NULL"
  }

  index "idx_users_username" {
    unique  = true
    columns = [column.username]
    where   = "username IS NOT NULL"
  }

  check "users_email_check" {
    expr = "email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$'"
  }

  check "users_contact_check" {
    expr = "email IS NOT NULL OR phone IS NOT NULL OR username IS NOT NULL"
  }
}

table "auth_user_oauth_connections" {
  schema = schema.public

  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    null = false
    type = uuid
  }

  column "provider" {
    null = false
    type = varchar(50)
  }

  column "provider_user_id" {
    null = false
    type = varchar(255)
  }

  column "provider_email" {
    null = true
    type = varchar(255)
  }

  column "access_token" {
    null = true
    type = text
  }

  column "refresh_token" {
    null = true
    type = text
  }

  column "token_expires_at" {
    null = true
    type = timestamptz
  }

  column "provider_data" {
    null = true
    type = jsonb
  }

  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "auth_user_oauth_connections_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
    on_delete   = CASCADE
  }

  index "uq_oauth_connection" {
    unique  = true
    columns = [column.provider, column.provider_user_id]
  }

  index "idx_oauth_user_id" {
    columns = [column.user_id]
  }

  index "idx_oauth_provider_lookup" {
    columns = [column.provider, column.provider_user_id]
  }

  index "idx_oauth_token_expiry" {
    columns = [column.token_expires_at]
    where   = "token_expires_at IS NOT NULL"
  }
}

table "auth_api_keys" {
  schema = schema.public

  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    null = false
    type = uuid
  }

  column "name" {
    null = false
    type = varchar(255)
  }

  column "key_prefix" {
    null = false
    type = varchar(8)
  }

  column "key_hash" {
    null = false
    type = varchar(255)
  }

  column "scopes" {
    null = false
    type = sql("TEXT[]")
  }

  column "expires_at" {
    null = true
    type = timestamptz
  }

  column "last_used_at" {
    null = true
    type = timestamptz
  }

  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "auth_api_keys_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
    on_delete   = CASCADE
  }

  index "uq_user_key_name" {
    unique  = true
    columns = [column.user_id, column.name]
  }

  index "idx_auth_api_keys_prefix" {
    columns = [column.key_prefix]
  }

  index "idx_auth_api_keys_user" {
    columns = [column.user_id]
  }

  index "idx_auth_api_keys_expiry" {
    columns = [column.expires_at]
    where   = "expires_at IS NOT NULL"
  }

  check "ck_key_expiry" {
    expr = "expires_at IS NULL OR expires_at > created_at"
  }
}

table "auth_revoked_tokens" {
  schema = schema.public

  column "jti" {
    null = false
    type = uuid
  }

  column "user_id" {
    null = false
    type = uuid
  }

  column "expires_at" {
    null = false
    type = timestamptz
  }

  column "revoked_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "reason" {
    null = true
    type = varchar(100)
  }

  primary_key {
    columns = [column.jti]
  }

  foreign_key "auth_revoked_tokens_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  index "uq_revoked_token" {
    unique  = true
    columns = [column.user_id, column.jti]
  }

  index "idx_auth_revoked_tokens_expiry" {
    columns = [column.expires_at]
  }

  index "idx_auth_revoked_tokens_user" {
    columns = [column.user_id]
  }

  check "ck_token_revocation" {
    expr = "expires_at > revoked_at"
  }
}
