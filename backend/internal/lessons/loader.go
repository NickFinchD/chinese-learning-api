package lessons

import "context"

func (s *Service) loadSteps(
	ctx context.Context,
	steps []LessonStep,
) ([]LessonStepResponse, error) {

	wordIDs := make([]int64, 0)

	for _, step := range steps {

		if step.StepType == "word" && step.EntityID != nil {
			wordIDs = append(wordIDs, *step.EntityID)
		}
	}

	words, err := s.wordProvider.GetByIDs(ctx, wordIDs)
	if err != nil {
		return nil, err
	}

	wordMap := make(map[int64]any)

	for _, word := range words {
		wordMap[word.ID] = word
	}

	result := make([]LessonStepResponse, 0, len(steps))

	for _, step := range steps {

		var data any

		switch step.StepType {

		case "word":

			if step.EntityID != nil {
				data = wordMap[*step.EntityID]
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
