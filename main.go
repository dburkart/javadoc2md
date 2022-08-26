package main

import "fmt"

func main() {
	s := BeginScanning("Foo", `
	/**
	 * This is a javadoc comment
	 *
	 * @param test This is a test
	 */
	`)
	
	go func() {
		for {
			s.State = s.State(s)
			if s.State == nil {
				break
			}
		}
	}()
	
	for {
		t := <- s.Tokens
		
		fmt.Println(t)
		
		if t.Type == TOK_EOF {
			break
		}
	}
}
