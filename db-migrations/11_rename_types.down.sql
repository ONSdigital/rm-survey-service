DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'survey_mode') THEN
        ALTER TYPE survey.survey_mode RENAME TO survey.surveymode;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'survey_type') THEN
        ALTER TYPE survey.survey_type RENAME survey.surveytype;
    END IF;
END
$$;