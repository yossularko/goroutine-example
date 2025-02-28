package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("App Start")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/goroutine", goRoutineService)

	app.Get("/goroutine2", goRoutineService2)

	app.Get("/goroutine3", goRoutineService3)

	app.Listen(":3000")
}

// separate goroutine tidak akan tereksekusi jika tidak menggunakan lifecycle http
// func main() {
// 	fmt.Println("App Start")

// 	goRoutineSeparate()
// }

func goRoutineService(c *fiber.Ctx) error {
	time.Sleep(time.Second * 2)
	fmt.Println("all go routine start")

	var wg sync.WaitGroup

	var data []string

	wg.Add(5)
	go func() {
		defer wg.Done()
		goRoutine1()
	}()
	go func() {
		defer wg.Done()
		goRoutine2()
	}()
	go func() {
		defer wg.Done()
		goRoutine3()
	}()
	go func() {
		defer wg.Done()
		goRoutine4()
	}()
	go func() {
		defer wg.Done()
		data = goRoutine5()
	}()

	wg.Wait()

	return c.Status(200).JSON(fiber.Map{"data": data})
}

func goRoutineService2(c *fiber.Ctx) error {
	goRoutineSeparate()

	return c.SendString("Hello, go routine 2!")
}

func goRoutineService3(c *fiber.Ctx) error {
	// tetap tereksekusi karena didalam lifecycle http (walaupun tanpa wait group)
	go func() {
		time.Sleep(time.Second * 3)
		fmt.Println("Hello go routine without wait group")
	}()

	return c.SendString("Hello, go routine 3!")
}

func goRoutineSeparate() {
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		fmt.Println("Hello go routine new 1")
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		fmt.Println("Hello go routine new 2")
	}()

	wg.Wait()

	// separate goroutine | tidak akan tereksekusi jika tidak menggunakan lifecycle http (langsung di main function)
	go func() {
		time.Sleep(time.Second * 3)
		fmt.Println("Hello go routine new 3")
	}()

	fmt.Println("go routine new end")
}

func goRoutine1() {
	go func() {
		time.Sleep(time.Second)
		fmt.Println("Hello go routine 1")
	}()

	time.Sleep(time.Second * 2)
}

func goRoutine2() {
	data := make(chan int)
	go func() {
		time.Sleep(time.Second)
		fmt.Println("Hello go routine 2")
		data <- 100
	}()

	newData := <-data

	fmt.Println("data:", newData)
}

func goRoutine3() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		fmt.Println("Hello go routine 3")
	}()

	wg.Wait()
}

func goRoutine4() {
	var wg sync.WaitGroup
	data := make(chan int, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		fmt.Println("Hello go routine 4a")
		data <- 100
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		fmt.Println("Hello go routine 4b")
		data <- 120
	}()

	wg.Wait()
	close(data)

	for v := range data {
		fmt.Println("data:", v)
	}
}

func goRoutine5() []string {
	fmt.Println("Hello go routine 5")
	values := []int{2, 5, 8, 12, 15, 24, 30}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	errChan := make(chan error, len(values))

	var newValues []string

	for _, value := range values {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			if err := checkInt(v); err != nil {
				errChan <- err
				return
			}

			fmt.Println("user id: ", v)

			mutex.Lock()
			defer mutex.Unlock()
			newValues = append(newValues, fmt.Sprintf("user id %d", v))
		}(value)
	}

	wg.Wait()
	close(errChan)

	for v := range errChan {
		fmt.Println("error:", v)
	}

	return newValues
}

func checkInt(value int) error {
	if value == 5 {
		return fmt.Errorf("test error")
	}

	return nil
}
