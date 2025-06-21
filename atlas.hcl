// Atlas configuration for vibegopher project
env "dev" {
  src = "file://schema.sql"
  url = "postgres://postgres:postgres@vibegopher-postgresql-dev1:5432/vibegopher_development?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"
  migration {
    dir = "file://atlas-migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "test" {
  src = "file://schema.sql"
  url = "postgres://postgres:postgres@postgresql-test:5432/vibegopher_test?sslmode=disable"
  dev = "docker://postgres/15/test?search_path=public"
  migration {
    dir = "file://atlas-migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "prod" {
  src = "file://schema.sql"
  url = env("DATABASE_URL")
  dev = "docker://postgres/15/prod?search_path=public"
  migration {
    dir = "file://atlas-migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
