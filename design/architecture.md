# Backend Architecture Design

## Technology Stack
- Go 1.21+ (Backend API)
- Temporal (Workflow Engine)
- GraphQL with gqlgen
- PostgreSQL (Primary Database)
- Redis (Caching)
- AWS S3 (Document Storage)

## Core Components

### 1. API Layer (GraphQL)
```go
// schema/schema.graphql
type Location {
  id: ID!
  name: String!
  address: String!
  tags: LocationTags
  qrCode: String!
}

type LocationTags {
  supervisor: Contact
  siteContact: Contact
  emergencyContact: EmergencyContact!
}

type Contact {
  name: String!
  phone: String!
  email: String
}

type EmergencyContact {
  name: String!
  phone: String!
}

type Contractor {
  id: ID!
  email: String!
  name: String!
  photo: String!
  mobile: String!
  company: String!
  emergencyContact: EmergencyContact!
  roles: [String!]!
}

type Visit {
  id: ID!
  locationId: ID!
  contractorId: ID!
  signInTime: Time!
  signOutTime: Time
  documentAcknowledged: Boolean!
}
```

### 2. Temporal Workflows

```go
// workflows/contractor_workflow.go
type ContractorWorkflow struct {
    SignInActivity  activities.SignInActivity
    DocumentActivity activities.DocumentActivity
}

func (w *ContractorWorkflow) Execute(ctx workflow.Context, input ContractorSignInInput) error {
    // Handle contractor sign-in process with retry and compensation logic
    signInCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
        StartToCloseTimeout: time.Minute * 5,
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval: time.Second,
            MaximumAttempts: 3,
        },
    })
    
    // Execute sign-in activity
    if err := workflow.ExecuteActivity(signInCtx, w.SignInActivity.Execute, input).Get(ctx, nil); err != nil {
        return err
    }
    
    // Execute document acknowledgment activity
    return workflow.ExecuteActivity(signInCtx, w.DocumentActivity.Execute, input).Get(ctx, nil)
}
```

### 3. Domain Models

```go
// models/location.go
type Location struct {
    ID        uuid.UUID     `json:"id"`
    Name      string        `json:"name"`
    Address   string        `json:"address"`
    Tags      LocationTags  `json:"tags"`
    QRCode    string        `json:"qrCode"`
    CreatedAt time.Time     `json:"createdAt"`
    UpdatedAt time.Time     `json:"updatedAt"`
}

// models/contractor.go
type Contractor struct {
    ID              uuid.UUID       `json:"id"`
    Email           string         `json:"email"`
    Name            string         `json:"name"`
    Photo           string         `json:"photo"`
    Mobile          string         `json:"mobile"`
    Company         string         `json:"company"`
    EmergencyContact EmergencyContact `json:"emergencyContact"`
    Roles           []string       `json:"roles"`
    CreatedAt       time.Time      `json:"createdAt"`
    UpdatedAt       time.Time      `json:"updatedAt"`
}
```

## Security Implementation

1. Authentication:
```go
// middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        
        // Validate JWT token
        claims, err := validateToken(token)
        if err != nil {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        
        c.Set("user", claims)
        c.Next()
    }
}
```

2. Rate Limiting:
```go
// middleware/ratelimit.go
func RateLimitMiddleware(store *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        key := fmt.Sprintf("ratelimit:%s", ip)
        
        count, err := store.Incr(key).Result()
        if err != nil {
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
        
        store.Expire(key, time.Hour)
        
        if count > 1000 {
            c.AbortWithStatus(http.StatusTooManyRequests)
            return
        }
        
        c.Next()
    }
}
```

## Scalability Features

1. Database Indexing:
```sql
-- migrations/000001_create_indexes.up.sql
CREATE INDEX idx_visits_contractor_id ON visits(contractor_id);
CREATE INDEX idx_visits_location_id ON visits(location_id);
CREATE INDEX idx_visits_sign_in_time ON visits(sign_in_time);
```

2. Caching Strategy:
```go
// services/cache.go
type CacheService struct {
    redis *redis.Client
}

func (c *CacheService) GetLocation(id string) (*Location, error) {
    key := fmt.Sprintf("location:%s", id)
    
    // Try cache first
    if cached, err := c.redis.Get(key).Result(); err == nil {
        var location Location
        if err := json.Unmarshal([]byte(cached), &location); err == nil {
            return &location, nil
        }
    }
    
    // Cache miss, fetch from DB
    location, err := c.db.GetLocation(id)
    if err != nil {
        return nil, err
    }
    
    // Cache for 15 minutes
    cached, _ := json.Marshal(location)
    c.redis.Set(key, cached, 15*time.Minute)
    
    return location, nil
}
```

## Document Management

1. S3 Integration:
```go
// services/document.go
type DocumentService struct {
    s3Client *s3.Client
    bucket   string
}

func (s *DocumentService) StoreDocument(ctx context.Context, key string, content []byte) error {
    _, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
        Body:   bytes.NewReader(content),
    })
    return err
}
```

## Real-time Updates

1. WebSocket Handler:
```go
// handlers/websocket.go
func (h *WebSocketHandler) HandleContractorUpdates(c *websocket.Conn) {
    defer c.Close()
    
    // Subscribe to Redis pub/sub for real-time updates
    pubsub := h.redis.Subscribe(ctx, "contractor_updates")
    defer pubsub.Close()
    
    for {
        msg, err := pubsub.ReceiveMessage(ctx)
        if err != nil {
            break
        }
        
        if err := c.WriteJSON(msg.Payload); err != nil {
            break
        }
    }
}
```