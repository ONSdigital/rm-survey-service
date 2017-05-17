DROP TABLE survey.classifiertype;
DROP TABLE survey.classifiertypeselector;
DROP TABLE survey.survey;
DROP SCHEMA survey;
REVOKE CONNECT ON DATABASE postgres FROM survey;
DROP ROLE survey;