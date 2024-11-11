#!/bin/bash

# 요청할 URL
URL="http://119.205.217.166:25000/identity/v3/auth/tokens" # 여기에 실제 URL을 입력하세요.

# JSON 본문
BODY='{
    "auth": {
        "identity": {
            "methods": [
                "password"
            ],
            "password": {
                "user": {
                    "name": "admin",
                    "domain": {
                        "id": "default"
                    },
                    "password": "cloud1234"
                }
            }
        },
        "scope": {
            "system": {
                "all": true
            }
        }
    }
}'

# POST 요청을 보내고 응답 헤더에서 X-Subject-Token 값 출력
TOKEN=$(wget --header="Content-Type: application/json" \
             --post-data="$BODY" \
             --no-check-certificate \
             -S -O- "$URL" 2>&1 | grep -i "X-Subject-Token" | awk '{print $2}')

# X-Subject-Token 출력
echo "X-Subject-Token: $TOKEN"
export Tkn=$TOKEN

#Project token get Url
PURL="http://119.205.217.166:25000/v3/auth/tokens"

# Project Name 입력
ProjectName="admin"

# JSON 본문
PBODY='{
    "auth": {
        "identity": {
            "methods": [
                "password"
            ],
            "password": {
                "user": {
                    "name": "admin",
                    "domain": {
                        "id": "default"
                    },
                    "password": "cloud1234"
                }
            }
        },
        "scope": {
            "project": {
                "domain": {
                    "name": "default"
                },
                "name": "$ProjectName"
            }
        }
    }
}'




PRTOKEN=$(wget --header="Content-Type: application/json" \
             --header="X-Auth-Token: $TOKEN" \
             --post-data="$BODY" \
             --no-check-certificate \
             -S -O- "$URL" 2>&1 | grep -i "X-Subject-Token" | awk '{print $2}')
#ProjectToken 출력
echo "X-Subject-Token(Project): $PRTOKEN"

export ProjectTkn=$PRTOKEN