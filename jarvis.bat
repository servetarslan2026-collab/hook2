@echo off
chcp 65001 >nul 2>&1
title Jarvis - AI Asistan

echo.
echo  ==========================================
echo   JARVIS - Kisilisel AI Asistan
echo  ==========================================
echo.

REM Ollama kontrol
ollama list >nul 2>&1
if %errorlevel% neq 0 (
    echo [HATA] Ollama calismiyor!
    echo Once kurulum scriptini calistirin: setup.ps1
    pause
    exit /b 1
)

REM Model kontrol
ollama list | findstr "qwen3.6" >nul 2>&1
if %errorlevel% neq 0 (
    echo [UYARI] qwen3.6:35b modeli bulunamadi!
    echo Indiriliyor...
    ollama pull qwen3.6:35b
)

REM Ollama servis kontrol
curl -s http://localhost:11434/api/tags >nul 2>&1
if %errorlevel% neq 0 (
    echo Ollama servisi baslatiliyor...
    start /B ollama serve
    timeout /t 5 /nobreak >nul
)

echo.
echo  Jarvis baslatiliyor...
echo  Cikis: Ctrl+C veya /quit
echo.

docker agent run jarvis.yaml --yolo

if %errorlevel% neq 0 (
    echo.
    echo [HATA] Jarvis baslatilamadi!
    echo.
    echo Cozumler:
    echo   1. Docker Desktop calisiyor mu?
    echo   2. setup.ps1 calistirdiniz mi?
    echo.
    pause
)
