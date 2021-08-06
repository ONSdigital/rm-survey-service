DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'survey_mode') THEN
        ALTER TYPE survey.survey_mode RENAME TO surveymode;
    END IF;
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'survey_type') THEN
        ALTER TYPE survey.survey_type RENAME TO surveytype;
    END IF;
END
$$;