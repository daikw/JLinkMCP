package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
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

	// Start the server using stdio transport
	log.Println("Starting JLink MCP server...")
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Failed to start server: %v", err)
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
