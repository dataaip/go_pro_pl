package guess_game

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func Guess_game() {
	fmt.Println("beg guess_game")
	rand.Seed(time.Now().UnixNano())
	target := rand.Intn(100) + 1

	reader := bufio.NewReader(os.Stdin)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("输入结束")
				return
			}
			fmt.Println("读取错误:", err)
			continue
		}
		input = strings.TrimSpace(input)

		guess, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("数字转换异常")
			continue
		}

		if guess < 1 || guess > 100 {
			fmt.Println("数字范围超出异常")
			continue
		}

		switch {
		case guess < target:
			fmt.Println("小了")
		case guess > target:
			fmt.Println("大了")
		default:
			fmt.Println("正确")
			return
		}
	}
}
