ALTER TABLE survey.survey ADD COLUMN surveymode character varying(20) NOT NULL;

INSERT INTO survey.survey(surveymode) VALUES('eQ') WHERE shortname IN ('RSI', 'MWSS', 'QBS','MBS')
INSERT INTO survey.survey(surveymode) VALUES('SEFT') WHERE shortname IN ('BRES', 'AIFDI', 'AOFDI', 'QIFDI', 'QOFDI', 'Sand&Gravel', 'Blocks', 'Bricks', 'PCS', 'ASHE', 'NBS', 'OFATS', 'GovERD')
