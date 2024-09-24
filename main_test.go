package main

import (
	"testing"
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

func TestInfixToPostfix(t *testing.T) {
	tests := []struct {
		expression string
		expected   string
		shouldFail bool
	}{
		{"( 2 - 8 ) * 5", "2 8 - 5 *", false},
		{"3 + 4 * 2 / ( 1 - 5 ) ^ 2 ^ 3", "3 4 2 * 1 5 - 2 3 ^ ^ / +", false},
		{"( 1 + 2", "", true}, // Несоответствие скобок
		{"1 + * 2", "", true}, // Некорректное выражение
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
