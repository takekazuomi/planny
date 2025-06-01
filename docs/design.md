---
title: Design
---

全体構造

```txt
http.ServerMux -+-> grpc service -+-> methods -> http.Handler
                |   (resource)    |
                |                 +-> methods -> http.Handler
                |
                +-> grpc service -+-> methods -> http.Handler
                |   (resource)    |
                |                 +-> methods -> http.Handler
                .
                .
                .
```

## 課題

- model の操作
- エラーの扱い

modelの操作を、handler(=grpc method) から切り離すべきか？
