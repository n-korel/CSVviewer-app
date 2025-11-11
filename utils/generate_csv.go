package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	rows := 100000
	
	file, err := os.Create(fmt.Sprintf("test_data_%d.csv", rows))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Заголовки
	writer.Write([]string{"ID", "Name", "Email", "City", "Age"})

	// Генерация строк
	for i := 1; i <= rows; i++ {
		writer.Write([]string{
			strconv.Itoa(i),
			fmt.Sprintf("User %d", i),
			fmt.Sprintf("user%d@email.com", i),
			fmt.Sprintf("City %d", rand.Intn(100)+1),
			strconv.Itoa(rand.Intn(40) + 20),
		})
	}

	fmt.Printf("Сгенерирован файл test_data_%d.csv с %d строками\n", rows, rows)
}