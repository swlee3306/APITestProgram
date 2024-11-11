# APITestProgram
인터넷이 없는 환경에서도 편하게 api 테스트 요청을 처리하기 위한 프로그램

# HTTP 클라이언트

이 프로젝트는 사용자가 입력한 URL에 대해 GET 및 POST 요청을 수행하는 Go 언어로 작성된 간단한 HTTP 클라이언트입니다. 사용자는 추가 헤더를 입력할 수 있으며, 응답을 JSON 형식으로 출력합니다.

## 기능

- GET 및 POST 요청을 지원합니다.
- 사용자로부터 URL과 추가 헤더를 입력받습니다.
- POST 요청의 경우 JSON 파일을 읽어서 요청 본문으로 사용합니다.
- 응답 헤더와 본문을 출력합니다.
- 종료 신호를 처리하여 프로그램을 안전하게 종료합니다.

## 사용 방법

1. 이 저장소를 클론합니다.

   ```bash
   git clone git@github.com:swlee3306/APITestProgram.git
   cd yourrepository
