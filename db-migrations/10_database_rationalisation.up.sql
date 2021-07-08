ALTER TABLE survey.classifiertypeselector RENAME COLUMN classifiertypeselectorpk TO classifier_type_selector_pk;
ALTER TABLE survey.classifiertypeselector RENAME COLUMN classifiertypeselector TO classifier_type_selector;
ALTER TABLE survey.classifiertypeselector RENAME COLUMN surveypk TO survey_pk;

ALTER TABLE survey.classifiertype RENAME COLUMN classifiertypepk TO classifier_type_pk;
ALTER TABLE survey.classifiertype RENAME COLUMN classifiertypeselectorfk TO classifier_type_selector_fk;
ALTER TABLE survey.classifiertype RENAME COLUMN classifiertype TO classifier_type;

ALTER TABLE survey.legalbasis RENAME COLUMN longname TO long_name;

ALTER TABLE survey.survey RENAME COLUMN surveypk TO survey_pk;
ALTER TABLE survey.survey RENAME COLUMN shortname TO short_name;
ALTER TABLE survey.survey RENAME COLUMN longname TO long_name;
ALTER TABLE survey.survey RENAME COLUMN surveyref TO survey_ref;
ALTER TABLE survey.survey RENAME COLUMN legalbasis TO legal_basis;
ALTER TABLE survey.survey RENAME COLUMN surveytype TO survey_type;
ALTER TABLE survey.survey RENAME COLUMN surveymode TO survey_mode;