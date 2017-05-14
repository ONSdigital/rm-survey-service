package models

type ClassifierTypeSelector struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func AllClassifierTypeSelectors(surveyID string) ([]*ClassifierTypeSelector, error) {
	// We need to run a separate query first to check if the survey exists so an HTTP 404
	// can be correctly returned if it doesn't exist. Without this check an HTTP 204 is
	// incorrectly returned for an invalid survey ID.
	var id string
	err := db.QueryRow("SELECT id from survey.survey WHERE id = $1", surveyID).Scan(&id)
	if err != nil {
		return nil, err
	}

	// Now we can get the classifier type selector records.
	rows, err := db.Query("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyid = survey.surveyid WHERE survey.id = $1 ORDER BY classifiertypeselector", surveyID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	classifierTypeSelectors := make([]*ClassifierTypeSelector, 0)

	for rows.Next() {
		classifierTypeSelector := new(ClassifierTypeSelector)
		err := rows.Scan(&classifierTypeSelector.ID, &classifierTypeSelector.Name)

		if err != nil {
			return nil, err
		}

		classifierTypeSelectors = append(classifierTypeSelectors, classifierTypeSelector)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return classifierTypeSelectors, nil
}
