- implemented: true

---

# タスク

- MCP ツールの名前は、全て `_` 区切りとする
- JLinkExe コマンドは全て NRF52 を利用する前提とする
- インターフェースは SWD とする
- 実行するスクリプトはこんな感じにする
    ```
    eoe 1
    connect
    r
    mem 0x10000060 8
    q
    ```
