package models

type SurveySummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Survey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Reference string `json:"surveyRef"`
}

func AllSurveys() ([]*SurveySummary, error) {
	rows, err := db.Query("SELECT id, name FROM survey.survey ORDER BY name ASC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	surveySummaries := make([]*SurveySummary, 0)

	for rows.Next() {
		surveySummary := new(SurveySummary)
		err := rows.Scan(&surveySummary.ID, &surveySummary.Name)

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
	err := db.QueryRow("SELECT id, name, surveyref from survey.survey WHERE id = $1", surveyID).Scan(&survey.ID, &survey.Name, &survey.Reference)
	if err != nil {
		return nil, err
	}

	return survey, nil
}

func GetSurveyByName(name string) (*Survey, error) {
	survey := new(Survey)
	err := db.QueryRow("SELECT id, name, surveyref from survey.survey WHERE LOWER(name) = LOWER($1)", name).Scan(&survey.ID, &survey.Name, &survey.Reference)
	if err != nil {
		return nil, err
	}

	return survey, nil
}

func GetSurveyByReference(reference string) (*Survey, error) {
	survey := new(Survey)
	err := db.QueryRow("SELECT id, name, surveyref from survey.survey WHERE LOWER(surveyref) = LOWER($1)", reference).Scan(&survey.ID, &survey.Name, &survey.Reference)
	if err != nil {
		return nil, err
	}

	return survey, nil
}

func getSurveyID(surveyID string) error {
	var id string
	err := db.QueryRow("SELECT id FROM survey.survey WHERE id = $1", surveyID).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}
