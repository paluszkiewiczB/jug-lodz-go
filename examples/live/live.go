package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type person struct {
	id   int
	name string
	age  int
}

func (p person) withName(n name) person {
	if p.id != n.id {
		panic(fmt.Sprintf("cannot merge person: %#v with name: %#v", p, n))
	}

	p.name = n.name
	return p
}

func (p person) withAge(a age) person {
	if p.id != a.id {
		panic(fmt.Sprintf("cannot merge person: %#v with age: %#v", p, a))
	}

	p.age = a.age
	return p
}

type age struct {
	id  int
	age int
}

func (a age) toPerson() person {
	return person{id: a.id, age: a.age}
}

type name struct {
	id   int
	name string
}

func (a name) toPerson() person {
	return person{id: a.id, name: a.name}
}

func main() {
	start := time.Now()
	defer func() {
		log.Printf("took: %s", time.Since(start))
	}()

	namesC := make(chan name, 5)
	agesC := make(chan age, 5)
	names := generateNames(10)
	ages := generateAges(10)

	go func() {
		defer close(namesC)
		for n, more := names(); more; n, more = names() {
			namesC <- n
		}
	}()

	go func() {
		defer close(agesC)
		for a, more := ages(); more; a, more = ages() {
			agesC <- a
		}
	}()

	people := make(chan person, 5)
	go merge(namesC, agesC, people)
	for p := range people {
		fmt.Printf("%#v\n", p)
	}

}

func merge(names <-chan name, ages <-chan age, people chan<- person) {
	defer close(people)

	merged := make(map[int]person)

	var (
		out   []person
		first person
		outC  chan<- person
	)

	for {
		if names == nil && ages == nil {
			break
		}

		if len(out) != 0 {
			first = out[0]
			outC = people
		}

		select {
		case n, ok := <-names:
			if !ok {
				log.Printf("names closed")
				names = nil
			}
			if p, ok := merged[n.id]; ok {
				out = append(out, p.withName(n))
				delete(merged, n.id)
			} else {
				merged[n.id] = n.toPerson()
			}
		case a, ok := <-ages:
			if !ok {
				log.Printf("ages closed")
				ages = nil
			}
			if p, ok := merged[a.id]; ok {
				out = append(out, p.withAge(a))
				delete(merged, a.id)
			} else {
				merged[a.id] = a.toPerson()
			}
		case outC <- first:
			out = out[1:]
			outC = nil
		}
	}
}

func generateNames(count int) func() (name, bool) {
	var i int
	return func() (name, bool) {
		defer func() {
			i++
		}()

		if i >= count {
			return name{}, false
		}

		log.Printf("generating name: %d", i)
		sleepMs(50, 25)
		return name{id: i, name: fmt.Sprintf("name %d", i)}, true
	}
}

func generateAges(count int) func() (age, bool) {
	var i int
	return func() (age, bool) {
		defer func() {
			i++
		}()

		if i >= count {
			return age{}, false
		}

		log.Printf("generating age: %d", i)
		sleepMs(250, 100)
		return age{id: i, age: i}, true
	}
}

func sleepMs(ms, variation int) {
	varTime := time.Duration(rand.Intn(variation))
	time.Sleep((time.Duration(ms) + varTime) * time.Millisecond)
}
