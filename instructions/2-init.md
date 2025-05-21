- implemented: true

---

# タスク

`github.com/mark3labs/mcp-go` を使って、J-Link Commander を制御する最小構成の MCP サーバを Go で実装してください。

## 前提

- このサーバは Raspberry Pi 上で動作し、USB接続された J-Link プローブを制御します
- `mcp-go` のプロトコルに準拠して、MCP クライアント（PC上）からの制御要求を受け取る
- J-Link Commander (`JLinkExe`) を subprocess 経由で操作します

## 実装内容

### 1. MCP サーバ構成

- `mcp-go` の `mcp.Server` を使ってサーバを起動
- 標準の HTTP サーバではなく、`mcp-go` の transport に準拠した gRPC サーバであること
- 対応する最小限のカスタムコマンドを1つだけ実装してください：

---

### 2. コマンド: `jlink.reset`

- 名前：`jlink.reset`
- 引数：なし
- 処理内容：
  1. 以下のJ-Link Commanderスクリプトを一時ファイル `cmds.jlink` に出力：
     ```
     device NRF52
     speed 4000
     r
     g
     q
     ```
  2. Goの `os/exec` を使って以下を実行：
     ```sh
     JLinkExe -CommanderScript cmds.jlink
     ```
  3. 実行後の標準出力を読み取り、ログとして返す
  4. 結果を `mcp.CommandResult` に `stdout` として格納してレスポンスする

---

## 補足

- エラー処理は簡易で構いません（`CommandResult.Error`にエラーメッセージを設定）
- `mcp.RegisterCommandHandler()` を使って `jlink.reset` を登録
- サーバは `mcp.Server.Start()` を使って起動
- `main.go` として書いてください

## 目的

この構成により、mcpクライアントから `jlink.reset` を呼び出すだけで、Raspi上のJ-Link経由でターゲットMCUをリセット・実行できます。
