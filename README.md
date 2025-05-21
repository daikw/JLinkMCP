# JLinkMCP

## 概要

J-Link プローブを制御する MCP サーバです。

## 使い方

```sh
go build -o jlink-mcp main.go
```

```json
{
  "mcpServers": {
    "jlink-mcp": {
      "command": "$your_dir/jlink-mcp",
      "args": []
    }
  }
}
```
