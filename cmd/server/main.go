package main

import (
	"context"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"

	"operators-mcp/internal/adapter/in/httpapi"
	"operators-mcp/internal/adapter/in/mcp"
	"operators-mcp/internal/adapter/in/ui"
	"operators-mcp/internal/adapter/out/filesystem"
	"operators-mcp/internal/adapter/out/persistence/sqlite"
	"operators-mcp/internal/application/blueprint"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	devMode := flag.Bool("dev", false, "proxy ui://designer to Vite dev server (run 'make web-dev' separately)")
	mcpAddr := flag.String("mcp.addr", ":8081", "MCP server listen address (IDE connects here)")
	httpAddr := flag.String("http.addr", ":8080", "HTTP server listen address (UI and API)")
	dbPath := flag.String("db", "data.db", "SQLite database path (e.g. data.db or :memory:)")
	flag.Parse()

	db, err := sqlite.Open(*dbPath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	root, _ := os.Getwd()
	projectStore := sqlite.NewProjectRepository(db)
	zoneStore := sqlite.NewZoneRepository(db)
	pathMatcher := filesystem.NewMatcher()
	treeLister := filesystem.NewLister()
	svc := blueprint.NewService(projectStore, zoneStore, pathMatcher, treeLister, root)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		cancel()
	}()

	go runMCPServer(ctx, *mcpAddr, svc, *devMode)
	runHTTPServer(ctx, *httpAddr, svc)
}

// runMCPServer runs the MCP server on its own port using mcp-go streamable HTTP transport.
func runMCPServer(ctx context.Context, addr string, svc *blueprint.Service, devMode bool) {
	s := server.NewMCPServer("operators-mcp", "0.0.1", server.WithToolCapabilities(true))
	mcp.RegisterTools(s, svc)

	designerResource := mcplib.NewResource(ui.DesignerURI, "Designer",
		mcplib.WithResourceDescription("Architecture Designer UI"),
		mcplib.WithMIMEType("text/html"),
	)
	s.AddResource(designerResource, designerResourceHandler(devMode))

	httpServer := server.NewStreamableHTTPServer(s)
	srv := &http.Server{Addr: addr, Handler: httpServer}
	go func() {
		log.Printf("MCP server listening on %s (use this URL in your IDE)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("MCP server: %v", err)
		}
	}()
	<-ctx.Done()
	_ = srv.Shutdown(context.Background())
}

func designerResourceHandler(devMode bool) func(context.Context, mcplib.ReadResourceRequest) ([]mcplib.ResourceContents, error) {
	var embedFS fs.FS
	if !devMode {
		embedFS = ui.Dist
	}
	devURL := ui.DefaultDevServerURL
	return func(ctx context.Context, req mcplib.ReadResourceRequest) ([]mcplib.ResourceContents, error) {
		if req.Params.URI != ui.DesignerURI {
			return nil, nil
		}
		html, mime, err := ui.DesignerContent(ctx, devMode, embedFS, devURL)
		if err != nil {
			return nil, err
		}
		return []mcplib.ResourceContents{
			mcplib.TextResourceContents{
				URI:      ui.DesignerURI,
				MIMEType: mime,
				Text:     html,
			},
		}, nil
	}
}

// runHTTPServer runs the HTTP server for UI and API only (no MCP).
func runHTTPServer(ctx context.Context, addr string, svc *blueprint.Service) {
	uiHandler, err := ui.SPAHandler(ui.Dist)
	if err != nil {
		log.Fatalf("UI handler: %v", err)
	}
	mux := http.NewServeMux()
	apiHandler := httpapi.NewHandler(svc)
	apiHandler.Mount(mux, "/api")
	mux.Handle("/", uiHandler)

	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		log.Printf("HTTP server listening on %s — UI: /, API: /api", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server: %v", err)
		}
	}()
	<-ctx.Done()
	_ = srv.Shutdown(context.Background())
}
