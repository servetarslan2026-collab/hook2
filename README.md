# Webhook Service

Webhook-as-a-Service platformu — Go + Svelte ile sıfırdan inşa edildi.

## 🚀 Hızlı Başlangıç

### Ön Koşullar
- [Docker](https://docs.docker.com/get-docker/) kurulu olmalı
- [Docker Compose](https://docs.docker.com/compose/install/) kurulu olmalı

### Çalıştırma

```bash
# 1. Repoyu klonla
git clone https://github.com/servetarslan2026-collab/hook2.git
cd hook2

# 2. Tüm servisleri başlat (ilk seferde birkaç dakika sürer)
docker-compose up --build

# 3. Hazır! Tarayıcıda aç:
#    Frontend:  http://localhost:5173
#    API:       http://localhost:3000
#    NATS UI:   http://localhost:8222
```

### İlk Kullanım
1. `http://localhost:5173/register` → Hesap oluştur
2. Giriş yap → Organization oluştur → Application oluştur
3. Event Type tanımla → Subscription oluştur (webhook URL'in ile)
4. Event gönder → Delivery loglarını kontrol et

---

## 📐 Mimari

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Frontend   │────▶│  Backend    │────▶│ PostgreSQL  │
│ Svelte 5    │     │  Go/Fiber   │     │   (veri)    │
│ Tailwind    │     │             │     └─────────────┘
└─────────────┘     │             │
                    │             │────▶ Redis (cache)
                    │             │
                    │             │────▶ NATS (kuyruk)
                    └─────────────┘
                            │
                    ┌───────▼───────┐
                    │  Worker       │
                    │  (webhook     │
                    │   teslimat)   │
                    └───────────────┘
                            │
                    ┌───────▼───────┐
                    │  Hedef URL    │
                    └───────────────┘
```

## 🛠 Teknoloji Stack

| Bileşen | Teknoloji |
|---------|-----------|
| Backend | Go 1.22 + Fiber v2 |
| Frontend | Svelte 5 + SvelteKit + Tailwind CSS 4 |
| Veritabanı | PostgreSQL 17 |
| Cache | Redis 7 |
| Kuyruk | NATS JetStream |
| Auth | JWT + API Key |
| Deploy | Docker Compose |

---

## 📡 API Referansı

### Kimlik Doğrulama
```bash
# Kayıt ol
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@test.com", "password": "12345678", "name": "Test"}'

# Giriş yap
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@test.com", "password": "12345678"}'
```

### Webhook Gönderme
```bash
# API Key ile event gönder
curl -X POST http://localhost:3000/api/v1/applications/{app_id}/events \
  -H "Content-Type: application/json" \
  -H "X-API-Key: whkey_YOUR_KEY" \
  -d '{
    "event_type": "order.created",
    "payload": {"order_id": "12345", "amount": 99.99}
  }'
```

### Tüm Endpoint'ler

| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | `/api/v1/auth/register` | Kayıt ol |
| POST | `/api/v1/auth/login` | Giriş yap |
| GET | `/api/v1/organizations` | Organizasyonları listele |
| POST | `/api/v1/organizations` | Organizasyon oluştur |
| GET | `/api/v1/organizations/:id` | Organizasyon detayı |
| GET | `/api/v1/organizations/:id/applications` | Uygulamaları listele |
| POST | `/api/v1/organizations/:id/applications` | Uygulama oluştur |
| GET | `/api/v1/applications/:id` | Uygulama detayı |
| GET | `/api/v1/applications/:id/event-types` | Event type'ları listele |
| POST | `/api/v1/applications/:id/event-types` | Event type oluştur |
| GET | `/api/v1/applications/:id/subscriptions` | Subscription'ları listele |
| POST | `/api/v1/applications/:id/subscriptions` | Subscription oluştur |
| POST | `/api/v1/applications/:id/events` | Event gönder |
| GET | `/api/v1/applications/:id/events` | Event'leri listele |
| GET | `/api/v1/applications/:id/deliveries` | Delivery loglarını listele |
| POST | `/api/v1/deliveries/:id/retry` | Delivery'yi tekrar dene |
| GET | `/api/v1/applications/:id/secrets` | API key'leri listele |
| POST | `/api/v1/applications/:id/secrets` | API key oluştur |

---

## 🔐 Webhook İmzalama

Her webhook isteği `X-Webhook-Signature` header'ı ile imzalanır:

```
X-Webhook-Signature: sha256=<hex-encoded-hmac>
```

İmza doğrulama:
```python
import hmac, hashlib

def verify_signature(payload, secret, signature):
    expected = hmac.new(secret.encode(), payload, hashlib.sha256).hexdigest()
    return hmac.compare_digest(f"sha256={expected}", signature)
```

---

## ⚙️ Ortam Değişkenleri

| Değişken | Varsayılan | Açıklama |
|----------|-----------|----------|
| `SERVER_PORT` | 3000 | API port |
| `DATABASE_URL` | postgres://... | PostgreSQL bağlantısı |
| `NATS_URL` | nats://localhost:4222 | NATS bağlantısı |
| `REDIS_URL` | redis://localhost:6379 | Redis bağlantısı |
| `JWT_SECRET` | (değiştir!) | JWT imza sırrı |
| `JWT_EXPIRY` | 24h | Token süresi |
| `WORKER_COUNT` | 4 | Worker sayısı |
| `FRONTEND_URL` | http://localhost:5173 | CORS origin |

---

## 📁 Proje Yapısı

```
webhook-service/
├── backend/
│   ├── cmd/server/main.go          # Ana giriş noktası + migrasyon
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/           # HTTP handler'lar (10 dosya)
│   │   │   ├── middleware/         # Auth, CORS, logger, rate limit
│   │   │   └── router.go          # Route tanımları
│   │   ├── auth/                   # JWT + API key
│   │   ├── config/                 # Ortam değişkenleri
│   │   ├── models/                 # Veri modelleri
│   │   ├── queue/                  # NATS JetStream
│   │   ├── store/                  # PostgreSQL CRUD
│   │   └── worker/                 # Webhook delivery worker
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── lib/api/client.ts       # API istemcisi + tipler
│   │   ├── lib/stores/auth.ts      # Auth state
│   │   └── routes/                 # Sayfalar (11 sayfa)
│   └── Dockerfile
├── docker-compose.yaml             # Tüm servisler
├── PROGRESS.md                     # İlerleme takibi
└── README.md                       # Bu dosya
```

---

## 📝 Lisans

MIT
