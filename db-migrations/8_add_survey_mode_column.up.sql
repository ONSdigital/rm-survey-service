DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'surveymode') THEN
        create type survey.surveymode AS ENUM ('eQ', 'SEFT');
    END IF;
END
$$;
ALTER TABLE survey.survey ADD COLUMN IF NOT EXISTS surveymode survey.surveymode;
UPDATE survey.survey SET surveymode = 'eQ' WHERE shortname IN ('RSI', 'MWSS', 'QBS','MBS');
UPDATE survey.survey SET surveymode = 'SEFT' WHERE shortname IN ('BRES', 'AIFDI', 'AOFDI', 'QIFDI', 'QOFDI',
                                                                  'Sand&Gravel', 'Blocks', 'Bricks', 'PCS', 'ASHE',
                                                                  'NBS', 'OFATS', 'GovERD');
ALTER TABLE survey.survey ALTER COLUMN surveymode SET NOT NULL;