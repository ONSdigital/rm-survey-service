SET schema 'survey';

CREATE SEQUENCE surveyidseq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 999999999999
    CACHE 1;

CREATE SEQUENCE classifiertypeidseq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 999999999999
    CACHE 1;

CREATE TABLE "survey" (
    surveyid bigint NOT NULL,
    survey character varying(20) NOT NULL
);

ALTER TABLE "survey" OWNER TO survey;

ALTER TABLE "survey" ADD CONSTRAINT survey_pkey PRIMARY KEY (surveyid);

CREATE TABLE "classifiertype" (
    classifiertypeid bigint NOT NULL,
    surveyid bigint NOT NULL,
    classifiertype character varying(50) NOT NULL
);

ALTER TABLE "classifiertype" OWNER TO survey;

ALTER TABLE "classifiertype"
    ADD CONSTRAINT classifiertype_pkey PRIMARY KEY (classifiertypeid);

ALTER TABLE "classifiertype"
    ADD CONSTRAINT surveyid_fkey FOREIGN KEY (surveyid) REFERENCES survey(surveyid);
