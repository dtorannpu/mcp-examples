package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/go-viper/mapstructure/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type Output struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func main() {
	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)
	cs, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: "http://localhost:8080"}, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect")
	}
	defer func(cs *mcp.ClientSession) {
		err := cs.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to close client session")
		}
	}(cs)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "greet", Arguments: Input{Name: "山田"}})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to call tool")
	}
	if res.IsError {
		log.Fatal().Msg("error calling tool")
	}

	var output Output
	err = mapstructure.Decode(res.StructuredContent, &output)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to decode response")
	}
	fmt.Println(output.Greeting)
}
