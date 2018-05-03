ALTER TABLE survey.survey ADD CONSTRAINT survey_reference_unique UNIQUE (surveyref);
ALTER TABLE survey.survey ADD CONSTRAINT survey_shortname_unique UNIQUE (shortname);
