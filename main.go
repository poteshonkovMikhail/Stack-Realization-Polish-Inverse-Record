package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Stack представляет собой стек
type Stack struct {
	mt    sync.RWMutex
	items []interface{}
}

// NewStack создает новый стек
func NewStack() *Stack {
	return &Stack{
		items: []interface{}{},
	}
}

// Push добавляет элемент на стек
func (s *Stack) Push(item interface{}) {
	s.mt.Lock()
	defer s.mt.Unlock()
	s.items = append(s.items, item)
}

// Pop удаляет и возвращает верхний элемент стека
func (s *Stack) Pop() (interface{}, error) {
	s.mt.Lock()
	defer s.mt.Unlock()
	if s.IsEmpty() {
		return nil, errors.New("стек пуст")
	}
	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index] // Удаляем верхний элемент
	return item, nil
}

func (s *Stack) Top() (interface{}, error) {
	s.mt.Lock()
	defer s.mt.Unlock()
	if s.IsEmpty() {
		return nil, errors.New("стек пуст")
	}
	return s.items[len(s.items)-1], nil
}

// IsEmpty проверяет, пуст ли стек
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Порядок операторов
func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "^":
		return 3
	default:
		return 0
	}
}

// Проверка на корректность выражения
func validateExpression(expression string) error {
	tokens := strings.Fields(expression)
	if len(tokens) == 0 {
		return errors.New("пустое выражение")
	}

	operatorCount := 0
	numberCount := 0
	parens := 0

	for i, token := range tokens {
		if _, err := strconv.Atoi(token); err == nil {
			numberCount++
		} else if token == "(" {
			parens++
		} else if token == ")" {
			parens--
			if parens < 0 {
				return errors.New("несоответствие скобок: больше закрывающих, чем открывающих")
			}
		} else if isOperator(token) {
			operatorCount++
			if i == 0 || i == len(tokens)-1 || isOperator(tokens[i-1]) {
				return errors.New("недопустимый оператор перед или после: " + token)
			}
		} else {
			return errors.New("недопустимый токен: " + token)
		}
	}
	if parens != 0 {
		return errors.New("несоответствие скобок: количество открывающих и закрывающих скобок не совпадает")
	}
	if operatorCount+1 != numberCount {
		return errors.New("ошибка в количестве операторов и чисел")
	}
	return nil
}

func isOperator(token string) bool {
	switch token {
	case "+", "-", "*", "/", "^":
		return true
	}
	return false
}

// Функция для перевода в постфиксную запись
func infixToPostfix(expression string) (string, error) {
	if err := validateExpression(expression); err != nil {
		return "", err
	}

	var output []string
	stack := NewStack()

	tokens := strings.Fields(expression)

	for _, token := range tokens {
		if _, err := strconv.Atoi(token); err == nil {
			output = append(output, token) // Если токен - число, добавляем его к выходу
		} else if token == "(" {
			stack.Push(token) // Если токен - открывающая скобка, помещаем его в стек
		} else if token == ")" {
			for !stack.IsEmpty() {
				popVal, _ := stack.Pop()
				if popVal == "(" {
					break
				}
				output = append(output, popVal.(string))
			}
			// Проверяем на несоответствие скобок
			if len(output) == 0 {
				return "", errors.New("несоответствие скобок")
			}
		} else { // Если токен - оператор
			for !stack.IsEmpty() {
				topVal, _ := stack.Top()
				if precedence(topVal.(string)) < precedence(token) {
					break
				}
				popVal, err := stack.Pop()
				if err != nil {
					return "", fmt.Errorf("%v", err)
				}
				output = append(output, popVal.(string))
			}
			stack.Push(token) // Добавляем текущий оператор в стек
		}
	}

	// Исчерпываем оставшиеся операторы в стеке
	for !stack.IsEmpty() {
		popVal, _ := stack.Pop()
		output = append(output, popVal.(string))
	}

	return strings.Join(output, " "), nil
}

func main() {
	expression := "( 1 + 2" // Пример входного выражения (в качестве разделителей пробелы)
	postfix, err := infixToPostfix(expression)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Println("Польская инверсная запись:", postfix)
	}
}
