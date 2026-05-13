#!/bin/bash
# Database backup script for Webhook Service
# Usage: ./scripts/backup.sh
# Cron: 0 2 * * * /app/scripts/backup.sh >> /var/log/backup.log 2>&1

set -e

BACKUP_DIR="/backup"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/webhook_service_${TIMESTAMP}.sql.gz"
RETENTION_DAYS=${BACKUP_RETENTION_DAYS:-7}

# Create backup
echo "[$(date)] Starting backup..."
pg_dump -U "${POSTGRES_USER:-webhook}" -d "${POSTGRES_DB:-webhook_service}" | gzip > "$BACKUP_FILE"

SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
echo "[$(date)] Backup created: $BACKUP_FILE ($SIZE)"

# Cleanup old backups
echo "[$(date)] Cleaning backups older than ${RETENTION_DAYS} days..."
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +${RETENTION_DAYS} -delete

REMAINING=$(find "$BACKUP_DIR" -name "*.sql.gz" | wc -l)
echo "[$(date)] Cleanup done. ${REMAINING} backup(s) remaining."

echo "[$(date)] Backup complete."
