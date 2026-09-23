package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuizFirstAttemptGetsShareOfCorrectAnswersInPoints(t *testing.T) {
	assignment := CreateAssignmentWithPerAttemptPenaltyAndPoints(t, 10)

	quiz := CreateQuiz(t, "Анна", 1, 12, 15)

	grade := assignment.Evaluate(quiz)

	assert.Equal(t, 80, grade.RawScore)
	assert.Equal(t, 80, grade.FinalScore)
	assert.Equal(t, "80", grade.Display)
}

func TestEssayWithoutPenaltyIsShownAsPass(t *testing.T) {
	assignment := CreateAssignmentWithoutPenaltyAndPassFail(t, 50)

	essay := CreateEssay(t, "Анна", 2, 30, 25, 20)

	grade := assignment.Evaluate(essay)

	assert.Equal(t, 75, grade.RawScore)
	assert.Equal(t, 75, grade.FinalScore)
	assert.Equal(t, "зачёт", grade.Display)
}

func TestProjectRetryIsCappedAndShownOnFivePointScale(t *testing.T) {
	assignment := CreateAssignmentWithRetryCapAndFivePoint(t, 70)

	project := CreateProject(t, "Борис", 2, 90, 100)

	grade := assignment.Evaluate(project)

	assert.Equal(t, 90, grade.RawScore)
	assert.Equal(t, 70, grade.FinalScore)
	assert.Equal(t, "4", grade.Display)
}

func TestGradebookDescribesWorksOfDifferentKinds(t *testing.T) {
	assignment := CreateAssignmentWithPerAttemptPenaltyAndPoints(t, 10)
	gradebook := NewGradebook(assignment)

	gradebook.Add(CreateQuiz(t, "Анна", 1, 12, 15))
	gradebook.Add(CreateEssay(t, "Борис", 2, 30, 25, 20))

	expected := []string{
		"Анна | попытка 1 | верных ответов: 12 из 15 | 80 -> 80 | 80",
		"Борис | попытка 2 | критерии: 30 + 25 + 20 | 75 -> 65 | 65",
	}

	assert.Equal(t, expected, gradebook.Lines())
}

func TestQuizWithMoreCorrectAnswersThanQuestionsIsRejected(t *testing.T) {
	_, err := NewQuizWork(
		"Анна",
		1,
		16,
		15,
	)

	require.Error(t, err)

	var domainErr *DomainRuleViolationError
	require.ErrorAs(t, err, &domainErr)

	assert.Equal(t, "quiz-answers-invalid", domainErr.Code)
}

// =========================
// Helpers
// =========================

func CreateAssignmentWithPerAttemptPenaltyAndPoints(
	t *testing.T,
	penaltyPerAttempt int,
) *Assignment {
	t.Helper()

	penalty, err := NewPerAttemptPenalty(penaltyPerAttempt)
	require.NoError(t, err)

	assignment, err := NewAssignment(
		"Задание",
		penalty,
		PointsScale{},
	)
	require.NoError(t, err)

	return assignment
}

func CreateAssignmentWithoutPenaltyAndPassFail(
	t *testing.T,
	passThreshold int,
) *Assignment {
	t.Helper()

	scale, err := NewPassFailScale(passThreshold)
	require.NoError(t, err)

	assignment, err := NewAssignment(
		"Задание",
		NoPenalty{},
		scale,
	)
	require.NoError(t, err)

	return assignment
}

func CreateAssignmentWithRetryCapAndFivePoint(
	t *testing.T,
	cap int,
) *Assignment {
	t.Helper()

	penalty, err := NewRetryCap(cap)
	require.NoError(t, err)

	assignment, err := NewAssignment(
		"Задание",
		penalty,
		FivePointScale{},
	)
	require.NoError(t, err)

	return assignment
}

func CreateQuiz(
	t *testing.T,
	student string,
	attempt int,
	correct int,
	total int,
) CheckedWork {
	t.Helper()

	quiz, err := NewQuizWork(
		student,
		attempt,
		correct,
		total,
	)
	require.NoError(t, err)

	return quiz
}

func CreateEssay(
	t *testing.T,
	student string,
	attempt int,
	criterionScores ...int,
) CheckedWork {
	t.Helper()

	essay, err := NewEssayWork(
		student,
		attempt,
		criterionScores,
	)
	require.NoError(t, err)

	return essay
}

func CreateProject(
	t *testing.T,
	student string,
	attempt int,
	teamScore int,
	contributionPercent int,
) CheckedWork {
	t.Helper()

	project, err := NewProjectWork(
		student,
		attempt,
		teamScore,
		contributionPercent,
	)
	require.NoError(t, err)

	return project
}
