package apirequest

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

//NOTE - get 요청
func HandleGetRequest() {
	fmt.Print("GET 요청할 URL 입력: ")

	url, err := scanData()
	if err != nil {
		fmt.Printf("데이터 읽기에 실패하였습니다.: %s", err)
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

//NOTE - post 요청
func HandlePostRequest() {
	fmt.Print("POST 요청할 URL 입력: ")

	url, err := scanData()
	if err != nil {
		fmt.Printf("데이터 읽기에 실패하였습니다.: %s", err)
		return
	}

	fmt.Print("JSON 파일 경로 입력: ")

	filePath, err := scanData()
	if err != nil {
		fmt.Printf("데이터 읽기에 실패하였습니다.: %s", err)
		return
	}

	body, err := ioutil.ReadFile(filePath)
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
