# Go In-Memory-Library (lite) #
training project

## Setup ##

```
go get github.com/arturmazanik/in-memory-cache/inmemorylibrary
```

## Example ## 

```
package main

import (
	"fmt"
	"time"

	inmemlib "github.com/arturmazanik/in-memory-cache/inmemorylibrary"
)

func main() {
	c := inmemlib.NewCache(5 * time.Second)
	c.Set("foo", "bar", 15)

	val, ok := c.Get("foo")
	fmt.Println(val, ok) // bar true

	c.Delete("foo")
	val, ok = c.Get("foo")
	fmt.Println(val, ok)
}
```