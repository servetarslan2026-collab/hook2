# Webhook Service - Proje Takip Dosyası

**Başlangıç:** 2026-05-13
**Repo:** https://github.com/servetarslan02/depo
**Tech Stack:** Go (Fiber) + Svelte 5 + Tailwind + PostgreSQL + NATS + Redis

---

## Aşama 1: Kritik Düzeltmeler (Derlenme için zorunlu)

- [ ] `go.sum` dosyası oluştur (`go mod tidy`)
- [ ] Route parametre uyumsuzluğunu düzelt (router `:org_id` → handler `Params("orgId")`)
- [ ] `.gitignore` ekle (node_modules, binary, .env, vs.)

## Aşama 2: Eksik Frontend Sayfalar (12 sayfa)

- [ ] `/organizations` — Tüm organizasyonlar listesi
- [ ] `/organizations/[id]/settings` — Organizasyon ayarları
- [ ] `/organizations/[id]/members` — Üye yönetimi
- [ ] `/organizations/[id]/service-tokens` — Service token yönetimi
- [ ] `/app/[id]/event-types` — Event type yönetimi (listele, oluştur, sil)
- [ ] `/app/[id]/settings` — Uygulama ayarları
- [ ] `/app/deliveries/[id]` — Delivery detay sayfası (request/response body)
- [ ] `/app/events/[id]` — Event detay sayfası (payload, metadata)
- [ ] `/app/subscriptions/new` — Yeni subscription oluşturma formu
- [ ] `/settings` — Kullanıcı profil ve şifre değiştirme
- [ ] `/register` — Kayıt sayfası (login'den ayrı)
- [ ] `/forgot-password` — Şifre sıfırlama sayfası

## Aşama 3: Eksik Özellikler

- [ ] Email doğrulama akışı (register sonrası)
- [ ] Organizasyona davet sistemi (email ile üye çağırma)
- [ ] Webhook imza doğrulama yardımcısı (frontend'de signature test)
- [ ] Pagination — tüm list sayfalarına ekle (events, deliveries, subscriptions)
- [ ] API key ile rate limiting (Redis tabanlı, farklı limitler)

## Aşama 4: İyileştirmeler

- [ ] OpenAPI/Swagger spec — tam ve doğru spec oluştur
- [ ] Test dosyaları — unit test (handler, store, worker)
- [ ] CI/CD pipeline — GitHub Actions (test, lint, build)
- [ ] Monitoring — Prometheus metrics endpoint (`/metrics`)
- [ ] Structured logging — tüm katmanlarda zap ile tutarlı log

## Aşama 5: Deploy & Production

- [ ] Production Docker Compose (Caddy reverse proxy, SSL)
- [ ] Environment değişkenleri dokümantasyonu
- [ ] Database backup stratejisi
- [ ] Health check endpoint'leri (DB, NATS, Redis)
- [ ] Graceful shutdown (tüm bileşenler)

---

## Tamamlananlar

- [x] Backend API yapısı (Go + Fiber) — 10 handler, 4 middleware
- [x] PostgreSQL veritabanı şeması + migrasyonlar
- [x] NATS JetStream kuyruk sistemi
- [x] Webhook delivery worker (3 retry: 1s, 30s, 5min)
- [x] HMAC-SHA256 imzalama
- [x] Dead letter queue
- [x] JWT + API Key authentication
- [x] Redis rate limiting (temel)
- [x] Frontend layout (sidebar, header, dark mode)
- [x] Login sayfası
- [x] Dashboard sayfası
- [x] Events listesi sayfası
- [x] Subscriptions listesi sayfası
- [x] Deliveries listesi sayfası
- [x] Secrets yönetimi sayfası
- [x] Send event sayfası
- [x] Organization create sayfası
- [x] Organization detail sayfası
- [x] Application dashboard sayfası
- [x] API client (ky) + tüm TypeScript tipler
- [x] Auth store (JWT token yönetimi)
- [x] Docker Compose (PostgreSQL + Redis + NATS + API + Frontend)
- [x] README dokümantasyonu

---

## Notlar

- Servet'in Rust deneyimi yok → Go ile devam
- Ölçek zirve olacak → NATS JetStream seçildi (RabbitMQ yerine)
- Hook0 SSPL lisanslı → Tasarım ilham alındı, kod sıfırdan yazıldı
- GitHub token'ı revoke edilmeli (chat'te paylaşıldı)
