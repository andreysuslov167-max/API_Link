#!/bin/bash

BASE_URL="http://localhost:3000"

echo " 1. регистрация "
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}')
echo $REGISTER_RESPONSE
TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Токен: $TOKEN"
echo

echo " 2. создание ссылки"
LINK_RESPONSE=$(curl -s -X POST $BASE_URL/api/links \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"url":"https://example.com"}')
echo $LINK_RESPONSE
SHORT_CODE=$(echo $LINK_RESPONSE | grep -o '"short_code":"[^"]*"' | cut -d'"' -f4)
echo "Короткий код: $SHORT_CODE"
echo

echo " 3. получение ссылки"
curl -s -X GET $BASE_URL/api/links/1 \
  -H "Authorization: Bearer $TOKEN"
echo

echo " 4. статистика "
curl -s -X GET $BASE_URL/api/links/1/stats \
  -H "Authorization: Bearer $TOKEN"
echo

echo " 5. Редирект "
curl -s -L $BASE_URL/$SHORT_CODE
echo
