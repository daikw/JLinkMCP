package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Command line flags
	transport := flag.String("transport", "stdio", "Transport type: stdio or sse")
	addr := flag.String("addr", ":8080", "Address for SSE server to listen on (only used with sse transport)")
	baseURL := flag.String("baseurl", "http://localhost:8080", "Base URL for SSE server (only used with sse transport)")
	flag.Parse()

	// Create a new MCP server
	mcpServer := server.NewMCPServer(
		"jlink-mcp",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// Register the jlink.reset command handler
	jlinkResetTool := mcp.NewTool(
		"jlink_reset",
		mcp.WithDescription("Reset a device using J-Link Commander"),
	)

	mcpServer.AddTool(jlinkResetTool, handleJLinkReset)

	// Start the server with the selected transport
	log.Printf("Starting JLink MCP server using %s transport...", *transport)

	// Handle server based on transport type
	switch *transport {
	case "sse":
		// Create SSE server
		sseServer := server.NewSSEServer(
			mcpServer,
			server.WithBaseURL(*baseURL),
		)

		// Setup signal handling for graceful shutdown
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		// Start server in a goroutine
		go func() {
			log.Printf("SSE server listening on %s", *addr)
			log.Printf("SSE endpoint: %s/sse", *baseURL)
			if err := sseServer.Start(*addr); err != nil {
				log.Fatalf("Failed to start SSE server: %v", err)
			}
		}()

		// Wait for interrupt signal
		<-stop
		log.Println("Shutting down SSE server...")

		// Shutdown gracefully
		if err := sseServer.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error during server shutdown: %v", err)
		}
		log.Println("Server stopped")

	case "stdio":
		// Use standard stdio transport as before
		log.Println("Starting JLink MCP server...")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Failed to start stdio server: %v", err)
		}

	default:
		log.Fatalf("Unsupported transport type: %s", *transport)
	}
}

// handleJLinkReset implements the jlink.reset command
func handleJLinkReset(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create a temporary file for the J-Link Commander script
	tmpDir := os.TempDir()
	scriptPath := filepath.Join(tmpDir, "cmds.jlink")

	// Write the J-Link Commander script
	scriptContent := `
speed 4000
eoe 1
connect
r
mem 0x10000060 8
q
`
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create script file: %v", err)), nil
	}
	defer os.Remove(scriptPath) // Clean up the temporary file

	// Execute the J-Link Commander with the script
	cmd := exec.Command("JLinkExe", "-device", "nRF52", "-speed", "4000", "-if", "SWD", "-CommanderScript", scriptPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return mcp.NewToolResultError(
			fmt.Sprintf("Failed to execute JLinkExe: %v\n%s", err, string(output)),
		), nil
	}

	// Return the command result
	return mcp.NewToolResultText(string(output)), nil
}
