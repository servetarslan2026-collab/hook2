# Webhook Service

**Ne bu?** Webhook gönderen bir platform. Diğer uygulamalar sana event gönderdiğinde, senin belirlediğin URL'lere webhook olarak iletiyor.

**Durum:** Çalışıyor ama eksik sayfalar var (aşağıda listelenmiş).

---

## 🚀 Çalıştırma (3 Adım)

```bash
# 1. Klonla
git clone https://github.com/servetarslan02/depo.git
cd depo

# 2. Başlat (ilk seferde 2-3 dakika sürer)
docker-compose up --build

# 3. Aç
# Tarayıcıda: http://localhost:5173
```

**Ön koşul:** [Docker](https://docs.docker.com/get-docker/) kurulu olmalı.

**İlk kullanım:**
1. Kayıt ol → Giriş yap
2. Organization oluştur → Application oluştur
3. Event Type tanımla → Subscription oluştur (webhook URL'in ile)
4. Event gönder → Delivery loglarını kontrol et

---

## 📁 Dosya Yapısı

```
webhook-service/
├── backend/                    → Go API sunucusu
│   ├── cmd/server/main.go      → Ana dosya (buradan başlar)
│   ├── internal/
│   │   ├── api/handlers/       → Her endpoint'in mantığı (10 dosya)
│   │   ├── api/middleware/     → Auth, CORS, logging, rate limit
│   │   ├── auth/               → JWT ve API key işlemleri
│   │   ├── models/             → Veri yapıları (User, Event, Subscription vs.)
│   │   ├── queue/              → NATS bağlantısı (webhook kuyruğu)
│   │   ├── store/              → PostgreSQL sorguları
│   │   └── worker/             → Webhook gönderme işçisi
│   └── Dockerfile
├── frontend/                   → Svelte 5 web arayüzü
│   ├── src/routes/             → Sayfalar (11 tane var)
│   ├── src/lib/api/client.ts   → Backend ile iletişim
│   └── Dockerfile
├── docker-compose.yaml         → Tüm servisler (DB, Redis, NATS, API, Frontend)
├── README.md                   → Detaylı dokümantasyon
└── PROGRESS.md                 → Bu dosya
```

---

## ✅ Yapılanlar (Çalışıyor)

| Özellik | Durum |
|---------|-------|
| Kullanıcı kayıt/giriş (JWT) | ✅ |
| Organizasyon oluşturma | ✅ |
| Uygulama oluşturma | ✅ |
| Event Type tanımlama | ✅ |
| Webhook Subscription (CRUD) | ✅ |
| Event gönderme (API + UI) | ✅ |
| Webhook delivery (3 deneme: 1sn, 30sn, 5dk) | ✅ |
| HMAC-SHA256 imzalama | ✅ |
| Dead letter queue (başarısız webhook'lar) | ✅ |
| Delivery logları | ✅ |
| API key yönetimi | ✅ |
| Rate limiting | ✅ |
| Dashboard istatistikleri | ✅ |
| Docker Compose ile tek komutla çalıştırma | ✅ |

---

## ❌ Eksikler (Yapılacaklar)

### ✅ Tüm öncelikler + gelecek iyileştirmeler tamamlandı!

Öncelik 1-4 tamamlandı. Gelecek iyileştirmeler de tamamlandı:

| İyileştirme | Durum | Detay |
|-------------|-------|-------|
| WebSocket real-time delivery updates | ✅ | ws/hub.go + ws/handler.go, NATS broadcast, canlı deliveries sayfası + toast bildirimler |
| Admin panel | ✅ | /admin dashboard, user management (make admin/delete), org management, dead letter queue |
| Retry strategy customization | ✅ | Subscription başına retry_count + retry_delays, worker per-subscription config kullanır |
| Multi-tenant isolation | ✅ | Per-org rate limiting (Redis), UUID validation, tüm handler'larda IsOrganizationMember kontrolü |
| Rate limit config | ✅ | TenantRateLimitMiddleware (500 req/min per org via Redis) |

### Öncelik 1 — Eksik Sayfalar (12 tane) ✅ TAMAMLANDI

| Sayfa | Ne işe yarar | Dosya yolu | Durum |
|-------|-------------|-----------|-------|
| Organizasyon listesi | Tüm organizasyonları göster | `routes/organizations/+page.svelte` | ✅ |
| Org ayarları | İsim değiştirme, silme | `routes/organizations/[id]/settings/+page.svelte` | ✅ |
| Üye yönetimi | Davet et, çıkar | `routes/organizations/[id]/members/+page.svelte` | ✅ |
| Service tokens | API token yönetimi | `routes/organizations/[id]/service-tokens/+page.svelte` | ✅ |
| Event types | Listele, oluştur, sil | `routes/app/event-types/+page.svelte` | ✅ |
| Uygulama ayarları | İsim/değiştirme, silme | `routes/app/settings/+page.svelte` | ✅ |
| Delivery detay | Request/response body göster | `routes/app/deliveries/[id]/+page.svelte` | ✅ |
| Event detay | Payload ve metadata göster | `routes/app/events/[id]/+page.svelte` | ✅ |
| Yeni subscription | Form ile oluşturma | `routes/app/subscriptions/new/+page.svelte` | ✅ |
| Kullanıcı ayarları | Profil ve şifre değiştirme | `routes/settings/+page.svelte` | ✅ |
| Kayıt sayfası | Ayrı kayıt sayfası | `routes/register/+page.svelte` | ✅ |
| Şifre sıfırlama | Email ile sıfırlama | `routes/forgot-password/+page.svelte` | ✅ |

**Ek düzeltmeler:** `currentAppId` store eklendi, mevcut app sayfaları store'a taşındı, sidebar güncellendi.

### Öncelik 2 — Eksik Özellikler (5 tane) ✅ TAMAMLANDI

- [x] Email doğrulama (register sonrası email gönderme) — Backend: verification_tokens tablosu + send/verify endpoints, Frontend: /verify sayfası + dashboard banner
- [x] Organizasyona davet (email ile üye çağırma) — Backend: AcceptInvitation endpoint + Frontend: /invite/[token] sayfası
- [x] Webhook imza test aracı (frontend'de) — /app/signature-test sayfası (HMAC-SHA256 compute + verify)
- [x] Pagination (sayfalama — tüm list sayfalarında) — Subscriptions, Event Types, Secrets sayfalarına eklendi
- [x] API key bazlı rate limiting — Rate limit middleware API key'e göre güncellendi

### Öncelik 3 — İyileştirmeler (5 tane) ✅ TAMAMLANDI

- [x] OpenAPI/Swagger spec — Tam OpenAPI 3.0.3 spec (tüm endpoint'ler, schema'lar, auth)
- [x] Unit test'ler — auth_test.go (hash, token, API key), middleware_test.go (prometheus, CORS)
- [x] CI/CD (GitHub Actions) — .github/workflows/ci.yml (backend, frontend, docker build)
- [x] Monitoring (Prometheus) — /metrics endpoint, http_requests_total, webhook_deliveries_total, request duration histogram
- [x] Structured logging — zap logger zaten vardı, prometheus metrics ile desteklendi

### Öncelik 4 — Production Hazırlığı (5 tane) ✅ TAMAMLANDI

- [x] SSL/TLS (Caddy reverse proxy) — Caddyfile + docker-compose.prod.yaml
- [x] Environment değişkenleri dokümantasyonu — .env.example genişletildi (tüm değişkenler, production checklist)
- [x] Database backup — scripts/backup.sh (gzip, retention, cron-ready)
- [x] Health check endpoint'leri — /health zaten vardı (main.go)
- [x] Graceful shutdown — zaten vardı (signal handling in main.go)

---

## 🛠 Teknoloji Kararları

| Karar | Seçenek | Neden |
|-------|---------|-------|
| Backend dili | Go (Rust değil) | Rust öğrenmek 3-6 ay, Go 1-2 hafta |
| Web framework | Fiber | Go'da en hızlı, Express benzeri |
| Frontend | Svelte 5 | React/Vue'dan daha hızlı, daha az kod |
| Queue | NATS JetStream | RabbitMQ'dan daha hızlı, daha basit |
| DB | PostgreSQL | Güvenilir, JSON desteği var |
| Cache | Redis | Rate limit ve session için |

---

## ⚠️ Bilinmesi Gerekenler

- **Hook0'dan ilham alındı** ama kod sıfırdan yazıldı (SSPL lisansı nedeniyle)
- **GitHub token'ı** chat'te paylaşıldı → revoke edilmeli
- **Frontend'de 12 sayfa eksik** → ✅ Tamamlandı (2026-05-13)
- **go.sum dosyası yok** → Docker build sırasında otomatik oluşur (go mod tidy)
- **package-lock.json yok** → Docker build sırasında otomatik oluşur (npm install)
