# TOOLS.md - Jarvis Tool Notları

## Shell Tool
- Windows PowerShell kullan
- Varsayılan timeout: 300 saniye
- Destructif komutlar: onay iste (eğer safer: true)

## Memory Tool
- SQLite veritabanı: `./jarvis_memory.db`
- Kalıcı bilgi saklama
- Her oturumda erişilebilir

## Scheduler Tool
- Zamanlanmış görevler oluştur
- Cron formatı veya relative time
- Hatırlatıcılar için ideal

## Background Jobs
- Uzun süren işlemleri arka plana at
- Sunucu, watcher, build işlemleri
- recall: true ile otomatik hatırlatma

## Snapshot Tool
- Otomatik dosya yedekleme
- Değişiklikleri geri alma (/undo)
- Git tabanlı

## MCP Tools
- **DuckDuckGo** — Web araması
- **Context7** — Dokümantasyon araması
- defer: true ile ilk kullanımda yüklenir

## Fetch Tool
- Web'den bilgi çekme
- Timeout: 60 saniye
- HTML → Markdown dönüşümü
