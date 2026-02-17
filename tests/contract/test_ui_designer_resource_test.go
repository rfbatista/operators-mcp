package contract

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"testing"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"operators-mcp/internal/ui"
)

func TestReadDesignerResource_ProductionWithEmbed_Success(t *testing.T) {
	ctx := context.Background()
	// Use a minimal FS with index.html (handler reads "static/index.html").
	html := `<html><body>Designer</body></html>`
	testFS := &staticFS{files: map[string]string{"static/index.html": html}}

	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	server.AddResource(&sdkmcp.Resource{URI: ui.DesignerURI, Name: "Designer", MIMEType: "text/html"},
		ui.NewDesignerResourceHandler(false, testFS))

	t1, t2 := sdkmcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "client", Version: "0.0.1"}, nil)
	session, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer session.Close()

	res, err := session.ReadResource(ctx, &sdkmcp.ReadResourceParams{URI: ui.DesignerURI})
	if err != nil {
		t.Fatalf("ReadResource: %v", err)
	}
	if len(res.Contents) == 0 {
		t.Fatal("expected at least one content")
	}
	c := res.Contents[0]
	if c.URI != ui.DesignerURI {
		t.Errorf("URI = %q, want %q", c.URI, ui.DesignerURI)
	}
	if c.MIMEType != "text/html" {
		t.Errorf("MIMEType = %q, want text/html", c.MIMEType)
	}
	if c.Text != html {
		t.Errorf("Text = %q, want %q", c.Text, html)
	}
}

func TestReadDesignerResource_ProductionAssetsMissing_StructuredError(t *testing.T) {
	ctx := context.Background()
	// Nil FS or empty FS -> assets missing.
	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	server.AddResource(&sdkmcp.Resource{URI: ui.DesignerURI, Name: "Designer", MIMEType: "text/html"},
		ui.NewDesignerResourceHandler(false, nil))

	t1, t2 := sdkmcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "client", Version: "0.0.1"}, nil)
	session, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer session.Close()

	_, err = session.ReadResource(ctx, &sdkmcp.ReadResourceParams{URI: ui.DesignerURI})
	if err == nil {
		t.Fatal("expected error when assets missing")
	}
	// Should be a structured error the client can display (not a crash).
	if err.Error() == "" {
		t.Error("error should have a message")
	}
}

func TestReadDesignerResource_DevModeWithServerRunning_Success(t *testing.T) {
	// Start a minimal HTTP server that serves HTML.
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html><html><body>Vite dev</body></html>`))
	})
	srv := &http.Server{Addr: ":5174", Handler: mux}
	go srv.ListenAndServe()
	t.Cleanup(func() { srv.Close() })
	time.Sleep(50 * time.Millisecond)

	ctx := context.Background()
	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	server.AddResource(&sdkmcp.Resource{URI: ui.DesignerURI, Name: "Designer", MIMEType: "text/html"},
		ui.NewDesignerProxyHandler("http://localhost:5174"))

	t1, t2 := sdkmcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "client", Version: "0.0.1"}, nil)
	session, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer session.Close()

	res, err := session.ReadResource(ctx, &sdkmcp.ReadResourceParams{URI: ui.DesignerURI})
	if err != nil {
		t.Fatalf("ReadResource: %v", err)
	}
	if len(res.Contents) == 0 || res.Contents[0].Text == "" {
		t.Fatal("expected HTML content from proxy")
	}
	if !strings.Contains(res.Contents[0].Text, "Vite dev") {
		t.Errorf("expected proxied content, got %q", res.Contents[0].Text[:min(80, len(res.Contents[0].Text))])
	}
}

func TestReadDesignerResource_DevModeServerNotRunning_StructuredError(t *testing.T) {
	ctx := context.Background()
	// Use a port that is not listening.
	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	server.AddResource(&sdkmcp.Resource{URI: ui.DesignerURI, Name: "Designer", MIMEType: "text/html"},
		ui.NewDesignerProxyHandler("http://127.0.0.1:59999"))

	t1, t2 := sdkmcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "client", Version: "0.0.1"}, nil)
	session, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer session.Close()

	_, err = session.ReadResource(ctx, &sdkmcp.ReadResourceParams{URI: ui.DesignerURI})
	if err == nil {
		t.Fatal("expected error when dev server not running")
	}
	if !strings.Contains(err.Error(), "unreachable") && !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("expected DEV_SERVER_UNREACHABLE-style message, got: %v", err)
	}
}

// staticFS is a minimal fs.FS for tests.
type staticFS struct {
	files map[string]string
}

func (s *staticFS) Open(name string) (fs.File, error) {
	body, ok := s.files[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return &staticFile{name: name, content: body}, nil
}

type staticFile struct {
	name    string
	content string
	offset  int64
}

func (f *staticFile) Stat() (fs.FileInfo, error) {
	return &staticInfo{name: f.name, size: int64(len(f.content))}, nil
}

func (f *staticFile) Read(b []byte) (int, error) {
	if f.offset >= int64(len(f.content)) {
		return 0, io.EOF
	}
	n := copy(b, f.content[f.offset:])
	f.offset += int64(n)
	return n, nil
}

func (f *staticFile) Close() error { return nil }

type staticInfo struct {
	name string
	size int64
}

func (i *staticInfo) Name() string       { return i.name }
func (i *staticInfo) Size() int64        { return i.size }
func (i *staticInfo) Mode() fs.FileMode  { return 0444 }
func (i *staticInfo) ModTime() time.Time { return time.Time{} }
func (i *staticInfo) IsDir() bool        { return false }
func (i *staticInfo) Sys() interface{}   { return nil }
