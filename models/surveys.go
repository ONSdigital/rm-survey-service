package models

type SurveySummary struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
}

type Survey struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
	Reference string `json:"surveyRef"`
}

func AllSurveys() ([]*SurveySummary, error) {
	rows, err := db.Query("SELECT id, shortname FROM survey.survey ORDER BY shortname ASC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	surveySummaries := make([]*SurveySummary, 0)

	for rows.Next() {
		surveySummary := new(SurveySummary)
		err := rows.Scan(&surveySummary.ID, &surveySummary.ShortName)

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
	err := db.QueryRow("SELECT id, shortname, longname, surveyref from survey.survey WHERE id = $1", surveyID).Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference)
	if err != nil {
		return nil, err
	}

	return survey, nil
}

func GetSurveyByShortName(shortName string) (*Survey, error) {
	survey := new(Survey)
	err := db.QueryRow("SELECT id, shortname, longname, surveyref from survey.survey WHERE LOWER(shortName) = LOWER($1)", shortName).Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference)
	if err != nil {
		return nil, err
	}

	return survey, nil
}

func GetSurveyByReference(reference string) (*Survey, error) {
	survey := new(Survey)
	err := db.QueryRow("SELECT id, shortname, longname, surveyref from survey.survey WHERE LOWER(surveyref) = LOWER($1)", reference).Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference)
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
