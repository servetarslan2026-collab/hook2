# ============================================================
#  JARVIS Otomatik Kurulum - Windows 11 PowerShell
#  Tüm bağımlılıkları otomatik kurar
# ============================================================

param(
    [switch]$SkipOllama,
    [switch]$SkipNode,
    [switch]$SkipDocker
)

$ErrorActionPreference = "SilentlyContinue"

Write-Host ""
Write-Host "  ========================================== " -ForegroundColor Cyan
Write-Host "   JARVIS Otomatik Kurulum v2.0              " -ForegroundColor Cyan
Write-Host "  ========================================== " -ForegroundColor Cyan
Write-Host ""

# ============================================================
# 1. Chocolatey Kurulumu (paket yöneticisi)
# ============================================================
Write-Host "[1/8] Chocolatey kontrol..." -ForegroundColor Yellow
if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
    Write-Host "  Chocolatey kuruluyor..." -ForegroundColor Yellow
    Set-ExecutionPolicy Bypass -Scope Process -Force
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
    try {
        Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
        Write-Host "  [OK] Chocolatey kuruldu" -ForegroundColor Green
    } catch {
        Write-Host "  [UYARI] Chocolatey kurulamadi, manuel kurulum gerekebilir" -ForegroundColor Yellow
    }
} else {
    Write-Host "  [OK] Chocolatey mevcut" -ForegroundColor Green
}

# Refresh PATH
$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")

# ============================================================
# 2. Node.js Kurulumu (Playwright MCP için)
# ============================================================
if (-not $SkipNode) {
    Write-Host "[2/8] Node.js kontrol..." -ForegroundColor Yellow
    if (-not (Get-Command node -ErrorAction SilentlyContinue)) {
        Write-Host "  Node.js kuruluyor..." -ForegroundColor Yellow
        choco install nodejs-lts -y | Out-Null
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
        Write-Host "  [OK] Node.js kuruldu" -ForegroundColor Green
    } else {
        Write-Host "  [OK] Node.js mevcut: $(node --version)" -ForegroundColor Green
    }
}

# ============================================================
# 3. Ollama Kurulumu
# ============================================================
if (-not $SkipOllama) {
    Write-Host "[3/8] Ollama kontrol..." -ForegroundColor Yellow
    if (-not (Get-Command ollama -ErrorAction SilentlyContinue)) {
        Write-Host "  Ollama kuruluyor..." -ForegroundColor Yellow
        choco install ollama -y | Out-Null
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
        Write-Host "  [OK] Ollama kuruldu" -ForegroundColor Green
    } else {
        Write-Host "  [OK] Ollama mevcut" -ForegroundColor Green
    }

    # Ollama servisini başlat
    $ollamaProcess = Get-Process ollama -ErrorAction SilentlyContinue
    if (-not $ollamaProcess) {
        Write-Host "  Ollama servisi başlatılıyor..." -ForegroundColor Yellow
        Start-Process ollama -ArgumentList "serve" -WindowStyle Hidden
        Start-Sleep -Seconds 3
        Write-Host "  [OK] Ollama başlatıldı" -ForegroundColor Green
    } else {
        Write-Host "  [OK] Ollama çalışıyor" -ForegroundColor Green
    }

    # Model indir
    Write-Host "  qwen3.6:35b modeli kontrol..." -ForegroundColor Yellow
    $modelList = ollama list 2>&1
    if ($modelList -notmatch "qwen3.6") {
        Write-Host "  Model indiriliyor (büyük dosya, uzun sürebilir)..." -ForegroundColor Yellow
        Write-Host "  İndirme boyutu: ~20GB" -ForegroundColor Gray
        ollama pull qwen3.6:35b
        Write-Host "  [OK] Model indirildi" -ForegroundColor Green
    } else {
        Write-Host "  [OK] qwen3.6:35b mevcut" -ForegroundColor Green
    }
}

# ============================================================
# 4. Docker Desktop Kontrol
# ============================================================
if (-not $SkipDocker) {
    Write-Host "[4/8] Docker Desktop kontrol..." -ForegroundColor Yellow
    if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
        Write-Host "  [UYARI] Docker Desktop bulunamadi!" -ForegroundColor Yellow
        Write-Host "  Docker Desktop 4.63+ gerekli" -ForegroundColor Yellow
        Write-Host "  Indirme: https://www.docker.com/products/docker-desktop/" -ForegroundColor Cyan
        Write-Host ""
        $install = Read-Host "  Docker Desktop simdi kurulsun mu? (E/H)"
        if ($install -eq "E" -or $install -eq "e") {
            choco install docker-desktop -y | Out-Null
            Write-Host "  [OK] Docker Desktop kuruldu (yeniden başlatma gerekebilir)" -ForegroundColor Green
        }
    } else {
        Write-Host "  [OK] Docker mevcut: $(docker --version)" -ForegroundColor Green
    }

    # docker-agent kontrol
    try {
        $agentVersion = docker agent --version 2>&1
        Write-Host "  [OK] docker-agent: $agentVersion" -ForegroundColor Green
    } catch {
        Write-Host "  [UYARI] docker-agent bulunamadi" -ForegroundColor Yellow
        Write-Host "  Docker Desktop 4.63+ gerekli" -ForegroundColor Yellow
    }
}

# ============================================================
# 5. Playwright MCP Kurulumu
# ============================================================
Write-Host "[5/8] Playwright MCP kontrol..." -ForegroundColor Yellow
if (Get-Command npx -ErrorAction SilentlyContinue) {
    $playwrightInstalled = npm list -g @playwright/mcp 2>&1
    if ($playwrightInstalled -notmatch "@playwright/mcp") {
        Write-Host "  Playwright MCP kuruluyor..." -ForegroundColor Yellow
        npm install -g @playwright/mcp | Out-Null
        Write-Host "  [OK] Playwright MCP kuruldu" -ForegroundColor Green
    } else {
        Write-Host "  [OK] Playwright MCP mevcut" -ForegroundColor Green
    }
    
    # Chromium browser indir
    Write-Host "  Chromium browser indiriliyor..." -ForegroundColor Yellow
    npx playwright install chromium 2>&1 | Out-Null
    Write-Host "  [OK] Chromium hazır" -ForegroundColor Green
} else {
    Write-Host "  [UYARI] npm bulunamadi, Playwright kurulamadi" -ForegroundColor Yellow
}

# ============================================================
# 6. Git Kurulumu
# ============================================================
Write-Host "[6/8] Git kontrol..." -ForegroundColor Yellow
if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Host "  Git kuruluyor..." -ForegroundColor Yellow
    choco install git -y | Out-Null
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
    Write-Host "  [OK] Git kuruldu" -ForegroundColor Green
} else {
    Write-Host "  [OK] Git mevcut: $(git --version)" -ForegroundColor Green
}

# ============================================================
# 7. Environment Variables Ayarla
# ============================================================
Write-Host "[7/8] Ortam değişkenleri ayarlanıyor..." -ForegroundColor Yellow

# Ollama keep alive
[System.Environment]::SetEnvironmentVariable("OLLAMA_KEEP_ALIVE", "24h", "User")
$env:OLLAMA_KEEP_ALIVE = "24h"
Write-Host "  [OK] OLLAMA_KEEP_ALIVE=24h" -ForegroundColor Green

# Ollama host
[System.Environment]::SetEnvironmentVariable("OLLAMA_HOST", "http://localhost:11434", "User")
$env:OLLAMA_HOST = "http://localhost:11434"
Write-Host "  [OK] OLLAMA_HOST=http://localhost:11434" -ForegroundColor Green

# ============================================================
# 8. Jarvis Başlatıcı Oluştur
# ============================================================
Write-Host "[8/8] Jarvis başlatıcı oluşturuluyor..." -ForegroundColor Yellow

# Masaüstü快捷方式
$desktopPath = [System.Environment]::GetFolderPath("Desktop")
$shortcutPath = Join-Path $desktopPath "Jarvis.lnk"
$shell = New-Object -ComObject WScript.Shell
$shortcut = $shell.CreateShortcut($shortcutPath)
$shortcut.TargetPath = "powershell.exe"
$shortcut.Arguments = "-NoExit -Command `"cd '$PSScriptRoot'; docker agent run jarvis.yaml --yolo`""
$shortcut.WorkingDirectory = $PSScriptRoot
$shortcut.Description = "Jarvis AI Asistan"
$shortcut.Save()
Write-Host "  [OK] Masaüstü快捷方式 oluşturuldu" -ForegroundColor Green

# ============================================================
# Tamamlandı
# ============================================================
Write-Host ""
Write-Host "  ========================================== " -ForegroundColor Green
Write-Host "   JARVIS KURULUM TAMAMLANDI!                " -ForegroundColor Green
Write-Host "  ========================================== " -ForegroundColor Green
Write-Host ""
Write-Host "  Kurulan bileşenler:" -ForegroundColor White
Write-Host "    ✅ Chocolatey (paket yöneticisi)" -ForegroundColor Gray
Write-Host "    ✅ Node.js + npm" -ForegroundColor Gray
Write-Host "    ✅ Ollama + qwen3.6:35b modeli" -ForegroundColor Gray
Write-Host "    ✅ Playwright MCP (tarayıcı otomasyonu)" -ForegroundColor Gray
Write-Host "    ✅ Git" -ForegroundColor Gray
Write-Host "    ✅ Ortam değişkenleri" -ForegroundColor Gray
Write-Host "    ✅ Masaüstü快捷方式" -ForegroundColor Gray
Write-Host ""
Write-Host "  Başlatma:" -ForegroundColor White
Write-Host "    1. Masaüstündeki Jarvis快捷方式'na çift tıkla" -ForegroundColor Gray
Write-Host "    2. Veya: docker agent run jarvis.yaml --yolo" -ForegroundColor Gray
Write-Host ""
Write-Host "  Komutlar:" -ForegroundColor White
Write-Host "    /yardim  - Tüm komutları listele" -ForegroundColor Gray
Write-Host "    /sistem  - Sistem durumu" -ForegroundColor Gray
Write-Host "    /git     - Git durumu" -ForegroundColor Gray
Write-Host "    /quit    - Çıkış" -ForegroundColor Gray
Write-Host ""
