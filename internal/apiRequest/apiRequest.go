package apirequest

import (
	"bytes"
	"crypto/tls"
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

//go:embed jsonfiles/*.json
var jsonFiles embed.FS

// NOTE - get 요청
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

// NOTE - post 요청
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

// NOTE - 원하는 body file 선택 후 post 요청
func handlePostRequestGetBody(body []byte) {
	fmt.Print("POST 요청할 URL 입력: ")

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

// NOTE - body list 출력 및 post 요청
func GetJsonList() {

	// json 폴더 내의 모든 JSON 파일을 읽어옴
	files, err := jsonFiles.ReadDir("jsonfiles")

	if err != nil {
		log.Fatalf("error reading directory: %v", err)
	}

	var filenamelist string
	var filenumber []int

	for i, v := range files {
		filenumber = append(filenumber, i)
		filenamelist += "(" + strconv.Itoa(i+1) + ")" + " " + v.Name() + " "
	}

	fmt.Printf("post 로 요청할 body 입력: %s\n" + filenamelist)

	bodylist, err := scanData()

	if err != nil {
		fmt.Printf("데이터 읽기에 실패하였습니다.: %s", err)
		return
	}

	bodynum, _ := strconv.Atoi(bodylist)
	bodynum -= 1

	for i, v := range filenumber {
		if v == bodynum {
			content, err := jsonFiles.ReadFile("jsonfiles/" + files[i].Name())

			if err != nil {
				log.Fatalf("error reading file: %v", err)
			}
			handlePostRequestGetBody(content)
		}
	}

	for _, file := range files {
		if !file.IsDir() {
			content, err := jsonFiles.ReadFile("jsonfiles/" + file.Name())
			if err != nil {
				log.Fatalf("error reading file: %v", err)
			}
			fmt.Printf("Content of %s:\n%s\n", file.Name(), content)
		}
	}
}

// NOTE - body list 출력 및 내용 출력
func ShowJsonFile() {
	// json 폴더 내의 모든 JSON 파일을 읽어옴
	files, err := jsonFiles.ReadDir("jsonfiles")

	if err != nil {
		log.Fatalf("error reading directory: %v", err)
	}

	var filenamelist string
	var filenumber []int

	for i, v := range files {
		filenumber = append(filenumber, i)
		filenamelist += "(" + strconv.Itoa(i+1) + ")" + " " + v.Name() + " "
	}

	fmt.Printf("출력할 json file 선택: %s\n", filenamelist)

	jsonlist, err := scanData()

	if err != nil {
		fmt.Printf("데이터 읽기에 실패하였습니다.: %s", err)
		return
	}

	jsonnum, _ := strconv.Atoi(jsonlist)
	jsonnum -= 1

	for i, v := range filenumber {
		if v == jsonnum {
			content, err := jsonFiles.ReadFile("jsonfiles/" + files[i].Name())

			if err != nil {
				log.Fatalf("error reading file: %v", err)
			}
			fmt.Printf("Content of %s:\n%s\n", files[i].Name(), content)
		}
	}

}
