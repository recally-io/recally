// Infrastructure Schema - Cache Table
// This file defines the cache table used for application-level caching

schema "public" {
}

table "cache" {
  schema = schema.public

  column "id" {
    null = false
    type = serial
  }

  column "domain" {
    null = false
    type = varchar(255)
  }

  column "key" {
    null = false
    type = text
  }

  column "value" {
    null = false
    type = jsonb
  }

  column "expires_at" {
    null = true
    type = timestamptz
  }

  column "created_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  index "uni_cache_domain_key" {
    unique  = true
    columns = [column.domain, column.key]
  }
}
