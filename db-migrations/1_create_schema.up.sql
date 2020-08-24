-- Create schema
DROP SCHEMA IF EXISTS survey CASCADE;
CREATE SCHEMA survey;

-- Create sequences
CREATE SEQUENCE IF NOT EXISTS survey.survey_surveypk_seq;
CREATE SEQUENCE IF NOT EXISTS survey.classifiertype_classifiertypepk_seq;
CREATE SEQUENCE IF NOT EXISTS survey.classifiertypeselector_classifiertypeselectorpk_seq;

ALTER SEQUENCE survey.survey_surveypk_seq RESTART WITH 1000;
ALTER SEQUENCE survey.classifiertype_classifiertypepk_seq RESTART WITH 1000;
ALTER SEQUENCE survey.classifiertypeselector_classifiertypeselectorpk_seq RESTART WITH 1000;

-- Create tables
CREATE TABLE survey.survey (surveypk serial NOT NULL, id uuid NOT NULL, shortname character varying(20) NOT NULL, longname character varying(100) NOT NULL, surveyref character varying(20) NOT NULL, legalbasis character varying(400) NOT NULL);
ALTER TABLE survey.survey ADD CONSTRAINT survey_pkey PRIMARY KEY (surveypk);
ALTER TABLE survey.survey ADD CONSTRAINT survey_id_key UNIQUE (id);
CREATE TABLE survey.classifiertypeselector (classifiertypeselectorpk serial NOT NULL, id uuid NOT NULL, surveyfk integer NOT NULL, classifiertypeselector character varying(50) NOT NULL);
ALTER TABLE survey.classifiertypeselector ADD CONSTRAINT classifiertypeselector_pkey PRIMARY KEY (classifiertypeselectorpk);
ALTER TABLE survey.classifiertypeselector ADD CONSTRAINT classifiertypeselector_id_key UNIQUE (id);
ALTER TABLE survey.classifiertypeselector ADD CONSTRAINT surveyfk_fkey FOREIGN KEY (surveyfk) REFERENCES survey.survey(surveypk);
CREATE TABLE survey.classifiertype (classifiertypepk serial NOT NULL, classifiertypeselectorfk integer NOT NULL, classifiertype character varying(50) NOT NULL);
ALTER TABLE survey.classifiertype ADD CONSTRAINT classifiertype_pkey PRIMARY KEY (classifiertypepk);
ALTER TABLE survey.classifiertype ADD CONSTRAINT classifiertypeselectorfk_fkey FOREIGN KEY (classifiertypeselectorfk) REFERENCES survey.classifiertypeselector(classifiertypeselectorpk);

