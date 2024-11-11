package apirequest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// NOTE - 입력받은 데이터를 저장해서 반환 하는 함수
func scanData() (string, error) {

	var indata string

	Reader := bufio.NewReader(os.Stdin)

	for {
		char, err := Reader.ReadByte()

		if err != nil {
			return "", err
		}

		// Enter 키를 누르면 종료
		if char == '\n' {
			break
		}

		// 백스페이스 처리
		if char == 127 { // ASCII 코드 127은 백스페이스
			if len(indata) > 0 {
				indata = indata[:len(indata)-1] // 마지막 문자 삭제
				fmt.Print("\r" + indata + " ")  // 현재 입력을 출력
				fmt.Print("\r" + indata)        // 커서를 맨 뒤로 이동
				continue
			}
		}

		indata += string(char)
	}

	return indata, nil
}

//NOTE - 추가할 header  읽어오는 함수
func getHeaders() map[string]string {
	headers := make(map[string]string)
	fmt.Println("추가할 헤더를 입력하세요 (형식: Key:Value). 입력을 마치려면 빈 줄을 입력하세요.")
	for {
		var header string
		fmt.Print("헤더 입력: ")
		fmt.Scanln(&header)
		if header == "" {
			break
		}
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		} else {
			fmt.Println("잘못된 형식입니다. 'Key:Value' 형식으로 입력하세요.")
		}
	}
	return headers
}

//NOTE - 읽어온 header 요청에 추가 하는 함수
func addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

//NOTE - 응답 본문 및 헤더 출력
func printResponse(resp *http.Response) {
	// 응답 헤더 출력
	fmt.Println("응답 헤더:")
	for key, value := range resp.Header {
		fmt.Printf("%s: %s\n", key, strings.Join(value, ", "))
	}

	// 응답 본문 출력
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var prettyBody bytes.Buffer
	err = json.Indent(&prettyBody, body, "", "    ")
	if err != nil {
		fmt.Println("JSON 예쁘게 출력 오류:", err)
		fmt.Println("응답 본문:", string(body))
		return
	}
	fmt.Println("응답 본문:")
	fmt.Println(prettyBody.String())
}
