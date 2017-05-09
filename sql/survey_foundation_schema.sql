SET schema 'survey';

CREATE TABLE "survey" (
    surveyid serial NOT NULL,
    survey character varying(20) NOT NULL,
    uuid uuid NOT NULL
);

ALTER TABLE "survey" OWNER TO survey;

ALTER TABLE "survey" ADD CONSTRAINT survey_pkey PRIMARY KEY (surveyid);
ALTER TABLE "survey" ADD CONSTRAINT survey_uuid_key UNIQUE (uuid);

CREATE TABLE "classifiertypeselector" (
    classifiertypeselectorid serial NOT NULL,
    surveyid integer NOT NULL,
    classifiertypeselector character varying(50) NOT NULL
);

ALTER TABLE "classifiertypeselector" OWNER TO survey;

ALTER TABLE "classifiertypeselector"
    ADD CONSTRAINT classifiertypeselector_pkey PRIMARY KEY (classifiertypeselectorid);

ALTER TABLE "classifiertypeselector"
    ADD CONSTRAINT surveyid_fkey FOREIGN KEY (surveyid) REFERENCES survey(surveyid);

CREATE TABLE "classifiertype" (
    classifiertypeid serial NOT NULL,
    classifiertypeselectorid integer NOT NULL,
    classifiertype character varying(50) NOT NULL
);

ALTER TABLE "classifiertype" OWNER TO survey;

ALTER TABLE "classifiertype"
    ADD CONSTRAINT classifiertype_pkey PRIMARY KEY (classifiertypeid);

ALTER TABLE "classifiertype"
    ADD CONSTRAINT classifiertypeselectorid_fkey FOREIGN KEY (classifiertypeselectorid) REFERENCES classifiertypeselector(classifiertypeselectorid);
