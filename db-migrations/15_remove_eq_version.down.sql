ALTER TABLE survey.survey ADD eq_version character varying(3) DEFAULT NULL;
UPDATE survey.survey SET eq_version = 'v3' where survey_mode='EQ';