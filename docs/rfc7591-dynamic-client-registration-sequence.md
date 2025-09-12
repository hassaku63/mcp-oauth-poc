# RFC 7591 Dynamic Client Registration — Sequence

```mermaid
sequenceDiagram
    autonumber
    participant Dev as Client
    participant AS as Authorization Server

    Dev->>AS: POST /register (JSON)

    alt Protected registration
        Dev->>AS: Authorization: Bearer <initial_access_token>
    else Open registration
        Dev->>AS: No Authorization header
    end

    alt Success
        AS-->>Dev: 201 Created (client + management)
    else Invalid metadata
        AS-->>Dev: 400 invalid_client_metadata
    end

    opt Management (RFC 7592)
        Dev->>AS: GET registration_client_uri (with RAT)
        AS-->>Dev: 200 OK (current metadata)

        Dev->>AS: PUT registration_client_uri (with RAT + updates)
        AS-->>Dev: 200 OK (updated metadata)

        Dev->>AS: DELETE registration_client_uri (with RAT)
        AS-->>Dev: 204 No Content
    end
```

- 注記:
  - 成功レスポンスはサーバポリシーにより `client_secret` を含まない場合があります（公開クライアントや非シークレット認証方式）。
  - `registration_access_token` と `registration_client_uri` によるクライアント構成 API は RFC 7592 に定義されています。
