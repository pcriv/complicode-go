# Complicode

Control code generator for invoices inside the Bolivian national tax service.

### Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/pablocrivella/complicode-go"
)

func main() {
	authCode := "29040011007"
	key := "9rCB7Sv4X29d)5k7N%3ab89p-3(5[A"
	date, _ := time.Parse("20060102", "20070702")
	invoice := complicode.Invoice{Number: 1503, Nit: 4189179011, Date: date, Amount: 2500}
	code := complicode.Generate(authCode, key, invoice)

	fmt.Printf(code) // => "6A-DC-53-05-14"
}
```

## License

Copyright 2020 [Pablo Crivella](https://pablocrivella.me).
Read [LICENSE](LICENSE) for details.