package main

import (
    "fmt"
    "sync"
)

func one(wg *sync.WaitGroup) {
    fmt.Println("You are tall")
    defer wg.Done() // remove goroutine from waitgroup counter
}
func two(wg *sync.WaitGroup) {
    fmt.Println("You are short")
    defer wg.Done()
}
func main() {
    // new waitgroup
    wg := new(sync.WaitGroup)
    // add two go routines
    wg.Add(2)

    go one(wg)
    go two(wg)

    // block execution until done
    wg.Wait()

}