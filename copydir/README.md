# copydir
A tool for copying contents of a directory to a destination directory, including all files and sub-directories.

Download:
```
> go get github.com/yaliv/go-pkg/copydir
```
Usage:  
**mybackup.go**
```
package main

import (
	"fmt"
	"os"

	"github.com/yaliv/go-pkg/copydir"
)

func main() {
	err := copydir.Copy("sample data", "backup", false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

```
