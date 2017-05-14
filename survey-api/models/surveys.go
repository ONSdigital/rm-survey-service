package models

type SurveySummary struct {
	ID     string `json:"id"`
	Survey string `json:"survey"`
}

type Survey struct {
	ID     string `json:"id"`
	Survey string `json:"survey"`
}

func AllSurveys() ([]*SurveySummary, error) {
	rows, err := db.Query("SELECT id, survey FROM survey.survey ORDER BY survey ASC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	surveySummaries := make([]*SurveySummary, 0)

	for rows.Next() {
		surveySummary := new(SurveySummary)
		err := rows.Scan(&surveySummary.ID, &surveySummary.Survey)

		if err != nil {
			return nil, err
		}

		surveySummaries = append(surveySummaries, surveySummary)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return surveySummaries, nil
}

func GetSurvey(surveyID string) (*Survey, error) {
	survey := new(Survey)
	err := db.QueryRow("SELECT id, survey from survey.survey WHERE id = $1", surveyID).Scan(&survey.ID, &survey.Survey)
	if err != nil {
		return nil, err
	}

	return survey, nil
}
