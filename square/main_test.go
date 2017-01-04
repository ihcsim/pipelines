package main

import "testing"

func TestExtract(t *testing.T) {
	var tests = []struct {
		inputs []int
	}{
		{inputs: []int{}},
		{inputs: []int{0}},
		{inputs: []int{1, 2}},
		{inputs: []int{-1, -2}},
		{inputs: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{inputs: []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
	}

	for _, test := range tests {
		actual := extract(test.inputs...)
		index := 0
		for n := range actual {
			if n != test.inputs[index] {
				t.Errorf("Mismatched numbers. Expected %d, but got %d", test.inputs[index], n)
			}
			index++
		}
	}
}

func TestSquare(t *testing.T) {
	var tests = []struct {
		inputs   []int
		expected []int
	}{
		{inputs: []int{}, expected: []int{0}},
		{inputs: []int{0}, expected: []int{0}},
		{inputs: []int{1, 2}, expected: []int{1, 4}},
		{inputs: []int{-1, -2}, expected: []int{1, 4}},
		{inputs: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, expected: []int{0, 1, 4, 9, 16, 25, 36, 49, 64, 81}},
		{inputs: []int{0, -1, -2, -3, -4, -5, -6, -7, -8, -9}, expected: []int{0, 1, 4, 9, 16, 25, 36, 49, 64, 81}},
	}

	for _, test := range tests {
		source := make(chan int)
		go func() {
			for _, num := range test.inputs {
				source <- num
			}
			close(source)
		}()

		var index int
		for actual := range square(source, nil) {
			if actual != test.expected[index] {
				t.Errorf("Mismatched numbers. Input %v. Expected %d, but got %d", test.expected, test.expected[index], actual)
			}
			index++
		}
	}
}

func TestSquareDone(t *testing.T) {
	inputs := []int{0, 1, 2, 3, 4, 5, 6}
	expected := []int{0, 1, 4, 9}

	in := make(chan int, len(inputs))
	for _, input := range inputs {
		in <- input
	}
	close(in)

	// keep sending value to in until the loop reaches stopAtIndex
	stopAtIndex := 3
	done := make(chan struct{})
	actuals := []int{}
	for out := range square(in, done) {
		actuals = append(actuals, out)
		if len(actuals)-1 == stopAtIndex {
			done <- struct{}{}
		}
	}

	for index, actual := range actuals {
		if actual != expected[index] {
			t.Errorf("Mismatched number. Expected %d, but got %d", expected[index], actual)
		}
	}
}

func TestMerge(t *testing.T) {
	inputs := []int{0, 1, 2, 3, 4, 5, 6}
	in1, in2 := make(chan int), make(chan int)
	go func() {
		for _, input := range inputs[:3] {
			in1 <- input
		}
		close(in1)
	}()

	go func() {
		for _, input := range inputs[4:] {
			in2 <- input
		}
		close(in2)
	}()

	for result := range merge(nil, in1, in2) {
		var found bool
		for _, input := range inputs {
			if result == input {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Mismatched number. Expected number to in %v. But got %d", inputs, result)
		}
	}
}

func TestMergeDone(t *testing.T) {
	done := make(chan struct{})
	in1, in2 := make(chan int), make(chan int)
	defer close(in1)
	defer close(in2)

	close(done)

	out := merge(done, in1, in2)
	if _, ok := <-out; ok {
		t.Error("Channel should be closed")
	}
}
