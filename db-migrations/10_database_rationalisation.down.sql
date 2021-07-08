ALTER TABLE survey.classifiertypeselector RENAME COLUMN classifier_type_selector_pk TO classifiertypeselectorpk;
ALTER TABLE survey.classifiertypeselector RENAME COLUMN classifier_type_selector TO classifiertypeselector;
ALTER TABLE survey.classifiertypeselector RENAME COLUMN survey_fk TO surveyfk;

ALTER TABLE survey.classifiertype RENAME COLUMN classifier_type_pk TO classifiertypepk;
ALTER TABLE survey.classifiertype RENAME COLUMN classifier_type_selector_fk TO classifiertypeselectorfk;
ALTER TABLE survey.classifiertype RENAME COLUMN classifier_type TO classifiertype;

ALTER TABLE survey.legalbasis RENAME COLUMN long_name TO longname;

ALTER TABLE survey.survey RENAME COLUMN survey_pk TO surveypk;
ALTER TABLE survey.survey RENAME COLUMN short_name TO shortname;
ALTER TABLE survey.survey RENAME COLUMN long_name TO longname;
ALTER TABLE survey.survey RENAME COLUMN survey_ref TO surveyref;
ALTER TABLE survey.survey RENAME COLUMN legal_basis TO legalbasis;
ALTER TABLE survey.survey RENAME COLUMN survey_type TO surveytype;
ALTER TABLE survey.survey RENAME COLUMN survey_mode TO surveymode;