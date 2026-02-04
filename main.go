package main

func main() {
	a, err := NewApp()
	if err != nil {
		panic(err)
	}

	a.Run()
}
