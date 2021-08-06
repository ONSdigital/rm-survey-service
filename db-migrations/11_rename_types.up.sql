DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'surveymode') THEN
        ALTER TYPE survey.surveymode RENAME TO survey.survey_mode;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'surveytype') THEN
        ALTER TYPE survey.surveytype RENAME survey.survey_type;
    END IF;
END
$$;