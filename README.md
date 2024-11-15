# APITestProgram
인터넷이 없는 환경에서도 편하게 api 테스트 요청을 처리하기 위한 프로그램

# HTTP 클라이언트

이 프로젝트는 사용자가 입력한 URL에 대해 GET 및 POST 요청을 수행하는 Go 언어로 작성된 간단한 HTTP 클라이언트입니다. 사용자는 추가 헤더를 입력할 수 있으며, 응답을 JSON 형식으로 출력합니다.

## 기능

- GET 및 POST 요청을 지원합니다.
- 사용자로부터 URL과 추가 헤더를 입력받습니다.
- POST 요청의 경우 JSON 파일을 읽어서 요청 본문으로 사용합니다.
- 미리 등록된 JSON 파일을 읽어서 요청 본문으로 사용합니다.
- 등록된 JSON 파일 내용을 읽어서 출력합니다.
- 응답 헤더와 본문을 출력합니다.
- 종료 신호를 처리하여 프로그램을 안전하게 종료합니다.

## 사용 방법

1. 이 저장소를 클론합니다.

   ```bash
   git clone git@github.com:swlee3306/APITestProgram.git
   cd yourrepository

2. Go 언어가 설치되어 있어야 합니다. Go 설치 가이드를 참조하세요.

3. 프로그램을 실행합니다.(bash)
  
  ```bash
   go run main.go
  ```
4. 요청할 방법을 선택합니다:

- (1) GET 요청

- (2) POST 요청

- (3) POST(등록된 body file)

- (4) json body 내용 보기

- (5) 종료

5. 선택한 요청 방법에 따라 다음을 수행합니다:

- GET 요청:

  - 요청할 URL을 입력합니다.

  - 추가 헤더가 필요하면 입력합니다.

- POST 요청:

  - 요청할 URL을 입력합니다.

  - JSON 파일의 경로를 입력합니다.

  - 추가 헤더가 필요하면 입력합니다.

  - 프로그램은 입력된 URL로 요청을 보내고, 응답을 출력합니다.

- POST(등록된 body file):

  - json 데이터 목록에서 사용할 body data 입력

  - 요청할 URL을 입력합니다.

  - 추가 헤더가 필요하면 입력합니다.

  - 프로그램은 입력된 URL로 요청을 보내고, 응답을 출력합니다.

- json body 내용 보기:

  - json 데이터 목록에서 출력할 body data 입력

  - 선택한 json 내용을 출력합니다.

6. 코드 설명

- ApiReqRun: 사용자가 요청 방법을 선택하고 해당 요청을 처리하는 루프를 실행합니다.

- HandleGetRequest: GET 요청을 처리하며, 입력된 URL과 헤더를 사용하여 요청을 보냅니다.

- HandlePostRequest: POST 요청을 처리하며, 입력된 URL과 JSON 파일을 사용하여 요청을 보냅니다.

- handlePostRequestGetBody: POST 요청을 처리하며, 미리 정의된 JSON BODY 값을 받아 요청을 보냅니다.

- GetJsonList: POST 요청을 처리하며, 미리 정의된 JSON BODY 값을 가지고 handlePostRequestGetBody 함수를 이용하여 요청을 보냅니다.

- ShowJsonFile: 사용자가 선택한 JSON BODY 값을 출력합니다.

- getHeaders: 사용자로부터 추가 헤더를 입력받아 맵 형태로 반환합니다.

- addHeaders: 요청에 추가 헤더를 설정합니다.

- printResponse: 응답 헤더와 본문을 출력하며, JSON 본문은 예쁘게 출력됩니다.

7. 주의사항

- HTTPS 요청 시 인증서 검증을 무시하도록 설정되어 있으므로, 신뢰할 수 있는 서버와 통신할 때만 사용해야 합니다.

- 프로덕션 환경에서는 적절한 인증서 검증을 구현하는 것이 중요합니다.

8. 라이센스

- 이 프로젝트는 MIT 라이센스 하에 배포됩니다.

## 별첨

- /APITestProgram/internal/apiRequest/jsonfiles 경로 안에 미리 등록할 json 형식에 body 를 저장해서 사용 가능합니다.
