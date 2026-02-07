# go-twitch [![GoDoc](https://godoc.org/github.com/adeithe/go-twitch?status.svg)](https://godoc.org/github.com/adeithe/go-twitch) [![Go Report Card](https://goreportcard.com/badge/github.com/adeithe/go-twitch?style=flat-square)](https://goreportcard.com/report/github.com/adeithe/go-twitch) [![Codecov](https://img.shields.io/codecov/c/github/adeithe/go-twitch?style=flat-square&logo=codecov&label=codecov)](https://codecov.io/github/adeithe/go-twitch)

The go-twitch library is a complete interface for Twitch services. It is designed to be easy to use and allows for easy integration into any project.

## API

The library provides access to the Twitch API in order to manage and retrieve information about Twitch resources, users, streams, and more.

### Documentation

The documentation for each API call includes a link to the related [Twitch Developer Documentation](https://dev.twitch.tv/docs/api/reference/), making it easy to find related information.

Where possible, required fields are enforced by the call initializers, making it easier to use the API correctly.

> [!IMPORTANT]
> In some cases, such as where at only one of multiple fields is required, enforcement may not be available. For this reason, it is recommended to refer to the documentation to ensure correct usage.

### Writing Tests

When writing tests for code that uses the Twitch API, it is often useful to mock the API responses. For this reason, the `apitest` package is included to facilitate easy mocking of Twitch API responses.

To use, provide the mocks HTTP client using the `api.WithHTTPClient` option. Failing to set the HTTP client will result in th API request being sent to Twitch.

### Examples

<details>
<summary>Fetch all current livestreams</summary>

```go
package main

import (
	"context"
	"fmt"

	"github.com/adeithe/go-twitch/api"
)

const (
	ClientID = "wbmytr93xzw8zbg0p1izqyzzc5mbiz"
	OAuthToken = "2gbdx6oar67tqtcmt49t3wpcgycthx"
)

func main() {
	ctx := context.Background()
	client := api.New(ClientID)

	var cursor string
	req := client.Streams.List().First(100)
	for {
		streams, err := req.After(cursor).Do(ctx, api.WithBearerToken(OAuthToken))
		if err != nil {
			panic(err)
		}

		if len(streams.Data) == 0 {
			break
		}

		for _, stream := range streams.Data {
			fmt.Printf("%s is streaming %s to %d viewers\n",
				stream.UserLogin,
				stream.GameName,
				stream.ViewerCount,
			)
		}
		cursor = streams.Pagination.Cursor
	}
}
```
</details>

## Go Version Support

In accordance with the [Go Project's Supported Versions Policy](https://go.dev/doc/devel/release#policy), each Go release is supported until there are two newer major releases. After this period, support for older versions may be dropped in future releases.

## Need Help?

If you find a bug, please open an issue on the [GitHub Issues Page](https://github.com/adeithe/go-twitch/issues).

For general questions or support, visit the go-twitch [Discord](https://discord.gg/TRr8Wbby89) channel.
