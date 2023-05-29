# gopher-burrow
Accelerate the development of Go REST Services and take some of the toil out of coding.

## Overview
This is a brand new project, so there is no release yet so there will likely be breaking changes.
There is definitely not even a high-level plan as to what tools and code snippets will be included.

Some high-level principles:
* Make repetitive tasks easier
* Follow a pattern similar to Controller/Service/Repository used in Spring Boot but the code should not look like a developer trying to make GO work like Java/Spring Boot
* Follow a test-driven development approach
* Code should be clean and easy to understand

## Installation
To install gopher-burrow, use `go get`:

```go get github.com/rkennel/gopher-burrow```

## Features
### Repository
Your Entity Repository can get all of the basic Spring Repository methods:
```
FindByID(id uuid.UUID) (*T, error)
FindAll() ([]T, error)
Save(t *T) (*T, error)
Delete(id uuid.UUID) error
Exists(id uuid.UUID) (bool, error)
Count() (int64, error)   
```

To take advantage of this, perform the following steps:
- Your Entity must include the BasicFields struct
- Your Entity must implement the BasicEntity struct
- Your Repository Implementation must include the BasicRepository struct

Example (Pardon the dated GI Joe References)
```
type GiJoe struct {
	BasicFields
	firstName string
	lastName  string
	codeName  string
	jobTitle  string
}

func (gijoe GiJoe) GetId() uuid.UUID {
	return gijoe.ID
}

type GiJoeRepository struct {
	BasicRepository[GiJoe]
}
```