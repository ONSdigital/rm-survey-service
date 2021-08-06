DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'surveymode') THEN
        ALTER TYPE survey.surveymode RENAME TO survey_mode;
    END IF;
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'surveytype') THEN
        ALTER TYPE survey.surveytype RENAME TO survey_type;
    END IF;
END
$$;