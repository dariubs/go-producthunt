# go-producthunt

`go-producthunt` is a Golang package for interacting with the Product Hunt GraphQL API, offering features like retrieving daily posts, product details, topic-based listings, and top-ranked products by date with strong error handling.


## Installation

To incorporate `go-producthunt` into your project, use the following command:

```bash
go get github.com/dariubs/go-producthunt
```

Ensure your Go environment is properly configured to use modules.


## Usage

Below is a basic example demonstrating how to use the package:

```go
package main

import (
	"fmt"
	"log"

	"github.com/dariubs/go-producthunt"
)

func main() {
	// Initialize the ProductHunt client with your API key.
	apiKey := "YOUR_API_KEY_HERE"
	client := producthunt.ProductHunt{APIKey: apiKey}

	// Retrieve the latest daily posts.
	products, err := client.GetDaily()
	if err != nil {
		log.Fatalf("Error fetching daily posts: %v", err)
	}

	// Display retrieved products.
	for _, product := range products {
		fmt.Printf("ID: %s, Name: %s, Tagline: %s\n", product.ID, product.Name, product.Tagline)
	}
}
```

Replace `"YOUR_API_KEY_HERE"` with your actual Product Hunt API key. This initializes the client and retrieves daily posts.



## Contributing

Contributions are welcome! Submit issues and pull requests via the [GitHub repository](https://github.com/dariubs/go-producthunt). Please follow coding conventions and include relevant tests and documentation with your contributions.


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

