#!/bin/bash
# seed-admin.sh — Create admin user after first deployment
# Usage: ./scripts/seed-admin.sh

set -e

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-webhook}"
DB_PASS="${DB_PASS:-webhook_secret}"
DB_NAME="${DB_NAME:-webhook_service}"

EMAIL="servetarslan02@gmail.com"
PASSWORD="Alayci_165"
NAME="Servet Arslan"

# Hash password with Go (bcrypt)
HASH=$(docker run --rm golang:1.22-alpine sh -c "
  cat <<'GO' > /tmp/hash.go
package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash, _ := bcrypt.GenerateFromPassword([]byte(\"${PASSWORD}\"), bcrypt.DefaultCost)
	fmt.Println(string(hash))
}
GO
  go run /tmp/hash.go
" 2>/dev/null)

if [ -z "$HASH" ]; then
  echo "❌ Failed to hash password"
  exit 1
fi

# Insert user
PGPASSWORD="${DB_PASS}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" <<SQL
INSERT INTO users (id, email, password_hash, name, is_admin, email_verified, created_at, updated_at)
VALUES (
  gen_random_uuid(),
  '${EMAIL}',
  '${HASH}',
  '${NAME}',
  true,
  true,
  NOW(),
  NOW()
)
ON CONFLICT (email) DO UPDATE SET
  is_admin = true,
  email_verified = true,
  updated_at = NOW();
SQL

echo "✅ Admin account created/updated: ${EMAIL}"
