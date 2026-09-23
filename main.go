package main

import "fmt"

func main() {
	// Тест
	quizPenalty, _ := NewPerAttemptPenalty(10)

	quiz, _ := NewAssignment(
		"Тест по теме 3",
		quizPenalty,
		PointsScale{},
	)

	quizBook := NewGradebook(quiz)

	quiz1, _ := NewQuizWork("Анна", 1, 12, 15)
	quiz2, _ := NewQuizWork("Борис", 3, 15, 15)

	quizBook.Add(quiz1)
	quizBook.Add(quiz2)

	// Эссе
	passFailScale, _ := NewPassFailScale(50)

	essay, _ := NewAssignment(
		"Эссе",
		NoPenalty{},
		passFailScale,
	)

	essayBook := NewGradebook(essay)

	essayWork, _ := NewEssayWork(
		"Анна",
		1,
		[]int{30, 25, 20},
	)

	essayBook.Add(essayWork)

	// Проект
	projectPenalty, _ := NewRetryCap(70)

	project, _ := NewAssignment(
		"Командный проект",
		projectPenalty,
		FivePointScale{},
	)

	projectBook := NewGradebook(project)

	project1, _ := NewProjectWork("Анна", 1, 90, 100)
	project2, _ := NewProjectWork("Борис", 2, 90, 100)

	projectBook.Add(project1)
	projectBook.Add(project2)

	// Вывод ведомостей
	assignments := []*Assignment{
		quiz,
		essay,
		project,
	}

	books := []*Gradebook{
		quizBook,
		essayBook,
		projectBook,
	}

	for i := range assignments {
		fmt.Println(assignments[i].Title())

		for _, line := range books[i].Lines() {
			fmt.Printf("  %s\n", line)
		}
	}

	// Проверка некорректных данных
	_, err := NewQuizWork(
		"Вера",
		1,
		16,
		15,
	)

	if err != nil {
		domainErr, ok := err.(*DomainRuleViolationError)
		if ok {
			fmt.Printf(
				"Некорректный результат теста отклонён: %s\n",
				domainErr.Code,
			)
		}
	}
}
