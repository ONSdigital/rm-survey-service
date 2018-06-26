DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'surveytype') THEN
        create type survey.surveytype AS ENUM ('Business', 'Social', 'Census');
    END IF;
END
$$;
ALTER TABLE survey.survey ADD COLUMN IF NOT EXISTS surveytype survey.surveytype;
UPDATE survey.survey SET surveytype = 'Business' WHERE surveytype IS NULL;
ALTER TABLE survey.survey ALTER COLUMN surveytype SET NOT NULL;