// ...existing code...
package main

import (
	"log"

	"github.com/strowk/foxy-contexts/pkg/app"
	"github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
	"github.com/strowk/foxy-contexts/pkg/stdio"

	// _ "github.com/strowk/foxy-contexts/streamablehttp"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// This example defines resource tool for MCP server
// , run it with:
// npx @modelcontextprotocol/inspector go run main.go
// , then in browser open http://localhost:6274
// , then click Connect
// , then click List Resources
// , then click hello-world

// --8<-- [start:resource]
/*
func NewGreatResource() fxctx.Resource {
    return fxctx.NewResource(
        mcp.Resource{
            Name:        "hello-world",
            Uri:         "hello-world://hello-world",
            MimeType:    utils.Ptr("application/json"),
            Description: utils.Ptr("Hello World Resource"),
            Annotations: &mcp.ResourceAnnotations{
                Audience: []mcp.Role{
                    mcp.RoleAssistant, mcp.RoleUser,
                },
            },
        },
        func(_ context.Context, uri string) (*mcp.ReadResourceResult, error) {
            return &mcp.ReadResourceResult{
                Contents: []interface{}{
                    mcp.TextResourceContents{
                        MimeType: utils.Ptr("application/json"),
                        Text:     `{"hello": "world"}`,
                        Uri:      uri,
                    },
                },
            }, nil
        },
    )
}
*/
func NewGreatResource() fxctx.Resource {
	mime := "application/json"
	desc := "Hello World Resource"

	return fxctx.NewResource(
		mcp.Resource{
			Name:        "hello-world",
			Uri:         "hello-world://hello-world",
			MimeType:    &mime,
			Description: &desc,
			Annotations: &mcp.ResourceAnnotations{
				Audience: []mcp.Role{
					mcp.RoleAssistant, mcp.RoleUser,
				},
			},
		},
		func(uri string) (*mcp.ReadResourceResult, error) {
			return &mcp.ReadResourceResult{
				Contents: []interface{}{
					mcp.TextResourceContents{
						MimeType: &mime,
						Text:     `{"hello": "world"}`,
						Uri:      uri,
					},
				},
			}, nil
		},
	)
}

// new: additional resources that describe the tools (inspector will list them)
func NewHelloResource() fxctx.Resource {
	mime := "application/json"
	desc := "Hello tool - returns greeting for provided name"

	return fxctx.NewResource(
		mcp.Resource{
			Name:        "tool-hello",
			Uri:         "tool-hello://tool-hello",
			MimeType:    &mime,
			Description: &desc,
			Annotations: &mcp.ResourceAnnotations{
				Audience: []mcp.Role{mcp.RoleAssistant, mcp.RoleUser},
			},
		},
		func(uri string) (*mcp.ReadResourceResult, error) {
			return &mcp.ReadResourceResult{
				Contents: []interface{}{
					mcp.TextResourceContents{
						MimeType: &mime,
						Text:     `{"name":"<string>"} -> {"text":"Hello, <name>!"}`,
						Uri:      uri,
					},
				},
			}, nil
		},
	)
}

func NewAddResource() fxctx.Resource {
	mime := "application/json"
	desc := "Add tool - returns sum of two numbers (a + b)"

	return fxctx.NewResource(
		mcp.Resource{
			Name:        "tool-add",
			Uri:         "tool-add://tool-add",
			MimeType:    &mime,
			Description: &desc,
			Annotations: &mcp.ResourceAnnotations{
				Audience: []mcp.Role{mcp.RoleAssistant, mcp.RoleUser},
			},
		},
		func(uri string) (*mcp.ReadResourceResult, error) {
			return &mcp.ReadResourceResult{
				Contents: []interface{}{
					mcp.TextResourceContents{
						MimeType: &mime,
						Text:     `{"a":<number>,"b":<number>} -> {"sum": <number>}`,
						Uri:      uri,
					},
				},
			}, nil
		},
	)
}

// ...existing code...

// --8<-- [end:resource]

// --8<-- [start:server]
func main() {
	listChanged := false
	subscribe := false

	err := app.
		NewBuilder().
		// adding the resources (including the new tool-describing resources)
		WithResource(NewGreatResource).
		WithResource(NewHelloResource).
		WithResource(NewAddResource).
		WithServerCapabilities(&mcp.ServerCapabilities{
			Resources: &mcp.ServerCapabilitiesResources{
				ListChanged: &listChanged,
				Subscribe:   &subscribe,
			},
		}).
		// setting up server
		WithName("my-mcp-server").
		WithVersion("0.0.1").
		WithTransport(stdio.NewTransport()).
		// WithTransport(sse.NewTransport()).
		// WithTransport(streamable_http.NewTransport(
		// 	streamable_http.Endpoint{
		// 		Hostname: "localhost",
		// 		Port:     8080,
		// 		Path:     "/mcp",
		// 	})),
		// Configuring fx logging to only show errors
		WithFxOptions(
			fx.Provide(func() *zap.Logger {
				cfg := zap.NewDevelopmentConfig()
				cfg.Level.SetLevel(zap.ErrorLevel)
				logger, _ := cfg.Build()
				return logger
			}),
			fx.Option(fx.WithLogger(
				func(logger *zap.Logger) fxevent.Logger {
					return &fxevent.ZapLogger{Logger: logger}
				},
			)),
		).Run()
	if err != nil {
		log.Fatal(err)
	}
}
