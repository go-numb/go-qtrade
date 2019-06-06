# go-qtrade

Qtrade API version1

## Description

go-qtrade is a go client library for [Qtrade API](https://qtrade-exchange.github.io/qtrade-docs).

## Installation

```
$ go get -u github.com/go-numb/go-qtrade
```

## Usage
``` 
package main

import (
 "fmt"
 "github.com/go-numb/go-qtrade"
)


func main() {
	c := qtrade.New("", "")
	res, err := c.Ticker(VEOBTC)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", res)
}
```

## Author

[@_numbP](https://twitter.com/_numbP)

## License

[MIT](https://github.com/go-numb/go-qtrade/blob/master/LICENSE)