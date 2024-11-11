package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// 종료 신호를 처리하기 위한 채널
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("요청할 방법을 선택하세요: (1) GET (2) POST (3) 종료")

		var choice string

		for {
			char, err := reader.ReadByte()
			if err != nil {
				break
			}

			// Enter 키를 누르면 종료
			if char == '\n' {
				break
			}

			// 백스페이스 처리
			if char == 127 { // ASCII 코드 127은 백스페이스
				if len(choice) > 0 {
					choice = choice[:len(choice)-1] // 마지막 문자 삭제
					fmt.Print("\r" + choice + " ")  // 현재 입력을 출력
					fmt.Print("\r" + choice)        // 커서를 맨 뒤로 이동
					continue
				}
			}

			choice += string(char)
		}

		switch choice {
		case "1":
			handleGetRequest()
		case "2":
			handlePostRequest()
		case "3":
			fmt.Println("종료합니다...")
			return
		default:
			fmt.Println("잘못된 선택입니다. 다시 선택해 주세요.")
		}

		// 종료 신호를 체크
		select {
		case <-stopChan:
			fmt.Println("종료합니다...")
			return
		default:
			continue
		}
	}
}

func handleGetRequest() {
	var url string
	fmt.Print("GET 요청할 URL 입력: ")
	fmt.Scanln(&url)

	headers := getHeaders()

	// HTTP 클라이언트 생성 (인증서 무시 설정)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 인증서 검증 무시
		},
	}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	addHeaders(req, headers)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	printResponse(resp)
}

func handlePostRequest() {
	var url string
	fmt.Print("POST 요청할 URL 입력: ")
	fmt.Scanln(&url)

	fmt.Print("JSON 파일 경로 입력: ")
	var filePath string
	fmt.Scanln(&filePath)

	body, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("파일 읽기 오류:", err)
		return
	}

	headers := getHeaders()

	// HTTP 클라이언트 생성 (인증서 무시 설정)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 인증서 검증 무시
		},
	}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	addHeaders(req, headers)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	printResponse(resp)
}

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

func addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

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
