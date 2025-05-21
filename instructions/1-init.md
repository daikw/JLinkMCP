- implemented: false
- ignored: true

---

Goで簡易MCPサーバを実装してください。

## 要件

- 対象: Raspberry Pi 上で動作
- 接続されている J-Link Pro を制御するために、JLinkExe を使用する
- HTTPサーバを立てて、REST API を介してターゲットデバイスを操作できるようにする
- 最小限の2つのエンドポイントを実装してください：

### 1. `POST /reset`
- JLinkExe に対して以下のスクリプトを流し込む：
  ```
  device NRF52
  speed 4000
  r
  g
  q
  ```
- 成功したら `"Reset OK"` を返す

### 2. `POST /mem/read`
- JSON で `{ "address": "0x20000000", "length": 4 }` の形式のボディを受け取る
- 指定されたアドレスと長さを使って JLinkExe でメモリを読み出す
- 例えば以下のようなスクリプトを作成して実行：
  ```
  device NRF52
  speed 4000
  mem32 0x20000000 1
  q
  ```
- 出力結果をパースして JSON で返す（例: `{ "value": 305419896 }`）

## 制約

- subprocess には `os/exec` を使うこと
- 実行用の一時スクリプト（.jlink）は都度生成・削除する
- エラー処理は簡易で構わない（500返す程度）

## 補足

- サーバはポート 8080 で待ち受ける
- 最終的に PC 側から `curl` で呼び出せるようにする

この要件に沿って、Goのファイル `main.go` を作成してください。
