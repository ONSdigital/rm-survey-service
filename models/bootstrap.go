package models

func bootstrapSQL() [18]string {
	sql := [18]string{
		"CREATE SCHEMA survey",
		"SET schema 'survey'",
		"CREATE TABLE survey (surveypk serial NOT NULL, id uuid NOT NULL, name character varying(20) NOT NULL, surveyref character varying(20) NOT NULL)",
		"ALTER TABLE survey ADD CONSTRAINT survey_pkey PRIMARY KEY (surveypk)",
		"ALTER TABLE survey ADD CONSTRAINT survey_id_key UNIQUE (id)",
		"CREATE TABLE classifiertypeselector (classifiertypeselectorpk serial NOT NULL, id uuid NOT NULL, surveyfk integer NOT NULL, classifiertypeselector character varying(50) NOT NULL)",
		"ALTER TABLE classifiertypeselector ADD CONSTRAINT classifiertypeselector_pkey PRIMARY KEY (classifiertypeselectorpk)",
		"ALTER TABLE classifiertypeselector ADD CONSTRAINT classifiertypeselector_id_key UNIQUE (id)",
		"ALTER TABLE classifiertypeselector ADD CONSTRAINT surveyfk_fkey FOREIGN KEY (surveyfk) REFERENCES survey(surveypk)",
		"CREATE TABLE classifiertype (classifiertypepk serial NOT NULL, classifiertypeselectorfk integer NOT NULL, classifiertype character varying(50) NOT NULL)",
		"ALTER TABLE classifiertype ADD CONSTRAINT classifiertype_pkey PRIMARY KEY (classifiertypepk)",
		"ALTER TABLE classifiertype ADD CONSTRAINT classifiertypeselectorfk_fkey FOREIGN KEY (classifiertypeselectorfk) REFERENCES classifiertypeselector(classifiertypeselectorpk)",
		"INSERT INTO survey (surveypk, id, name, surveyref) VALUES (1, 'cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87', 'BRES', '221');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (1, 'efa868fb-fb80-44c7-9f33-d6800a17c4da', 1, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (2, 'e119ffd6-6fc1-426c-ae81-67a96f9a71ba', 1, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (1, 1, 'COLLECTION_EXERCISE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (2, 1, 'RU_REF');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (3, 2, 'LEGAL_BASIS');",
	}

	return sql
}
