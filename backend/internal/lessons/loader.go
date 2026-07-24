package lessons

import "context"

func (s *Service) loadSteps(
	ctx context.Context,
	steps []LessonStep,
) ([]LessonStepResponse, error) {

	wordIDs := make([]int64, 0)
	quizIDs := make([]int64, 0)
	grammarIDs := make([]int64, 0)
	sentenceIDs := make([]int64, 0)

	for _, step := range steps {

		switch step.StepType {

		case "word":
			if step.EntityID != nil {
				wordIDs = append(wordIDs, *step.EntityID)
			}

		case "quiz":
			if step.EntityID != nil {
				quizIDs = append(quizIDs, *step.EntityID)
			}

		case "grammar":
			if step.EntityID != nil {
				grammarIDs = append(grammarIDs, *step.EntityID)
			}

		case "sentence_builder":
			if step.EntityID != nil {
				sentenceIDs = append(sentenceIDs, *step.EntityID)
			}
		}
	}

	words, err := s.wordProvider.GetByIDs(ctx, wordIDs)
	if err != nil {
		return nil, err
	}

	quizzes, err := s.quizProvider.GetByIDs(ctx, quizIDs)
	if err != nil {
		return nil, err
	}

	grammarNotes, err := s.grammarProvider.GetByIDs(ctx, grammarIDs)
	if err != nil {
		return nil, err
	}

	sentenceExercises, err := s.sentenceProvider.GetByIDs(ctx, sentenceIDs)
	if err != nil {
		return nil, err
	}

	wordMap := make(map[int64]any)
	quizMap := make(map[int64]any)
	grammarMap := make(map[int64]any)
	sentenceMap := make(map[int64]any)

	for _, word := range words {
		wordMap[word.ID] = word
	}

	for _, quiz := range quizzes {
		quizMap[quiz.ID] = quiz
	}

	for _, note := range grammarNotes {
		grammarMap[note.ID] = note
	}

	for _, exercise := range sentenceExercises {
		sentenceMap[exercise.ID] = exercise
	}

	result := make([]LessonStepResponse, 0, len(steps))

	for _, step := range steps {

		var data any

		switch step.StepType {

		case "word":
			if step.EntityID != nil {
				data = wordMap[*step.EntityID]
			}

		case "quiz":
			if step.EntityID != nil {
				data = quizMap[*step.EntityID]
			}

		case "grammar":
			if step.EntityID != nil {
				data = grammarMap[*step.EntityID]
			}

		case "sentence_builder":
			if step.EntityID != nil {
				data = sentenceMap[*step.EntityID]
			}

		default:
			data = nil
		}

		result = append(result, LessonStepResponse{
			ID:        step.ID,
			StepType:  step.StepType,
			SortOrder: step.SortOrder,
			Data:      data,
		})
	}

	return result, nil
}
