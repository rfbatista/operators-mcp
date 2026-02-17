package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"operators-mcp/internal/blueprint"
	"operators-mcp/internal/mcp"
	"operators-mcp/internal/ui"
)

func main() {
	devMode := flag.Bool("dev", false, "enable dev mode: proxy ui://designer to Vite dev server")
	flag.Parse()

	cfg := mcp.Config{DevMode: *devMode}
	server := mcp.NewServer(cfg)

	root, _ := os.Getwd()
	store := blueprint.NewStore()
	blueprint.RegisterTools(server, root, store)

	if *devMode {
		server.AddResource(&sdkmcp.Resource{
			URI:      ui.DesignerURI,
			Name:     "Designer",
			MIMEType: "text/html",
		}, ui.NewDesignerProxyHandler(ui.DefaultDevServerURL))
	} else {
		server.AddResource(&sdkmcp.Resource{
			URI:      ui.DesignerURI,
			Name:     "Designer",
			MIMEType: "text/html",
		}, ui.NewDesignerResourceHandler(false, ui.Dist))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		cancel()
	}()

	if err := server.Run(ctx, &sdkmcp.StdioTransport{}); err != nil && ctx.Err() == nil {
		log.Fatalf("server: %v", err)
	}
}
