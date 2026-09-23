package main

import (
	"fmt"
	"strings"
)

// =========================
// Domain error
// =========================

type DomainRuleViolationError struct {
	Code    string
	Message string
}

func (e *DomainRuleViolationError) Error() string {
	return e.Message
}

// =========================
// Checked work
// =========================

type CheckedWork interface {
	StudentName() string
	AttemptNumber() int
	RawScore() int
	Details() string
}

// =========================
// Common work data
// =========================

type workBase struct {
	studentName   string
	attemptNumber int
}

func newWorkBase(studentName string, attemptNumber int) (workBase, error) {
	if strings.TrimSpace(studentName) == "" {
		return workBase{}, &DomainRuleViolationError{
			Code:    "student-name-empty",
			Message: "The student name must not be empty.",
		}
	}

	if attemptNumber < 1 {
		return workBase{}, &DomainRuleViolationError{
			Code:    "attempt-number-invalid",
			Message: "The attempt number starts from 1.",
		}
	}

	return workBase{
		studentName:   studentName,
		attemptNumber: attemptNumber,
	}, nil
}

func (w workBase) StudentName() string {
	return w.studentName
}

func (w workBase) AttemptNumber() int {
	return w.attemptNumber
}

// =========================
// Quiz
// =========================

type QuizWork struct {
	workBase
	correctAnswers int
	totalQuestions int
}

func NewQuizWork(
	studentName string,
	attemptNumber int,
	correctAnswers int,
	totalQuestions int,
) (*QuizWork, error) {

	base, err := newWorkBase(studentName, attemptNumber)
	if err != nil {
		return nil, err
	}

	if totalQuestions <= 0 ||
		correctAnswers < 0 ||
		correctAnswers > totalQuestions {

		return nil, &DomainRuleViolationError{
			Code:    "quiz-answers-invalid",
			Message: "Correct answers must be between 0 and the positive number of questions.",
		}
	}

	return &QuizWork{
		workBase:       base,
		correctAnswers: correctAnswers,
		totalQuestions: totalQuestions,
	}, nil
}

func (w QuizWork) RawScore() int {
	return w.correctAnswers * 100 / w.totalQuestions
}

func (w QuizWork) Details() string {
	return fmt.Sprintf(
		"верных ответов: %d из %d",
		w.correctAnswers,
		w.totalQuestions,
	)
}

// =========================
// Essay
// =========================

type EssayWork struct {
	workBase
	criterionScores []int
}

func NewEssayWork(
	studentName string,
	attemptNumber int,
	criterionScores []int,
) (*EssayWork, error) {

	base, err := newWorkBase(studentName, attemptNumber)
	if err != nil {
		return nil, err
	}

	if len(criterionScores) == 0 {
		return nil, &DomainRuleViolationError{
			Code:    "essay-criteria-invalid",
			Message: "There must be at least one criterion.",
		}
	}

	sum := 0

	for _, score := range criterionScores {
		if score < 0 {
			return nil, &DomainRuleViolationError{
				Code:    "essay-criteria-invalid",
				Message: "Criterion scores must be non-negative and give at most 100 in total.",
			}
		}

		sum += score
	}

	if sum > 100 {
		return nil, &DomainRuleViolationError{
			Code:    "essay-criteria-invalid",
			Message: "Criterion scores must be non-negative and give at most 100 in total.",
		}
	}

	scoresCopy := append([]int(nil), criterionScores...)

	return &EssayWork{
		workBase:        base,
		criterionScores: scoresCopy,
	}, nil
}

func (w EssayWork) RawScore() int {
	sum := 0

	for _, score := range w.criterionScores {
		sum += score
	}

	return sum
}

func (w EssayWork) Details() string {
	parts := make([]string, len(w.criterionScores))

	for i, score := range w.criterionScores {
		parts[i] = fmt.Sprintf("%d", score)
	}

	return fmt.Sprintf(
		"критерии: %s",
		strings.Join(parts, " + "),
	)
}

// =========================
// Project
// =========================

type ProjectWork struct {
	workBase
	teamScore           int
	contributionPercent int
}

func NewProjectWork(
	studentName string,
	attemptNumber int,
	teamScore int,
	contributionPercent int,
) (*ProjectWork, error) {

	base, err := newWorkBase(studentName, attemptNumber)
	if err != nil {
		return nil, err
	}

	if teamScore < 0 ||
		teamScore > 100 ||
		contributionPercent < 0 ||
		contributionPercent > 100 {

		return nil, &DomainRuleViolationError{
			Code:    "project-review-invalid",
			Message: "The team score and the contribution percent must be between 0 and 100.",
		}
	}

	return &ProjectWork{
		workBase:            base,
		teamScore:           teamScore,
		contributionPercent: contributionPercent,
	}, nil
}

func (w ProjectWork) RawScore() int {
	return w.teamScore * w.contributionPercent / 100
}

func (w ProjectWork) Details() string {
	return fmt.Sprintf(
		"балл команды %d, вклад %d%%",
		w.teamScore,
		w.contributionPercent,
	)
}

// =========================
// Penalties
// =========================

type Penalty interface {
	Apply(rawScore int, attemptNumber int) int
}

// No penalty

type NoPenalty struct{}

func (NoPenalty) Apply(rawScore int, attemptNumber int) int {
	return rawScore
}

// Per-attempt penalty

type PerAttemptPenalty struct {
	penalty int
}

func NewPerAttemptPenalty(penalty int) (*PerAttemptPenalty, error) {
	if penalty < 0 || penalty > 100 {
		return nil, &DomainRuleViolationError{
			Code:    "penalty-value-invalid",
			Message: "The penalty value must be between 0 and 100.",
		}
	}

	return &PerAttemptPenalty{
		penalty: penalty,
	}, nil
}

func (p PerAttemptPenalty) Apply(rawScore int, attemptNumber int) int {
	result := rawScore - p.penalty*(attemptNumber-1)

	if result < 0 {
		return 0
	}

	return result
}

// Retry cap

type RetryCap struct {
	cap int
}

func NewRetryCap(cap int) (*RetryCap, error) {
	if cap < 0 || cap > 100 {
		return nil, &DomainRuleViolationError{
			Code:    "penalty-value-invalid",
			Message: "The penalty value must be between 0 and 100.",
		}
	}

	return &RetryCap{
		cap: cap,
	}, nil
}

func (p RetryCap) Apply(rawScore int, attemptNumber int) int {
	if attemptNumber == 1 {
		return rawScore
	}

	if rawScore > p.cap {
		return p.cap
	}

	return rawScore
}

// =========================
// Grade scales
// =========================

type GradeScale interface {
	Display(finalScore int) string
}

// Points

type PointsScale struct{}

func (PointsScale) Display(finalScore int) string {
	return fmt.Sprintf("%d", finalScore)
}

// Pass / Fail

type PassFailScale struct {
	threshold int
}

func NewPassFailScale(threshold int) (*PassFailScale, error) {
	if threshold < 0 || threshold > 100 {
		return nil, &DomainRuleViolationError{
			Code:    "pass-threshold-invalid",
			Message: "The pass threshold must be between 0 and 100.",
		}
	}

	return &PassFailScale{
		threshold: threshold,
	}, nil
}

func (s PassFailScale) Display(finalScore int) string {
	if finalScore >= s.threshold {
		return "зачёт"
	}

	return "незачёт"
}

// Five-point

type FivePointScale struct{}

func (FivePointScale) Display(finalScore int) string {
	if finalScore >= 85 {
		return "5"
	}

	if finalScore >= 70 {
		return "4"
	}

	if finalScore >= 50 {
		return "3"
	}

	return "2"
}

// =========================
// Grade
// =========================

type Grade struct {
	RawScore   int
	FinalScore int
	Display    string
}

// =========================
// Assignment
// =========================

type Assignment struct {
	title   string
	penalty Penalty
	scale   GradeScale
}

func NewAssignment(
	title string,
	penalty Penalty,
	scale GradeScale,
) (*Assignment, error) {

	if strings.TrimSpace(title) == "" {
		return nil, &DomainRuleViolationError{
			Code:    "assignment-title-empty",
			Message: "The assignment title must not be empty.",
		}
	}

	return &Assignment{
		title:   title,
		penalty: penalty,
		scale:   scale,
	}, nil
}

func (a Assignment) Title() string {
	return a.title
}

func (a Assignment) Evaluate(work CheckedWork) Grade {
	rawScore := work.RawScore()

	finalScore := a.penalty.Apply(
		rawScore,
		work.AttemptNumber(),
	)

	display := a.scale.Display(finalScore)

	return Grade{
		RawScore:   rawScore,
		FinalScore: finalScore,
		Display:    display,
	}
}

// =========================
// Gradebook
// =========================

type Gradebook struct {
	assignment *Assignment
	works      []CheckedWork
}

func NewGradebook(assignment *Assignment) *Gradebook {
	return &Gradebook{
		assignment: assignment,
		works:      make([]CheckedWork, 0),
	}
}

func (g *Gradebook) Add(work CheckedWork) {
	g.works = append(g.works, work)
}

func (g Gradebook) Lines() []string {
	lines := make([]string, 0, len(g.works))

	for _, work := range g.works {
		grade := g.assignment.Evaluate(work)

		lines = append(
			lines,
			fmt.Sprintf(
				"%s | попытка %d | %s | %d -> %d | %s",
				work.StudentName(),
				work.AttemptNumber(),
				work.Details(),
				grade.RawScore,
				grade.FinalScore,
				grade.Display,
			),
		)
	}

	return lines
}
