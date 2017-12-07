package models

// SurveySummary represents a summary of a survey.
type SurveySummary struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
}

// Survey represents the details of a survey.
type Survey struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
	Reference string `json:"surveyRef"`
}

type ModelsInterface interface {
	AllSurveys() *SurveySummary
	GetSurvey() *Survey
	GetSurveyByShortName() *Survey
	GetSurveyByReference() *Survey
}

type Model struct {
}

// AllSurveys returns summaries of all known surveys. The surveys are returned in ascending short name order.
func (m *Model) AllSurveys() ([]*SurveySummary, error) {
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

// GetSurvey returns the details of the survey identified by the string surveyID.
func (m *Model) GetSurvey(surveyID string) (*Survey, error) {
	survey := new(Survey)
	err := db.QueryRow("SELECT id, shortname, longname, surveyref from survey.survey WHERE id = $1", surveyID).Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference)
	if err != nil {
		return nil, err
	}

	return survey, nil
}

// GetSurveyByShortName returns the details of the survey identified by the string shortName.
func (m *Model) GetSurveyByShortName(shortName string) (*Survey, error) {
	survey := new(Survey)

	query := "SELECT id, shortname, longname, surveyref from survey.survey WHERE LOWER(shortName) = LOWER($1)"
	err := db.QueryRow(query, shortName).
		Scan(&survey.ID,
			&survey.ShortName,
			&survey.LongName,
			&survey.Reference)

	if err != nil {
		return nil, err
	}

	return survey, nil
}

// GetSurveyByReference returns the details of the survey identified by the string reference.
func (m *Model) GetSurveyByReference(reference string) (*Survey, error) {
	survey := new(Survey)
	err := db.QueryRow("SELECT id, shortname, longname, surveyref from survey.survey WHERE LOWER(surveyref) = LOWER($1)", reference).Scan(&survey.ID, &survey.ShortName, &survey.LongName, &survey.Reference)
	if err != nil {
		return nil, err
	}

	return survey, nil
}

func (m *Model) getSurveyID(surveyID string) error {
	var id string
	return db.QueryRow("SELECT id FROM survey.survey WHERE id = $1", surveyID).Scan(&id)
}
