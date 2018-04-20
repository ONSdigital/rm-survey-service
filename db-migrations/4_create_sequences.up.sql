CREATE SEQUENCE IF NOT EXISTS survey.survey_surveypk_seq;
CREATE SEQUENCE IF NOT EXISTS survey.classifiertype_classifiertypepk_seq;
CREATE SEQUENCE IF NOT EXISTS survey.classifiertypeselector_classifiertypeselectorpk_seq;


ALTER SEQUENCE survey.survey_surveypk_seq RESTART WITH 1000;
ALTER SEQUENCE survey.classifiertype_classifiertypepk_seq RESTART WITH 1000;
ALTER SEQUENCE survey.classifiertypeselector_classifiertypeselectorpk_seq RESTART WITH 1000;
