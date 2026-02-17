package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"operators-mcp/internal/adapter/in/httpapi"
	"operators-mcp/internal/adapter/in/mcp"
	"operators-mcp/internal/adapter/in/ui"
	"operators-mcp/internal/adapter/out/filesystem"
	"operators-mcp/internal/adapter/out/persistence/memory"
	"operators-mcp/internal/application/blueprint"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	devMode := flag.Bool("dev", false, "enable dev mode: proxy ui://designer to Vite dev server")
	httpUI := flag.Bool("http", false, "expose the React UI on an HTTP server")
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address for the UI (e.g. :8080)")
	flag.Parse()

	cfg := mcp.Config{DevMode: *devMode}
	server := mcp.NewServer(cfg)

	root, _ := os.Getwd()
	projectStore := memory.NewProjectStore()
	zoneStore := memory.NewStore()
	pathMatcher := filesystem.NewMatcher()
	treeLister := filesystem.NewLister()
	svc := blueprint.NewService(projectStore, zoneStore, pathMatcher, treeLister, root)
	mcp.RegisterTools(server, svc)

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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		cancel()
	}()

	if *httpUI {
		// HTTP mode: single server serving React UI at / and MCP streamable transport at /mcp.
		mcpHandler := sdkmcp.NewStreamableHTTPHandler(func(*http.Request) *sdkmcp.Server { return server }, nil)
		uiHandler, err := ui.SPAHandler(ui.Dist)
		if err != nil {
			log.Fatalf("UI handler: %v", err)
		}
		mux := http.NewServeMux()
		apiHandler := httpapi.NewHandler(svc)
		apiHandler.Mount(mux, "/api")
		mux.Handle("/mcp", mcpHandler)
		mux.Handle("/mcp/", mcpHandler)
		mux.Handle("/", uiHandler)

		srv := &http.Server{Addr: *httpAddr, Handler: mux}
		go func() {
			log.Printf("HTTP server listening on %s (UI at /, MCP at /mcp)", *httpAddr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTP server: %v", err)
			}
		}()
		go func() {
			<-ctx.Done()
			_ = srv.Shutdown(context.Background())
		}()
		<-ctx.Done()
		return
	}

	if err := server.Run(ctx, &sdkmcp.StdioTransport{}); err != nil && ctx.Err() == nil {
		log.Fatalf("server: %v", err)
	}
}
