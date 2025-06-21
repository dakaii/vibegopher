// Atlas configuration for vibegopher project
env "dev" {
  src = "file://schema.sql"
  url = "postgres://postgres:postgres@vibegopher-postgresql-dev1:5432/vibegopher_development?sslmode=disable"
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
  migration {
    dir = "file://atlas-migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
