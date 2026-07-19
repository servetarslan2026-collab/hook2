#!/usr/bin/env powershell
# Jarvis baslatici - YAML dosyasini dogru sekilde olusturur
$yaml = @"
# JARVIS - Windows 11 AI Asistan
# Ollama + qwen3.6:35b

models:
  jarvis:
    provider: ollama
    model: qwen3.6:35b
    max_tokens: 32768

agents:
  root:
    model: jarvis
    name: Jarvis
    description: Kisilisel AI asistan
    max_iterations: 100
    add_environment_info: true
    instruction: |
      Sen Jarvis'sin. Turkce konus, net ve kisa ol.
      Windows 11 PowerShell kullan, Linux komutlari KULLANMA.
      Dosya: Get-ChildItem / dir
      Okuma: Get-Content / type
      Arama: Select-String / findstr
      Kopyala: Copy-Item
      Sil: Remove-Item
      Dizin: New-Item -ItemType Directory -Force
      Process: Get-Process / Stop-Process
      Path separator: \
    toolsets:
      - type: filesystem
        ignore_vcs: true
      - type: shell
        timeout: 300
      - type: git
      - type: think
      - type: todo
      - type: memory
        path: ./jarvis_memory.db
      - type: scheduler
      - type: fetch
        timeout: 60
      - type: background_jobs
      - type: session_plan
      - type: user_prompt
"@
# UTF-8 without BOM
$utf8 = New-Object System.Text.UTF8Encoding $false
[System.IO.File]::WriteAllText("jarvis.yaml", $yaml, $utf8)
Write-Host "jarvis.yaml olusturuldu" -ForegroundColor Green
docker agent run jarvis.yaml --yolo
