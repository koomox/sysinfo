# sysinfo
golang sysinfo
# Code Example           
```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/koomox/sysinfo"
)

func main() {
	si, err := sysinfo.Get()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, err := json.MarshalIndent(si, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(data))
}
```