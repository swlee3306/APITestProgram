package apirequest

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func ApiReqRun() {
		// 종료 신호를 처리하기 위한 채널
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	
		for {
			fmt.Println("요청할 방법을 선택하세요: (1) GET (2) POST (3) POST(등록된 body file) (4) json body 내용 보기 (5) 종료")
			method, err := scanData()
	
			if err != nil {
				fmt.Printf("데이터 읽기에 실패하였습니다.: %s", err)
				method = "0"
			}
	
			switch method {
			case "1":
				HandleGetRequest()
			case "2":
				HandlePostRequest()
			case "3":
				GetJsonList()
			case "4":
				ShowJsonFile()
			case "5":
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