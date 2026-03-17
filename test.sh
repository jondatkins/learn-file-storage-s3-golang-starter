#!/usr/bin/env bash
EMAIL="admin@tubely.com"
PASSWORD="password"

FORM_DATA=$(jq -n --arg email "$EMAIL" --arg password "$PASSWORD" '{email: $email, password: $password}')
echo $FORM_DATA

RESPONSE=$(curl -s -X POST http://localhost:8091/api/login \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d "$FORM_DATA")

JWT_TOKEN=$(echo "$RESPONSE" | jq -r '.token')
echo $RESPONSE | jq
echo $JWT_TOKEN
# GET /api/videos
# Headers:
# Authorization: Bearer ${jwtToken}
curl -w "\n" --data '{ "body": "foo", "user_id": "bfc7ea05-9c51-4664-ad54-5218c8fa94ed" }' --header 'Content-Type: application/json' http://localhost:8091/api/videos

curl -s -X GET http://localhost:8091/api/videos \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json"
