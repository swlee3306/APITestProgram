# APITestProgram

**고급 HTTP 클라이언트 및 API 테스트 도구**

인터넷이 없는 환경에서도 편하게 API 테스트 요청을 처리하기 위한 Go 언어 기반 HTTP 클라이언트입니다. 포트폴리오용으로 보안, 사용자 경험, 기능을 대폭 강화한 전문적인 API 테스트 도구입니다.

## 🚀 주요 기능

### 🔐 보안 강화
- **TLS 인증서 검증 옵션**: 기본값 검증 활성화, 환경별 설정 지원
- **보안 헤더 자동 추가**: X-Content-Type-Options, X-Frame-Options 등
- **안전한 요청/응답 로깅**: 민감한 정보 마스킹 처리

### 🎨 사용자 경험 개선
- **대화형 CLI 인터페이스**: 컬러 출력 및 진행 표시기
- **설정 파일 기반 기본값**: YAML/JSON 설정 파일 지원
- **히스토리 기능**: 이전 요청 재실행 및 관리
- **진행 표시기**: 요청 진행 상황 실시간 표시

### ⚡ 고급 기능
- **확장된 HTTP 메서드**: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- **요청 체인 기능**: 여러 API 요청을 순차적으로 실행
- **변수 지원**: 이전 응답 데이터를 다음 요청에서 사용
- **조건부 실행**: 응답 조건에 따른 흐름 제어
- **데이터 추출**: 응답에서 필요한 데이터 자동 추출

### 📊 응답 처리
- **응답 검증 및 파싱**: JSON, XML, HTML 자동 파싱
- **상태 코드 분석**: HTTP 상태 코드별 처리
- **에러 처리**: 네트워크 오류 및 HTTP 오류 처리
- **응답 포맷팅**: 가독성 높은 응답 출력

### 🔧 개발자 도구
- **환경 변수 지원**: 설정을 환경 변수로 관리
- **설정 파일 관리**: YAML/JSON 설정 파일 지원
- **보안 로깅 시스템**: 민감한 정보 마스킹 처리
- **히스토리 관리**: 이전 요청 재실행 및 관리
- **UI/UX 개선**: 컬러 출력 및 사용자 친화적 인터페이스
- **테스트 자동화**: API 테스트 자동화 지원

## 📦 설치 및 실행

### 1. 저장소 클론
```bash
git clone https://github.com/swlee3306/APITestProgram.git
cd APITestProgram
```

### 2. 의존성 설치
```bash
go mod tidy
```

### 3. 프로그램 실행
```bash
# 기본 실행
go run main.go

# 설정 파일과 함께 실행
go run main.go --config config.yaml

# 환경 변수와 함께 실행
export APITEST_CONFIG_FILE=config.yaml
go run main.go
```

## 🛠️ 사용 방법

### 기본 HTTP 요청

#### 1. GET 요청
```bash
# 기본 GET 요청
curl -X GET https://httpbin.org/get

# 헤더와 함께 GET 요청
curl -X GET https://httpbin.org/get \
  -H "Authorization: Bearer token" \
  -H "User-Agent: APITestProgram"
```

#### 2. POST 요청
```bash
# JSON 데이터와 함께 POST 요청
curl -X POST https://httpbin.org/post \
  -H "Content-Type: application/json" \
  -d '{"name": "John", "email": "john@example.com"}'
```

#### 3. PUT/DELETE 요청
```bash
# PUT 요청
curl -X PUT https://httpbin.org/put \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "name": "Updated Name"}'

# DELETE 요청
curl -X DELETE https://httpbin.org/delete/1
```

### 요청 체인 기능

#### 1. 체인 생성
```go
package main

import (
    "APITestProgram/chain"
    "fmt"
)

func main() {
    // 체인 매니저 생성
    cm := chain.NewChainManager()
    
    // 체인 생성
    chain := cm.CreateChain("User API Test", "사용자 등록 및 로그인 테스트")
    
    // 등록 스텝 추가
    registrationStep := chain.ChainStep{
        Name:   "사용자 등록",
        Method: "POST",
        URL:    "https://api.example.com/register",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: `{"username": "testuser", "email": "test@example.com"}`,
        Extract: []chain.ExtractRule{
            {
                Name:     "user_id",
                Source:   "body",
                Pattern:  "id",
                Required: true,
            },
        },
    }
    
    // 로그인 스텝 추가 (등록에서 추출한 user_id 사용)
    loginStep := chain.ChainStep{
        Name:   "사용자 로그인",
        Method: "POST",
        URL:    "https://api.example.com/login",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: `{"user_id": "{{user_id}}", "password": "password123"}`,
        Extract: []chain.ExtractRule{
            {
                Name:     "auth_token",
                Source:   "body",
                Pattern:  "token",
                Required: true,
            },
        },
    }
    
    // 스텝을 체인에 추가
    cm.AddStep(chain.ID, registrationStep)
    cm.AddStep(chain.ID, loginStep)
    
    // 체인 실행
    results, err := cm.ExecuteChain(chain.ID)
    if err != nil {
        fmt.Printf("체인 실행 실패: %v\n", err)
        return
    }
    
    // 결과 출력
    for i, result := range results {
        fmt.Printf("스텝 %d: %s\n", i+1, result.StepID)
        fmt.Printf("  상태: %d\n", result.StatusCode)
        fmt.Printf("  성공: %t\n", result.Success)
        if len(result.Extracted) > 0 {
            fmt.Println("  추출된 변수:")
            for k, v := range result.Extracted {
                fmt.Printf("    %s = %s\n", k, v)
            }
        }
    }
}
```

#### 2. 대화형 CLI 사용
```go
package main

import "APITestProgram/chain"

func main() {
    cli := chain.NewChainCLI()
    cli.Run()
}
```

### 설정 파일 사용

#### config.yaml
```yaml
# 기본 설정
defaults:
  timeout: 30s
  retries: 3
  delay: 1s

# 보안 설정
security:
  verify_tls: true
  auto_security_headers: true
  mask_sensitive_data: true

# 로깅 설정
logging:
  level: info
  format: json
  file: logs/apitest.log

# 환경별 설정
environments:
  development:
    base_url: "https://dev-api.example.com"
    headers:
      X-Environment: "development"
  
  production:
    base_url: "https://api.example.com"
    headers:
      X-Environment: "production"
```

#### .env 파일
```bash
# 데이터베이스 설정
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=apitest_db
DB_USERNAME=root
DB_PASSWORD=password

# API 설정
API_BASE_URL=https://api.example.com
API_TIMEOUT=30s
API_RETRIES=3
```

## 🏗️ 프로젝트 구조

```
APITestProgram/
├── main.go                 # 메인 실행 파일
├── bin/                    # 컴파일된 바이너리
│   ├── ApiTestProgram-linux
│   ├── ApiTestProgram-MacBook_M1
│   ├── ApiTestProgram-win.exe
│   └── GetProjectTkn-수정.sh
├── chain/                 # 요청 체인 기능
│   ├── request_chain.go   # 체인 관리
│   ├── chain_cli.go      # 대화형 CLI
│   └── chain_test.go     # 테스트
├── cli/                   # CLI 인터페이스
│   └── interactive.go    # 대화형 CLI 구현
├── config/                # 설정 관리
│   ├── certificates.go   # 인증서 설정
│   └── settings.go      # 설정 관리
├── history/               # 히스토리 기능
│   └── history.go        # 요청 히스토리 관리
├── http/                  # HTTP 클라이언트
│   └── methods.go        # HTTP 메서드 정의
├── internal/              # 내부 패키지
│   ├── apiRequest/       # API 요청 처리
│   │   ├── apiRequest.go
│   │   ├── common.go
│   │   ├── run.go
│   │   └── jsonfiles/     # JSON 파일 저장소
│   │       ├── projectTkn.json
│   │       └── token.json
├── logging/               # 보안 로깅
│   └── secure_logger.go  # 보안 로깅 시스템
├── security/              # 보안 기능
│   ├── tls.go            # TLS 설정
│   └── headers.go        # 보안 헤더
├── ui/                    # UI 및 컬러 출력
│   └── colors.go         # 컬러 출력 및 UI
├── go.mod                # Go 모듈 파일
├── go.sum                # 의존성 체크섬
└── README.md             # 프로젝트 문서
```

## 🔧 고급 사용법

### 1. 환경별 설정
```bash
# 개발 환경
export APITEST_ENV=development
go run main.go

# 프로덕션 환경
export APITEST_ENV=production
go run main.go
```

### 2. 커스텀 헤더 설정
```yaml
# config.yaml
default_headers:
  User-Agent: "APITestProgram/1.0"
  Accept: "application/json"
  X-API-Version: "v1"
```

### 3. 프록시 설정
```bash
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=https://proxy.example.com:8080
go run main.go
```

### 4. SSL 인증서 검증
```bash
# SSL 검증 활성화 (기본값)
export APITEST_VERIFY_TLS=true
go run main.go

# SSL 검증 비활성화 (개발용)
export APITEST_VERIFY_TLS=false
go run main.go
```

### 5. 히스토리 관리
```go
// 히스토리 기능 사용
import "APITestProgram/history"

// 히스토리 매니저 생성
hm := history.NewHistoryManager()

// 요청 히스토리 저장
hm.SaveRequest("POST", "https://api.example.com/login", headers, body)

// 히스토리 조회
requests := hm.GetHistory()
for _, req := range requests {
    fmt.Printf("Method: %s, URL: %s, Time: %s\n", 
        req.Method, req.URL, req.Timestamp)
}

// 특정 히스토리 재실행
hm.ReplayRequest("request_id_123")
```

### 6. 보안 로깅
```go
// 보안 로깅 시스템 사용
import "APITestProgram/logging"

// 보안 로거 생성
logger := logging.NewSecureLogger()

// 민감한 필드 설정
logger.SetSensitiveFields([]string{"password", "token", "secret"})

// 요청 로깅 (민감한 정보 자동 마스킹)
logger.LogRequest("POST", "https://api.example.com/login", 
    map[string]string{
        "Content-Type": "application/json",
        "Authorization": "Bearer secret_token_123", // 자동 마스킹됨
    }, 
    `{"username": "user", "password": "secret123"}`, // 패스워드 마스킹됨
)

// 응답 로깅
logger.LogResponse(200, responseHeaders, responseBody, time.Since(startTime))
```

## 🧪 테스트

### 단위 테스트 실행
```bash
# 모든 테스트 실행
go test ./...

# 특정 패키지 테스트
go test ./chain
go test ./http
go test ./security

# 테스트 커버리지 확인
go test -cover ./...
```

### 통합 테스트
```bash
# 체인 기능 테스트
go test ./chain -v

# HTTP 클라이언트 테스트
go test ./http -v

# 보안 기능 테스트
go test ./security -v
```

## 📊 성능 및 모니터링

### 로깅 설정
```yaml
logging:
  level: debug
  format: json
  output: file
  file: logs/apitest.log
  max_size: 100MB
  max_age: 7d
  max_backups: 3
```

### 메트릭 수집
```go
// 요청 메트릭 수집
metrics := &Metrics{
    TotalRequests:   0,
    SuccessRequests: 0,
    FailedRequests:  0,
    AverageLatency:  0,
}
```

## 🔒 보안 고려사항

### 1. TLS 설정
- 기본적으로 TLS 인증서 검증 활성화
- 개발 환경에서만 검증 비활성화 가능
- 커스텀 CA 인증서 지원

### 2. 데이터 보호
- 민감한 정보 자동 마스킹
- 로그에서 패스워드, 토큰 등 제거
- 안전한 설정 파일 저장

### 3. 입력 검증
- URL 형식 검증
- 헤더 값 검증
- JSON 데이터 검증

## 🚀 배포 및 운영

### Docker 사용
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o apitest main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/apitest .
COPY --from=builder /app/config.yaml .
CMD ["./apitest"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  apitest:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APITEST_ENV=production
    volumes:
      - ./config.yaml:/app/config.yaml
      - ./logs:/app/logs
```

## 🤝 기여하기

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 라이센스

이 프로젝트는 MIT 라이센스 하에 배포됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

## 🙏 감사의 말

- [Go](https://golang.org/) - 프로그래밍 언어
- [Gin](https://gin-gonic.com/) - HTTP 웹 프레임워크
- [Viper](https://github.com/spf13/viper) - 설정 관리
- [Cobra](https://github.com/spf13/cobra) - CLI 프레임워크

## 📞 지원 및 문의

- 이슈 리포트: [GitHub Issues](https://github.com/swlee3306/APITestProgram/issues)
- 이메일: swlee3306@gmail.com
- 문서: [Wiki](https://github.com/swlee3306/APITestProgram/wiki)

---

**APITestProgram** - 전문적인 API 테스트를 위한 강력한 도구 🚀