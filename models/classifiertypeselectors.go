package models

type ClassifierTypeSelectorSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClassifierTypeSelector struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	ClassifierTypes []string `json:"classifierTypes"`
}

func AllClassifierTypeSelectors(surveyID string) ([]*ClassifierTypeSelectorSummary, error) {
	// We need to run a query first to check if the survey exists so an HTTP 404 can be correctly
	// returned if it doesn't exist. Without this check an HTTP 204 is incorrectly returned for an
	// invalid survey ID.
	err := getSurveyID(surveyID)
	if err != nil {
		return nil, err
	}

	// Now we can get the classifier type selector records.
	rows, err := db.Query("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyfk = survey.surveypk WHERE survey.id = $1 ORDER BY classifiertypeselector ASC", surveyID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	classifierTypeSelectorSummaries := make([]*ClassifierTypeSelectorSummary, 0)

	for rows.Next() {
		classifierTypeSelectorSummary := new(ClassifierTypeSelectorSummary)
		err := rows.Scan(&classifierTypeSelectorSummary.ID, &classifierTypeSelectorSummary.Name)

		if err != nil {
			return nil, err
		}

		classifierTypeSelectorSummaries = append(classifierTypeSelectorSummaries, classifierTypeSelectorSummary)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return classifierTypeSelectorSummaries, nil
}

func GetClassifierTypeSelector(surveyID string, classifierTypeSelectorID string) (*ClassifierTypeSelector, error) {
	// We need to run two queries first to check if the survey and classifier type selector both exist
	// so an HTTP 404 can be correctly returned if the don't exist. Without this check and HTTP 204 is
	// incorrectly return for an invalid survey ID or classifier type selector ID.
	err := getSurveyID(surveyID)
	if err != nil {
		return nil, err
	}

	err = getClassifierTypeSelectorID(classifierTypeSelectorID)
	if err != nil {
		return nil, err
	}

	// Now we can get the classifier type selector and classifier type records.
	rows, err := db.Query("SELECT id, classifiertypeselector, classifiertype FROM survey.classifiertype INNER JOIN survey.classifiertypeselector ON classifiertype.classifiertypeselectorfk = classifiertypeselector.classifiertypeselectorpk WHERE classifiertypeselector.id = $1 ORDER BY classifiertype ASC", classifierTypeSelectorID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	classifierTypeSelector := new(ClassifierTypeSelector)

	// Using make here ensures the JSON contains an empty array if there are no classifier
	// types, rather than null.
	classifierTypes := make([]string, 0)
	var classifierType string

	for rows.Next() {
		err := rows.Scan(&classifierTypeSelector.ID, &classifierTypeSelector.Name, &classifierType)

		if err != nil {
			return nil, err
		}

		classifierTypes = append(classifierTypes, classifierType)
	}

	classifierTypeSelector.ClassifierTypes = classifierTypes

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return classifierTypeSelector, nil
}

func getClassifierTypeSelectorID(classifierTypeSelectorID string) error {
	var id string
	err := db.QueryRow("SELECT id FROM survey.classifiertypeselector WHERE id = $1", classifierTypeSelectorID).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}
