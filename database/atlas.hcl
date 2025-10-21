env "local" {
  src = "file://database/schema"
  url = getenv("DATABASE_URL")
  dev = "docker://paradedb/paradedb:latest-pg16"  # CRITICAL: Use ParadeDB, not vanilla Postgres

  migration {
    dir = "file://database/migrations"
  }

  # Ignore River's internal migration tables
  exclude {
    schema_pattern = "^river_.*"
    table_pattern  = "^river_.*"
  }

  lint {
    destructive {
      error = true  # Prevent accidental drops in production
    }
  }

  diff {
    skip {
      drop_schema = true  # Safety: never auto-drop schemas
    }
  }
}

env "production" {
  src = "file://database/schema"
  url = getenv("DATABASE_URL")

  migration {
    dir = "file://database/migrations"
  }

  exclude {
    schema_pattern = "^river_.*"
    table_pattern  = "^river_.*"
  }

  lint {
    destructive {
      error = true
    }
  }
}
