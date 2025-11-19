# アーキテクチャ

## システム全体像

```mermaid
graph TB
    Client[HTTP Client] -->|HTTP Request| Router[Gin Router]
    Router -->|Route| Handler[Handler Layer]
    Handler -->|Business Logic| Service[Service Layer]
    Service -->|Data Access| Repository[Repository Layer]
    Repository -->|Query| DB[(Database)]

    OpenAPI[OpenAPI Spec] -->|Code Generation| Generated[Generated API Code]
    Generated -->|Implements| Handler

    subgraph "Application Layer"
        Handler
        Service
    end

    subgraph "Infrastructure Layer"
        Repository
        DB
    end

    subgraph "Domain Layer"
        Domain[Domain Models]
    end

    Service -.->|Uses| Domain
    Repository -.->|Uses| Domain
```

## レイヤー構造と依存関係

```mermaid
graph LR
    subgraph "Presentation Layer"
        A[HTTP Handler<br/>handler/announcement.go]
    end

    subgraph "Application Layer"
        B[Service<br/>service/announcement.go]
    end

    subgraph "Domain Layer"
        C[Domain Model<br/>domain/announcement.go]
    end

    subgraph "Infrastructure Layer"
        D[Repository<br/>repository/annoucement_repository.go]
        E[Mock Repository<br/>repository/mock_announcement_repository.go]
    end

    subgraph "External"
        F[OpenAPI Spec<br/>openapi/openapi.yaml]
        G[Generated Code<br/>generated/api.gen.go]
        H[Gin Framework]
        I[GORM]
    end

    A -->|depends on| B
    A -->|uses| H
    A -->|implements| G
    G -->|generated from| F

    B -->|depends on| C
    B -->|depends on| D

    D -->|depends on| C
    D -->|uses| I
    E -->|depends on| C

    style C fill:#e1f5ff
    style B fill:#fff4e1
    style A fill:#ffe1f5
    style D fill:#e1ffe1
    style E fill:#e1ffe1
```

## リクエストフロー

```mermaid
sequenceDiagram
    participant Client
    participant Router as Gin Router
    participant Handler as Handler Layer
    participant Service as Service Layer
    participant Repository as Repository Layer
    participant DB as Database

    Client->>Router: GET /announcements
    Router->>Handler: AnnouncementsList()
    Handler->>Service: GetAnnouncements()
    Service->>Repository: GetAnnouncements()
    Repository->>DB: SELECT * FROM announcements
    DB-->>Repository: Result Set
    Repository-->>Service: []domain.Announcement
    Service-->>Handler: []domain.Announcement
    Handler->>Handler: Convert to JSON
    Handler-->>Router: HTTP 200 + JSON
    Router-->>Client: HTTP Response
```

## コンポーネント詳細

### Handler Layer (`internal/handler/`)

- **責務**: HTTP リクエストの受信とレスポンスの返却
- **依存**: Service Layer
- **実装**: Gin フレームワークを使用
- **主要コンポーネント**: `Handler`構造体

### Service Layer (`internal/service/`)

- **責務**: ビジネスロジックの実装
- **依存**: Domain Layer, Repository Interface
- **主要コンポーネント**: `AnnouncementService`

### Domain Layer (`internal/domain/`)

- **責務**: ドメインモデルの定義
- **依存**: なし（最下層）
- **主要コンポーネント**: `Announcement`構造体

### Repository Layer (`internal/repository/`)

- **責務**: データアクセスの実装
- **依存**: Domain Layer
- **実装**: GORM を使用
- **主要コンポーネント**:
  - `announcementRepository` (実装)
  - `MockAnnouncementRepository` (テスト用)

## 依存関係の方向

```mermaid
graph TD
    A[Domain Layer] -->|No Dependencies| A
    B[Service Layer] -->|Depends on| A
    C[Handler Layer] -->|Depends on| B
    D[Repository Layer] -->|Depends on| A

    style A fill:#e1f5ff
    style B fill:#fff4e1
    style C fill:#ffe1f5
    style D fill:#e1ffe1
```

**依存関係の原則**:

- 外側のレイヤーは内側のレイヤーに依存する
- 内側のレイヤーは外側のレイヤーに依存しない
- Domain Layer は最下層で、他のレイヤーに依存しない

## 技術スタック

- **Web Framework**: Gin
- **ORM**: GORM
- **API 仕様**: OpenAPI 3.1.0
- **コード生成**: oapi-codegen
- **言語**: Go
