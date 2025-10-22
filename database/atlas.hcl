env "local" {
  src = "file://schema"
  url = getenv("DATABASE_URL")
  dev = "docker://postgres/16/dev?search_path=public"

  migration {
    dir = "file://migrations"
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
  src = "file://schema"
  url = getenv("DATABASE_URL")

  migration {
    dir = "file://migrations"
  }

  lint {
    destructive {
      error = true
    }
  }
}
