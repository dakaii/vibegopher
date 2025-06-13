# Factory System

This factory system separates entity creation from database persistence, making it easy to add new factory entities without duplicating database save logic.

## Key Features

- **Separation of Concerns**: Entity creation logic is separate from database saving
- **Faker Integration**: Uses `jaswdr/faker` for realistic test data generation
- **Generic Save Functions**: Reusable `Save` and `SaveMany` functions work with any entity
- **Simple API**: Easy to understand and extend

## Architecture

### Core Components

1. **`factory.go`**: Contains generic `Save[T]` and `SaveMany[T]` functions
2. **Entity Factories**: Individual files like `userfactory.go` that define entity creation logic
3. **Faker Integration**: Uses faker library for generating realistic test data

### Functions Pattern

Each entity factory follows this pattern:

- `Build{Entity}()` - Creates entity without saving
- `Build{Entity}With{Param}()` - Creates entity with specific parameters
- `Build{Entities}(count)` - Creates multiple entities without saving
- `Create{Entity}(db)` - Creates and saves entity to database
- `Create{Entity}With{Param}(db, param)` - Creates and saves entity with parameters
- `Create{Entities}(db, count)` - Creates and saves multiple entities

## Usage Examples

### User Factory

```go
// Create entities without saving to database
user := factory.BuildUser()
users := factory.BuildUsers(5)
userWithPassword := factory.BuildUserWithPassword("custom-password")

// Create and save entities to database
savedUser := factory.CreateUser(db)
savedUsers := factory.CreateUsers(db, 5)
savedUserWithPassword := factory.CreateUserWithPassword(db, "custom-password")
```

### Post Factory (Example)

```go
// Create entities without saving to database
post := factory.BuildPost()
posts := factory.BuildPosts(3)
postWithAuthor := factory.BuildPostWithAuthor(authorID)

// Create and save entities to database
savedPost := factory.CreatePost(db)
savedPosts := factory.CreatePosts(db, 3)
savedPostWithAuthor := factory.CreatePostWithAuthor(db, authorID)
```

## Adding New Factory Entities

To add a new factory entity, create a new file (e.g., `commentfactory.go`) and follow this pattern:

```go
package factory

import (
    "time"
    "github.com/google/uuid"
    "github.com/jaswdr/faker/v2"
    "gorm.io/gorm"
)

// BuildComment creates a comment entity with faker data
func BuildComment() Comment {
    fake := faker.New()

    return Comment{
        ID:        uuid.New(),
        Content:   fake.Lorem().Sentence(10),
        AuthorID:  uuid.New(),
        PostID:    uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

// CreateComment creates and saves a comment to the database
func CreateComment(db *gorm.DB) Comment {
    comment := BuildComment()

    savedComment, err := Save(db, comment)
    if err != nil {
        panic(err)
    }

    return savedComment
}

// Add more Build/Create functions as needed...
```

## Benefits

1. **No Code Duplication**: Database save logic is centralized in `factory.go`
2. **Easy to Extend**: Adding new entities requires minimal boilerplate
3. **Realistic Test Data**: Faker generates realistic data for better tests
4. **Flexible**: Can create entities with or without saving to database
5. **Type Safe**: Uses Go generics for type safety in save functions

## Dependencies

- `github.com/jaswdr/faker/v2` - For generating fake data
- `gorm.io/gorm` - For database operations
- `github.com/google/uuid` - For UUID generation
