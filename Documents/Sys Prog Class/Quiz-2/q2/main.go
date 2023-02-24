package main
import "fmt"

func main() {

  // create channel
  age := make(chan int)
  name := make(chan string)

  // call function with goroutine
  go channelData(age, name)

  // retrieve data from channel
  fmt.Println("Age: ", <-age)
  fmt.Println("Name: ", <-name)

}

func channelData(age chan int, name chan string) {

  // send data into channel
  age <- 20
  name <- "Tamika Chen"
}