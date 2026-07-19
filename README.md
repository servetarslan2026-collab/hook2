# 🤖 JARVIS v2.0 - Windows 11 Kişisel AI Asistan

Docker Agent + Ollama + qwen3.6:35b ile çalışan tam donanımlı kişisel AI asistanı.

## 🚀 Hızlı Kurulum

### Ön Gereksinimler

| Gereksinim | Sürüm | Link |
|---|---|---|
| Docker Desktop | 4.63+ | [docker.com](https://www.docker.com/products/docker-desktop/) |
| Ollama | Son sürüm | [ollama.ai](https://ollama.ai/) |
| Windows | 11 | - |
| RAM | 24GB+ (35B model için) | - |

### Kurulum Adımları

```powershell
# 1. Ollama'yı başlat
ollama serve

# 2. Modeli indir (~20GB)
ollama pull qwen3.6:35b

# 3. Jarvis'i başlat
docker agent run jarvis.yaml --request-timeout 0 --yolo

# Veya PowerShell script ile
.\setup.ps1

# Veya CMD ile
jarvis.bat
```

## 📁 Dosya Yapısı

```
├── jarvis.yaml                  # Ana agent konfigürasyonu (v2.0)
├── jarvis.bat                   # CMD başlatıcı
├── setup.ps1                    # PowerShell kurulum scripti
├── jarvis-timeout-patch.diff    # Stream timeout patch'i
├── AGENTS.md                    # Agent yönergeleri
├── TOOLS.md                     # Tool notları
└── README.md                    # Bu dosya
```

## ✨ v2.0 Yenilikleri

### 🔧 Yapılandırma
- **version: 12** — En son config schema
- **metadata** — Versiyon, yazar, açıklama
- **welcome_message** — Hoş geldin mesajı
- **add_date: true** — Tarih bilgisi otomatik eklenir

### 🛡️ Güvenlik & İzinler
- **permissions** — Fine-grained tool izin kontrolü
- **allow/deny/ask** — Hangi tool'ların otomatik çalışacağını belirle

### 🪝 Hooks (Otomasyon)
- **session_start** — Oturum başlangıcında sistem bilgisi toplar
- **turn_start** — Her tur başında git durumu + prompt files
- **turn_end** — Tarih bilgisi ekler
- **pre_tool_use / post_tool_use** — Snapshot (otomatik yedekleme)

### ⌨️ Komutlar (Slash Commands)
| Komut | Açıklama |
|---|---|
| `/yardim` | Tüm komutları listele |
| `/sistem` | CPU, RAM, Disk durumu |
| `/git` | Git durum özeti |
| `/proje` | Proje dizin analizi |
| `/temizle` | Geçici dosyaları temizle |
| `/gunluk` | Bugünün özeti |
| `/hava` | Hava durumu |
| `/hatirlat` | Hatırlatıcı kur |

### 🧠 Toolset'ler (Tam Liste)
| Tool | Amaç |
|---|---|
| filesystem | Dosya okuma/yazma/düzenleme |
| shell | PowerShell komutları (timeout: 5dk) |
| background_jobs | Arka plan işlemleri |
| git | Git versiyon kontrolü |
| think | Düşünme/planlama |
| todo | Yapılacaklar listesi |
| memory | Kalıcı bellek (SQLite) |
| scheduler | Zamanlayıcı |
| session_plan | Oturum planlama |
| sessioncontext | Oturum bağlamı |
| user_prompt | Kullanıcıdan input isteme |
| fetch | Web'den bilgi çekme |
| snapshot | Otomatik dosya yedekleme |
| mcp:duckduckgo | Web araması |
| mcp:context7 | Dokümantasyon araması |

### 🧩 Alt Agent
- **Araştırmacı** — Web araştırması ve dokümantasyon için optimize edilmiş

## 🎯 Kullanım Senaryoları

### Kodlama
```
> Python ile bir REST API yaz, FastAPI kullan
> Bu Node.js projesindeki hataları bul ve düzelt
> Docker Compose dosyası oluştur
```

### Sistem Yönetimi
```
> Çalışan servisleri listele
> Disk kullanımını göster
> Geçici dosyaları temizle
```

### Araştırma
```
> Docker best practices ara
> Bu kütüphanenin dokümantasyonunu bul
> Benzer projeler var mı araştır
```

### Zamanlanmış İşler
```
> Her gün sabah 9'da sistem durumu raporu ver
> 30 dakika sonra toplantı hatırlatması yap
> Haftalık git aktivite özeti çıkar
```

## ⚠️ Sorun Giderme

| Sorun | Çözüm |
|---|---|
| `Ollama calismiyor` | `ollama serve` |
| `Model bulunamadi` | `ollama pull qwen3.6:35b` |
| `5 dk timeout` | `--request-timeout 0` + timeout patch |
| `Yanlış shell komutları` | Otomatik PowerShell (instruction'da kurallar var) |
| `Docker bulunamadi` | Docker Desktop başlat |
| `RAM yetmiyor` | Daha küçük model kullan (qwen3.6:7b) |
| `Yavaş çalışıyor` | GPU sürücüsünü kontrol et |

## 🔧 Timeout Patch

Hardcoded 5 dakikalık stream timeout'u kaldırmak için:

```powershell
git clone https://github.com/docker/docker-agent.git
cd docker-agent
git apply ..\jarvis-timeout-patch.diff
task build
.\bin\docker-agent.exe run ..\jarvis.yaml --request-timeout 0
```

## 📝 Notlar

- `--request-timeout 0` = timeout yok
- `--yolo` = tüm tool_calls otomatik onayla
- `jarvis_memory.db` = kalıcı bellek (otomatik oluşur)
- `version: 12` = en son config schema
- Snapshot'lar otomatik oluşturulur (değişiklikleri geri almak için)
- Türkçe konuşur, Türkçe komut ver
