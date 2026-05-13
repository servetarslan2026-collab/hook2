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

### Öncelik 1 — Eksik Sayfalar (12 tane)

Bu sayfalar henüz yok. Eklenmesi lazım:

| Sayfa | Ne işe yarar | Dosya yolu |
|-------|-------------|-----------|
| Organizasyon listesi | Tüm organizasyonları göster | `routes/organizations/+page.svelte` |
| Org ayarları | İsim değiştirme, silme | `routes/organizations/[id]/settings/+page.svelte` |
| Üye yönetimi | Davet et, çıkar | `routes/organizations/[id]/members/+page.svelte` |
| Service tokens | API token yönetimi | `routes/organizations/[id]/service-tokens/+page.svelte` |
| Event types | Listele, oluştur, sil | `routes/app/[id]/event-types/+page.svelte` |
| Uygulama ayarları | İsim/değiştirme, silme | `routes/app/[id]/settings/+page.svelte` |
| Delivery detay | Request/response body göster | `routes/app/deliveries/[id]/+page.svelte` |
| Event detay | Payload ve metadata göster | `routes/app/events/[id]/+page.svelte` |
| Yeni subscription | Form ile oluşturma | `routes/app/subscriptions/new/+page.svelte` |
| Kullanıcı ayarları | Profil ve şifre değiştirme | `routes/settings/+page.svelte` |
| Kayıt sayfası | Ayrı kayıt sayfası | `routes/register/+page.svelte` |
| Şifre sıfırlama | Email ile sıfırlama | `routes/forgot-password/+page.svelte` |

### Öncelik 2 — Eksik Özellikler (5 tane)

- [ ] Email doğrulama (register sonrası email gönderme)
- [ ] Organizasyona davet (email ile üye çağırma)
- [ ] Webhook imza test aracı (frontend'de)
- [ ] Pagination (sayfalama — tüm list sayfalarında)
- [ ] API key bazlı rate limiting

### Öncelik 3 — İyileştirmeler (5 tane)

- [ ] OpenAPI/Swagger spec
- [ ] Unit test'ler
- [ ] CI/CD (GitHub Actions)
- [ ] Monitoring (Prometheus)
- [ ] Structured logging

### Öncelik 4 — Production Hazırlığı (5 tane)

- [ ] SSL/TLS (Caddy reverse proxy)
- [ ] Environment değişkenleri dokümantasyonu
- [ ] Database backup
- [ ] Health check endpoint'leri
- [ ] Graceful shutdown

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
- **Frontend'de 12 sayfa eksik** → çalışır ama bazı sayfalar 404 verir
- **go.sum dosyası yok** → Docker build sırasında otomatik oluşur (go mod tidy)
- **package-lock.json yok** → Docker build sırasında otomatik oluşur (npm install)
