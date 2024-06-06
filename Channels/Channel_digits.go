package main

/*
Решить задачу на каналы по построению пайплайна обработки стрима чисел
а) входной стрим по тику раз в 500 мс, случайное целое число помещается в входной канал
б) процедура разделения - входный канал разделяется на 2 (четные и нечетные)
в) для нечетных - если число кратно 3, то пишем в общий выходной канал
г) для четных, если число большо 1000 то пишем пару число+слово Биг, если меньше то пишем число+смалл
д) выходной поток пишем в коносль
е) программа завершается по получению сигнала от ОС (gracefull shutdown)
*/
import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	delay = time.Millisecond * 500
)

func crtlC() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT)
	<-sigChan
	fmt.Println("CTRL+C Pressed to interrupt")
	os.Exit(0)
}
func generateNumbers(ch chan<- int) {
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ch <- rand.Intn(2000)
		}
	}
}

func processEvenNumbers(in <-chan int, out chan<- string) {
	for num := range in {
		if num%2 == 0 && num > 1000 {
			out <- fmt.Sprintf("Big %d", num)
		} else if num%2 == 0 {
			out <- fmt.Sprintf("Small %d", num)
		}
	}
}

func processOddNumbers(in <-chan int, out chan<- string) {
	for num := range in {
		if num%2 != 0 && num%3 == 0 {
			out <- fmt.Sprintf("%d", num)
		}
	}
}

func main() {

	go crtlC()
	rand.Seed(time.Now().UnixNano())
	input := make(chan int)
	output := make(chan string)

	go generateNumbers(input)
	go processEvenNumbers(input, output)
	go processOddNumbers(input, output)

	for result := range output {
		fmt.Println(result)
	}
}
