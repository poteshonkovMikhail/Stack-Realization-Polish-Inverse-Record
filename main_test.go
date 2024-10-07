package main

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestStack(t *testing.T) {
	stack := NewStack()

	// Проверяем, что стек пуст
	if !stack.IsEmpty() {
		t.Error("ожидали, что стек будет пуст")
	}

	// Тестируем Push
	stack.Push(1)
	if stack.IsEmpty() {
		t.Error("ожидали, что стек не будет пуст после добавления элемента")
	}

	// Тестируем Top
	top, err := stack.Top()
	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}
	if top != 1 {
		t.Errorf("ожидали 1, получили %v", top)
	}

	// Тестируем Pop
	popped, err := stack.Pop()
	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}
	if popped != 1 {
		t.Errorf("ожидали 1, получили %v", popped)
	}

	// Проверяем, что стек снова пуст
	if !stack.IsEmpty() {
		t.Error("ожидали, что стек будет пуст после удаления элемента")
	}

	// Тестируем ошибку на Pop из пустого стека
	_, err = stack.Pop()
	if err == nil {
		t.Error("ожидали ошибку при попытке удалить элемент из пустого стека")
	}
}

func generateRandomExpression() (string, string) {
	rand.NewSource(time.Now().UnixNano())
	operators := []string{"+", "-", "*", "/"}
	numCount := rand.Intn(5) + 2 // случайное количество чисел от 2 до 6
	expression := []string{}
	expected := ""

	// Генерация чисел и операторов
	for i := 0; i < numCount; i++ {
		num := rand.Intn(10) + 1 // случайное число от 1 до 10
		expression = append(expression, strconv.Itoa(num))

		if i < numCount-1 {
			// Проверка для предотвращения деления на ноль
			var op string
			for {
				op = operators[rand.Intn(len(operators))]
				if op != "/" || (i < numCount-1 && num != 0) { // если делим, число не должно быть нулем
					break
				}
			}
			expression = append(expression, op)
		}
	}

	// Преобразуем в инфиксное выражение, добавляя скобки случайным образом
	if rand.Intn(2) == 0 {
		// Вставляем скобки
		start := rand.Intn(numCount-1) * 2 // Учитываем, что операторы располагаются между числами
		end := start + 2                   // конца с учётом следующего числа
		if end > len(expression) {
			end = len(expression)
		}
		// Создаем строку с скобками
		expected = strings.Join(expression[:start], " ") + " ( " + strings.Join(expression[start:end], " ") + " ) " + strings.Join(expression[end:], " ")
	} else {
		expected = strings.Join(expression, " ")
	}

	return expected, "" // Ожидаемое значение возвращаем пустым
}

func TestInfixToPostfix(t *testing.T) {
	tests := []struct {
		expression string
		expected   string
		shouldFail bool
	}{
		{"( 3 + 5 ) * 2", "3 5 + 2 *", false},
		{"( 4 - 6 ) / 2", "4 6 - 2 /", false},
		{"7 - ( 2 + 1", "", true}, // Несоответствие скобок
		{"( 3 * + 5 )", "", true}, // Некорректное выражение
		{"5 * ( 3 + 4 ) - 1", "5 3 4 + * 1 -", false},
		{"9 / ( 3 * 3 )", "9 3 3 * /", false},
		{"( 7 / 2", "", true}, // Несоответствие скобок
		{"1 + 2 )", "", true}, // Несоответствие скобок
		{"( 5 - 2 ) + ( 4 / 2 )", "5 2 - 4 2 / +", false},
		{"( 1 + 2 ) * ( 3 - 4 )", "1 2 + 3 4 - *", false},
		{"5 + * 3", "", true}, // Некорректное выражение
		{"10 / ( 2 - ( 4 * 2 ) )", "10 2 4 2 * - /", false},
		{"( ( 1 + 2 )", "", true}, // Несоответствие скобок
		{"8 + 5 - ", "", true},    // Некорректное выражение
		{"( ( 2 ) ) + 3", "2 3 +", false},
		{"((2)) + 3", "2 3 +", true},
	}

	// Генерация случайных тестов
	for i := 0; i < 5; i++ { // Генерируем 5 случайных выражений
		expr, _ := generateRandomExpression()
		// Здесь указано, что ожидаемое значение не определено
		tests = append(tests, struct {
			expression string
			expected   string
			shouldFail bool
		}{expression: expr, expected: "", shouldFail: false}) // Здесь вы можете добавить собственную логику проверок или оставьте это пустым
	}

	for _, test := range tests {
		result, err := infixToPostfix(test.expression)
		if test.shouldFail {
			if err == nil {
				t.Errorf("ожидали ошибку для выражения: %s", test.expression)
			}
		} else {
			if err != nil {
				t.Errorf("неожиданная ошибка для выражения: %s - %v", test.expression, err)
			}
			if result != test.expected {
				t.Errorf("ожидали: %s, получили: %s", test.expected, result)
			}
		}
	}
}
