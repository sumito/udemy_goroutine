package main

import (
	"fmt"
	"runtime"
	"sync"
)

func sayHallo(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Hello")
}

func memConsumed() uint64 {
	runtime.GC()

	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	return s.Sys
}

func Hello(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	fmt.Printf("Hello from %v!\n", id)
}

type value struct {
	mu    sync.Mutex
	value int
	name  string
}

var onceA, onceB sync.Once

func A() {
	fmt.Println("A")
	onceB.Do(B)
}
func B() {
	fmt.Println("B")
	onceB.Do(A)
}

func main() {
	/*
		var wg sync.WaitGroup
		wg.Add(1)
		go sayHallo(&wg)
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("World")
		}()
		wg.Wait()
		time.Sleep(2 * time.Second)
	*/

	/*
		//168
			var wg sync.WaitGroup

			say := "hello"
			wg.Add(1)

			go func() {
				defer wg.Done()
				say = "Good Bye"
			}()

			wg.Wait()

			fmt.Println(say)
	*/

	/*
		//169
		var wg sync.WaitGroup

		tasks := []string{"A", "B", "C"}
		for _, task := range tasks {
			wg.Add(1)
			go func(task string) {
				defer wg.Done()
				fmt.Println(task)
			}(task)
		}
		wg.Wait()
	*/

	/*
		//170
		var ch <-chan interface{}
		var wg sync.WaitGroup
		noop := func() {
			wg.Done()
			<-ch
		}
		const numGoroutines = 1000000
		wg.Add(numGoroutines)
		before := memConsumed()
		for i := 0; i < numGoroutines; i++ {
			go noop()

		}
		wg.Wait()
		after := memConsumed()
		fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
	*/

	/*
		//172 WaitGroup
		var wg = sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("1st Goroutine Start")
			time.Sleep(1 * time.Second)
			fmt.Println("1st Goroutine Done")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("2nd Goroutine Start")
			time.Sleep(1 * time.Second)
			fmt.Println("2nd Goroutine Done")
		}()
		wg.Wait()
	*/

	/*
		//172
		var wg sync.WaitGroup
		var CPU int = runtime.NumCPU()
		wg.Add(CPU)
		for i := 0; i <= CPU; i++ {
			go Hello(&wg, i)
		}
		wg.Wait()
	*/

	/*
		//172 競合 mutex
		var wg sync.WaitGroup
		var memoryAccess sync.Mutex
		var data int

		wg.Add(1)
		go func() {
			defer wg.Done()
			memoryAccess.Lock()
			data++
			memoryAccess.Unlock()
		}()

		wg.Wait()

		memoryAccess.Lock()
		if data == 0 {
			fmt.Println(data)
		} else {
			fmt.Println(data)
		}
		memoryAccess.Unlock()
	*/

	/*
		//173 デッドロック
		var wg sync.WaitGroup

		printSum := func(v1, v2 *value) {
			defer wg.Done()

			v1.mu.Lock()
			fmt.Printf("%v がロックを取得しました\n", v1.name)
			defer v1.mu.Unlock()

			time.Sleep(2 * time.Second)

			v2.mu.Lock()
			fmt.Printf("%v がロックを取得しました\n", v2.name)
			defer v2.mu.Unlock()

			fmt.Println(v1.value + v2.value)
		}

		var a value = value{name: "a"}
		var b value = value{name: "b"}

		wg.Add(2)

		go printSum(&a, &b)
		go printSum(&b, &a)

		wg.Wait()
	*/

	/*
		//175 リソース枯渇
		var wg sync.WaitGroup
		var lock sync.Mutex

		const timer = 1 * time.Second

		greedyWorkder := func() {
			defer wg.Done()
			count := 0

			begin := time.Now()

			for time.Since(begin) <= timer {
				lock.Lock()
				time.Sleep(3 * time.Nanosecond)
				lock.Unlock()
				count++
			}

			fmt.Printf("greedyWorkder: %v \n", count)
		}

		politeWorker := func() {
			defer wg.Done()

			count := 0

			begin := time.Now()

			for time.Since(begin) <= timer {
				lock.Lock()
				time.Sleep(1 * time.Nanosecond)
				lock.Unlock()
				lock.Lock()
				time.Sleep(1 * time.Nanosecond)
				lock.Unlock()
				lock.Lock()
				time.Sleep(1 * time.Nanosecond)
				lock.Unlock()

				count++
			}

			fmt.Printf("politeWorker: %v\n", count)

		}

		wg.Add(2)

		go greedyWorkder()
		go politeWorker()

		wg.Wait()
	*/

	/*
		//176 MutexとRWMutex

		var count int
		var lock sync.RWMutex
		var wg sync.WaitGroup

		increment := func(wg *sync.WaitGroup, l sync.Locker) {
			l.Lock()
			defer l.Unlock()
			defer wg.Done()

			fmt.Println("increment")
			count++
			time.Sleep(1 * time.Second)
		}

		read := func(wg *sync.WaitGroup, l sync.Locker) {
			l.Lock()
			defer l.Unlock()
			defer wg.Done()

			fmt.Println("read")
			fmt.Println(count)
			time.Sleep(1 * time.Second)
		}

		start := time.Now()

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go increment(&wg, &lock)
		}
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go read(&wg, lock.RLocker())
		}

		wg.Wait()

		end := time.Now()

		fmt.Println(end.Sub(start))
	*/

	/*
		//172 Cond
		var mutex sync.Mutex
		cond := sync.NewCond(&mutex)

		for _, name := range []string{"A", "B", "C", "D"} {
			go func(name string) {
				mutex.Lock()
				defer mutex.Unlock()

				cond.Wait()
				fmt.Println(name)
			}(name)
		}

		fmt.Println("Ready...")
		time.Sleep(time.Second)
		fmt.Println("Go!")

		//1件ずつ signal
		//for i := 0; i < 3; i++ {
		//	time.Sleep(time.Second)
		//	cond.Signal()
		//}

		//bloadcat
		cond.Broadcast()
		//

		time.Sleep(time.Second)
		fmt.Println("Done")
	*/

	//173 ライブロック
	/*
		cond := sync.NewCond(&sync.Mutex{})

		go func() {
			for range time.Tick(1 * time.Second) {
				cond.Broadcast()
			}
		}()

		var flag [2]bool

		takeStep := func() {
			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
			fmt.Println(flag)
		}

		var wg sync.WaitGroup

		p0 := func() {
			defer wg.Done()
			flag[0] = true
			takeStep()

			for flag[1] {
				takeStep()
				flag[0] = false
				takeStep()

				if flag[0] != flag[1] {
					break
				}

				takeStep()

				flag[0] = true

				takeStep()
			}
		}

		p1 := func() {
			defer wg.Done()
			flag[1] = true
			takeStep()

			for flag[1] {
				takeStep()
				flag[1] = false
				takeStep()

				if flag[1] != flag[1] {
					break
				}

				takeStep()

				flag[1] = true

				takeStep()
			}
		}

		wg.Add(2)

		go p0()
		go p1()

		wg.Wait()
	*/

	//179 Once
	/*
		count := 0
		increment := func() {
			fmt.Println("increment")
			count++
		}
		decrement := func() {
			fmt.Println("dcrement")
			count--
		}
		var once sync.Once
		once.Do(increment)
		once.Do(decrement)
		fmt.Println(count)
	*/
	//deadlock!
	//onceA.Do(A)

	//180 Pool
	/*
		type Person struct {
			Name string
		}

		mypool := &sync.Pool{
			New: func() interface{} {
				fmt.Println("Create new instance")
				return new(Person)
			},
		}

		mypool.Put(&Person{Name: "1"})
		mypool.Put(&Person{Name: "2"})
		instance1 := mypool.Get()
		instance2 := mypool.Get()
		instance3 := mypool.Get().(*Person)

		fmt.Println(instance1, instance2, instance3)

		instance3.Name = "3"

		fmt.Println(instance1, instance2, instance3)

		mypool.Put(instance1)
		mypool.Put(instance2)
		mypool.Put(instance3)

		instance4 := mypool.Get()
		fmt.Println(instance4)
	*/

	//180 pool - 2
	/*
		count := 0

		mypool := &sync.Pool{
			New: func() interface{} {
				count++
				fmt.Println("Creating")
				return struct{}{}
			},
		}

		mypool.Put("manualy added:1")
		mypool.Put("manualy added:2")

		var wg sync.WaitGroup

		wg.Add(10000)

		for i := 0; i < 10000; i++ {
			time.Sleep(1 * time.Millisecond)

			go func() {
				defer wg.Done()
				instance := mypool.Get()
				mypool.Put(instance)
			}()
		}
		wg.Wait()

		fmt.Printf("created instance: %d\n", count)
	*/

	//182 Map
	/*
		//invalid logic
		var wg sync.WaitGroup
		m := map[string]int{"A": 0, "B": 1}

		for i := 0; i < 10; i++ {
			wg.Add(2)

			go func() {
				defer wg.Done()
				m["A"] = rand.Intn(100)
				m["B"] = rand.Intn(100)
			}()
		}
		wg.Wait()
	*/
	/*
		//182 map - smap
		smap := &sync.Map{}
		smap.Store("Hello", "World")
		smap.Store(1, 2)

		smap.Delete(1)

		v, ok := smap.Load("Hello")
		if ok {
			fmt.Println(v)
		}

		smap.LoadOrStore("Hello", "Wooooooooold")
		smap.LoadOrStore(2, 3)

		smap.Range(func(key, value interface{}) bool {
			fmt.Println(key, value)
			return true
		})
	*/
	/*
		//183 Atomic
		var count int64

		increment := func() {
			atomic.AddInt64(&count, 1)
		}

		increment()

		fmt.Println(count)
	*/

	//185 channel
	/*
		ch := make(chan int)
		go func() {
			defer close(ch)

			for i := 0; i < 10; i++ {
				ch <- i
			}

		}()
		for integer := range ch {
			fmt.Println(integer)
		}
	*/

	/*
		//183 channel close
		begin := make(chan interface{})

		var wg sync.WaitGroup

		for i := 0; i < 5; i++ {
			wg.Add(1)
			fmt.Printf("Start goroutine %d\n", i)

			go func(i int) {
				defer wg.Done()
				<-begin
				fmt.Printf("%d has begin\n", i)
			}(i)
		}
		fmt.Println("unblocking goroutine!")
		close(begin)
		wg.Wait()
	*/

	/*
		//187 channel buffer
		ch := make(chan int, 5)
		go func() {
			for i := 0; i < 5; i++ {
				fmt.Printf("Writing to channelk: %v\n", i)
				ch <- i
			}
		}()
		for integer := range ch {
			time.Sleep(1 * time.Second)
			fmt.Printf("reading to channel: %v\n", integer)
		}
	*/

	/*
		//188 channelライフサイクルのカプセル化
		chanOwner := func() <-chan int {
			resultStream := make(chan int, 5)

			go func() {
				defer close(resultStream)
				for i := 0; i < 5; i++ {
					resultStream <- i
				}
			}()
			return resultStream
		}

		resultStream := chanOwner()

		for result := range resultStream {
			fmt.Printf("Receieved: %v\n", result)
		}

		fmt.Println("Done")
	*/

	//191 select for-select & timeout & canncel & default
	//select
	/*
		a := make(chan int)
		b := make(chan int)
		close(a)

		select {
		case <-b:¥
		case <-a:
		}
	*/

	/*
			start := time.Now()
			ch1 := make(chan int)
			ch2 := make(chan int)

			done := make(chan interface{})

			go func() {
				time.Sleep(2 * time.Second)
				close(done)
			}()

			go func() {
				defer close(ch1)
				for i := 0; i < 10; i++ {
					time.Sleep(1 * time.Second)
					ch1 <- i
				}
			}()
			go func() {
				defer close(ch2)
				for i := 0; i < 10; i++ {
					time.Sleep(1 * time.Second)
					ch2 <- i
				}
			}()

		loop:
			for {
				select {
				case <-done:
					break loop
				case <-time.After(1 * time.Second):
					break loop
				case v, ok := <-ch1:
					if !ok {
						break loop
					}
					fmt.Printf("ch1: %v\n", v)
				case v, ok := <-ch2:
					if !ok {
						break loop
					}
					fmt.Printf("ch2: %v\n", v)
				}

			}
			end := time.Now()
			fmt.Println(end.Sub(start))
	*/

	/*
		//default
		var ch <-chan int
		select {
		case <-ch:
		default:
			fmt.Println("Default")

		}
	*/

}
