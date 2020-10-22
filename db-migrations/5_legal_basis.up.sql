CREATE TABLE survey.legalbasis (
  ref VARCHAR(20) PRIMARY KEY,
  longname VARCHAR(400) UNIQUE
);

INSERT INTO survey.legalbasis ( ref, longname ) VALUES ( 'GovERD', 'GovERD' );
INSERT INTO survey.legalbasis ( ref, longname ) VALUES ( 'STA1947', 'Statistics of Trade Act 1947' );
INSERT INTO survey.legalbasis ( ref, longname ) VALUES ( 'STA1947_BEIS', 'Statistics of Trade Act 1947 - BEIS' );
INSERT INTO survey.legalbasis ( ref, longname ) VALUES ( 'Vol', 'Voluntary Not Stated' );
INSERT INTO survey.legalbasis ( ref, longname ) VALUES ( 'Vol_BEIS', 'Voluntary - BEIS' );
INSERT INTO survey.legalbasis ( ref, longname ) VALUES ( 'Voluntary', 'Voluntary' );

UPDATE survey.survey SET legalbasis = 'GovERD' WHERE legalbasis = 'GovERD';
UPDATE survey.survey SET legalbasis = 'STA1947' WHERE legalbasis = 'Statistics of Trade Act 1947';
UPDATE survey.survey SET legalbasis = 'STA1947_BEIS' WHERE legalbasis = 'Statistics of Trade Act 1947 - BEIS';
UPDATE survey.survey SET legalbasis = 'Vol' WHERE legalbasis = 'Voluntary Not Stated';
UPDATE survey.survey SET legalbasis = 'Vol_BEIS' WHERE legalbasis = 'Voluntary - BEIS';

ALTER TABLE survey.survey ADD CONSTRAINT survey_legalbasis_fk FOREIGN KEY ( legalbasis ) REFERENCES survey.legalbasis ( ref );



