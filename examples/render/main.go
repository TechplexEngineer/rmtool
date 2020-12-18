package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"akeil.net/akeil/rm"
	"akeil.net/akeil/rm/pkg/render"
)

func main() {
	base := "testdata"
	id := "25e3a0ce-080a-4389-be2a-f6aa45ce0207"
	n := rm.NewNotebook(base, id)
	err := n.Read()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for i, p := range n.Pages {
		wg.Add(1)
		go func(i int, p *rm.Page) {
			defer wg.Done()
			log.Printf("Read page %v", i)
			err := p.ReadDrawing()
			if err != nil {
				log.Fatal(err)
			}

			err = p.Drawing.Validate()
			if err != nil {
				log.Printf("Found validation error: %v", err)
			}

			out := fmt.Sprintf("./out/drawing-%v.png", i)
			f, err := os.Create(out)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			w := bufio.NewWriter(f)
			err = render.RenderDrawing(p.Drawing, w)
			if err != nil {
				log.Fatal(err)
			}
			w.Flush()

		}(i, p)
	}

	wg.Wait()
	log.Println("exit ok")
}
