
# DeepClone library for go
[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/github.com/demosphen200/go-deepclone)


This repository contains a set of functions to deep clone values of any type.


DeepClone supports all types of go language. When cloning, new instances of all objects are created, including slice (slice and values) and map (map, map values, not keys)

## Usage


```go
import (
	"github.com/demosphen200/go-deepclone"
)

type SomeStruct struct {
	A string
	B []int
}

src := SomeStruct{
    A: "123",
    B: []int{1, 2, 3},
}
dst := SomeStruct{}
err := deepclone.DeepClone(&src, &dst)
if err != nil {
    panic("cannot clone")
}
// <-- here dst is deep clone of src 


srcSlice := []SomeStruct{src}
var dstSlice []SomeStruct
err = deepclone.DeepClone(&srcSlice, &dstSlice)
if err != nil {
    panic("cannot clone")
}

// <-- here dstSlice is deep clone of srcSlice 

```

## How it works

DeepClone uses reflection.
Primitive types are copied by value, structures are cloned field by field.
For array, slice and map, a new instance is created and then each element from the old one is cloned and placed into the new one.
For chan - a new chan is created with the same buffer size.
Due to limitations of reflect, only exported fields are cloned in structures.
the rest remain with default values.

## Custom cloners

To clone structures containing non-exported fields, use custom cloning functions. For the type time.Time this is already registered.


Example of registering a custom clone function. Just replace time.Time with your data type.
```go
RegisterStructCloner(
    func(src *time.Time) (*time.Time, error) {
        var tm = *src
        return &tm, nil
    },
)
```

## Where to use

DeepClone is great for testing purposes,
quick start in new applications,
where it is necessary to quickly implement new functions,
as well as for applications where there are no high loads.

DeepClone should not be used in high-load applications,
due to its low performance compared to specialized cloning functions,
designed for specific data types.

## Repository

https://github.com/demosphen200/go-deepclone is the canonical Git repository.

