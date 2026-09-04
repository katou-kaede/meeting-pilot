# WebSocket

EchoではHTTP接続をWebSocketへUpgradeでき、gorilla/websocket が利用できます。

API通信が全部WebSocketに変わったわけではありません。

今後は2種類を併用します。

### HTTP

DBを更新・取得する通常のAPIです。

```
PATCH /api/meetings/:id/current-agenda
GET   /api/meetings/:id/session
```

- 議題を変更する
- Session情報を取得する
- 会議を完了する

### WebSocket

サーバーからブラウザへ、変更をリアルタイム通知する接続です。
```
ws://localhost:8080/ws/meetings/:id
```

流れはこうなります。
```
操作したブラウザ
  ↓ HTTP PATCH
GoサーバーがDB更新
  ↓ WebSocketで通知
同じ会議を開いている他のブラウザ
  ↓
Session情報を再取得して画面更新
```

AddClient() は、WebSocketで通知を送る相手を会議ごとに記録しています。

つまり、

- データ更新・取得：HTTP
- 更新されたことの通知：WebSocket

という役割分担です。

## 実装順
- WebSocket Hubを作る
- Session用WebSocketハンドラーを作る
- ルートを追加
- Reactから接続
- 議題切り替え成功後に配信
- 受信側でSessionデータを再取得

## 実装
- ライブラリ追加
```bash
go get github.com/gorilla/websocket
```
