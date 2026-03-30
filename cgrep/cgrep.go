package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 1. 사용할 옵션(플래그)들을 정의합니다.
	// 기본값은 false, 마지막은 도움말(Help) 설명입니다.
	ignoreCase := flag.Bool("i", false, "대소문자 무시 (Ignore case)")
	showLineNum := flag.Bool("n", false, "줄 번호 출력 (Show line number)")
	afterNum := flag.Int("A", 0, "이후 줄 출력")
	beforeNum := flag.Int("B", 0, "이전 줄 출력")

	// 2. 입력받은 인자들을 분석(Parsing)합니다. (필수!)
	flag.Parse()

	// 3. flag.Args()는 옵션(-i, -n 등)을 제외하고 남은 순수 인자들만 반환합니다.
	remainingArgs := flag.Args()

	if len(remainingArgs) < 2 {
		fmt.Println("사용법: cgrep [-i] [-n] [-A] [-B] [파일명] [단어]")
		return
	}

	fileName := remainingArgs[0]
	keyword := remainingArgs[1]

	// 결과 확인
	fmt.Println("--- 입력 분석 결과 ---")
	fmt.Printf("파일명: %s\n", fileName)
	fmt.Printf("키워드: %s\n", keyword)

	// 포인터 변수이므로 *를 붙여서 값을 확인합니다.
	if *ignoreCase {
		fmt.Println("옵션: 대소문자를 무시합니다.")
		keyword = strings.ToLower(keyword)
	}
	if *showLineNum {
		fmt.Println("옵션: 줄 번호를 출력합니다.")
	}

	absPath, err := filepath.Abs(fileName)
	if err != nil {
		fmt.Println("파일을 찾을 수 없습니다. ", err)
		return
	}

	f, err := os.Open(absPath)
	if err != nil {
		fmt.Println("파일을 읽는 중 오류가 발생하였습니다. ", err)
		return
	}

	defer f.Close()

	sc := bufio.NewScanner(f)
	lineNum := 1

	var beforeSl []string
	var afterCount int

	for sc.Scan() {
		originLine := sc.Text()
		compLine := originLine
		if *ignoreCase {
			compLine = strings.ToLower(originLine)
		}

		if strings.Contains(compLine, keyword) {
			afterCount = *afterNum

			for _, val := range beforeSl {
				fmt.Println(val)
			}
			beforeSl = nil

			if *showLineNum {
				fmt.Printf("%d: %s\n", lineNum, originLine)
			} else {
				fmt.Println(originLine)
			}
		} else {
			if afterCount > 0 {
				fmt.Println(originLine)
				afterCount--

			} else {
				if *beforeNum > 0 {
					beforeSl = append(beforeSl, originLine)
					if len(beforeSl) > *beforeNum {
						beforeSl = beforeSl[1:]
					}
				}
			}
		}
		lineNum++
	}

	if err = sc.Err(); err != nil {
		fmt.Printf("파일을 읽는 중 오류 발생: %v", err)
	}
}
