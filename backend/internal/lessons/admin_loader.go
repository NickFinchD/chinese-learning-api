package lessons

import "context"

// AdminGetSteps returns a lesson's steps with a human-readable label per
// step (instead of loadSteps' full player-facing payload), for the admin
// step editor.
func (s *Service) AdminGetSteps(ctx context.Context, lessonID int64) ([]AdminStepResponse, error) {

	steps, err := s.repository.GetSteps(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	return s.buildAdminStepResponses(ctx, steps)
}

func (s *Service) buildAdminStepResponses(ctx context.Context, steps []LessonStep) ([]AdminStepResponse, error) {

	wordIDs := make([]int64, 0)
	quizIDs := make([]int64, 0)
	grammarIDs := make([]int64, 0)
	sentenceIDs := make([]int64, 0)

	for _, step := range steps {

		if step.EntityID == nil {
			continue
		}

		switch step.StepType {
		case "word":
			wordIDs = append(wordIDs, *step.EntityID)
		case "quiz":
			quizIDs = append(quizIDs, *step.EntityID)
		case "grammar":
			grammarIDs = append(grammarIDs, *step.EntityID)
		case "sentence_builder":
			sentenceIDs = append(sentenceIDs, *step.EntityID)
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

	wordLabel := make(map[int64]string, len(words))
	for _, w := range words {
		wordLabel[w.ID] = w.Hanzi + " — " + w.Translation
	}

	quizLabel := make(map[int64]string, len(quizzes))
	for _, q := range quizzes {
		quizLabel[q.ID] = q.Question
	}

	grammarLabel := make(map[int64]string, len(grammarNotes))
	for _, n := range grammarNotes {
		grammarLabel[n.ID] = n.Title
	}

	sentenceLabel := make(map[int64]string, len(sentenceExercises))
	for _, e := range sentenceExercises {
		sentenceLabel[e.ID] = e.Translation
	}

	result := make([]AdminStepResponse, 0, len(steps))

	for _, step := range steps {

		var entityID int64
		var label string

		if step.EntityID != nil {
			entityID = *step.EntityID

			switch step.StepType {
			case "word":
				label = wordLabel[entityID]
			case "quiz":
				label = quizLabel[entityID]
			case "grammar":
				label = grammarLabel[entityID]
			case "sentence_builder":
				label = sentenceLabel[entityID]
			}
		}

		result = append(result, AdminStepResponse{
			ID:        step.ID,
			StepType:  step.StepType,
			EntityID:  entityID,
			SortOrder: step.SortOrder,
			Label:     label,
		})
	}

	return result, nil
}
