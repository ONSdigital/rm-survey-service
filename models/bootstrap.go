package models

func bootstrapSQL() [107]string {
	sql := [107]string{
		"DROP SCHEMA survey CASCADE",
		"CREATE SCHEMA survey",
		"SET schema 'survey'",
		"CREATE TABLE survey (surveypk serial NOT NULL, id uuid NOT NULL, shortname character varying(20) NOT NULL, longname character varying(100) NOT NULL, surveyref character varying(20) NOT NULL, legalbasis character varying(400) NOT NULL)",
		"ALTER TABLE survey ADD CONSTRAINT survey_pkey PRIMARY KEY (surveypk)",
		"ALTER TABLE survey ADD CONSTRAINT survey_id_key UNIQUE (id)",
		"CREATE TABLE classifiertypeselector (classifiertypeselectorpk serial NOT NULL, id uuid NOT NULL, surveyfk integer NOT NULL, classifiertypeselector character varying(50) NOT NULL)",
		"ALTER TABLE classifiertypeselector ADD CONSTRAINT classifiertypeselector_pkey PRIMARY KEY (classifiertypeselectorpk)",
		"ALTER TABLE classifiertypeselector ADD CONSTRAINT classifiertypeselector_id_key UNIQUE (id)",
		"ALTER TABLE classifiertypeselector ADD CONSTRAINT surveyfk_fkey FOREIGN KEY (surveyfk) REFERENCES survey(surveypk)",
		"CREATE TABLE classifiertype (classifiertypepk serial NOT NULL, classifiertypeselectorfk integer NOT NULL, classifiertype character varying(50) NOT NULL)",
		"ALTER TABLE classifiertype ADD CONSTRAINT classifiertype_pkey PRIMARY KEY (classifiertypepk)",
		"ALTER TABLE classifiertype ADD CONSTRAINT classifiertypeselectorfk_fkey FOREIGN KEY (classifiertypeselectorfk) REFERENCES classifiertypeselector(classifiertypeselectorpk)",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (1, 'cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87', 'BRES', 'Business Register and Employment Survey', '221', 'Statistics of Trade Act 1947');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (1, 'efa868fb-fb80-44c7-9f33-d6800a17c4da', 1, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (2, 'e119ffd6-6fc1-426c-ae81-67a96f9a71ba', 1, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (1, 1, 'COLLECTION_EXERCISE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (2, 1, 'RU_REF');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (3, 2, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (4, 2, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (2, '75b19ea0-69a4-4c58-8d7f-4458c8f43f5c', 'RSI', 'Monthly Business Survey - Retail Sales Index', '023', 'Statistics of Trade Act 1947');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (3, 'f8bb4b96-e63a-11e7-80c1-9a214cf093ae', 2, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (4, 'f8bb4ccc-e63a-11e7-80c1-9a214cf093ae', 2, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (5, 3, 'FORM_TYPE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (6, 4, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (7, 4, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (3, '41320b22-b425-4fba-a90e-718898f718ce', 'AIFDI', 'Annual Inward Foreign Direct Investment Survey', '062', 'Statistics of Trade Act 1947');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (5, 'f8bb4a6a-e63a-11e7-80c1-9a214cf093ae', 3, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (6, 'f8bb492a-e63a-11e7-80c1-9a214cf093ae', 3, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (8, 6, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (9, 6, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (4, '04dbb407-4438-4f89-acc4-53445d75330c', 'AOFDI', 'Annual Outward Foreign Direct Investment Survey', '063', 'Statistics of Trade Act 1947');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (7, 'f8bb47cc-e63a-11e7-80c1-9a214cf093ae', 4, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (8, 'f8bb44ac-e63a-11e7-80c1-9a214cf093ae', 4, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (10, 8, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (11, 8, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (5, 'c3eaeff3-d570-475d-9859-32c3bf87800d', 'QIFDI', 'Quarterly Inward Foreign Direct Investment Survey', '064', 'Statistics of Trade Act 1947');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (9, 'f8bb4380-e63a-11e7-80c1-9a214cf093ae', 5, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (10, 'f8bb4254-e63a-11e7-80c1-9a214cf093ae', 5, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (12, 9, 'FORM_TYPE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (13, 10, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (14, 10, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (6, '57a43c94-9f81-4f33-bad8-f94800a66503', 'QOFDI', 'Quarterly Outward Foreign Direct Investment Survey', '065', 'Statistics of Trade Act 1947');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (11, 'f8bb4128-e63a-11e7-80c1-9a214cf093ae', 6, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (12, 'f8bb3fe8-e63a-11e7-80c1-9a214cf093ae', 6, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (15, 11, 'FORM_TYPE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (16, 12, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (17, 12, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (7, 'c48d6646-eb6f-4c7c-9f37-f7b41c8d2bc6', 'Sand&Gravel', 'Quarterly Survey of Building Materials Sand and Gravel', '066', 'Statistics of Trade Act 1947 - BEIS');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (13, 'f8bb3e6c-e63a-11e7-80c1-9a214cf093ae', 7, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (14, 'f8bb3ab6-e63a-11e7-80c1-9a214cf093ae', 7, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (18, 13, 'FORM_TYPE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (19, 14, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (20, 14, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (8, '9b6872eb-28ee-4c09-b705-c3ab1bb0f9ec', 'Blocks', 'Monthly Survey of Building Materials Concrete Building Blocks', '073', 'Statistics of Trade Act 1947 - BEIS');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (15, 'f8bb3980-e63a-11e7-80c1-9a214cf093ae', 8, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (16, 'f8bb3840-e63a-11e7-80c1-9a214cf093ae', 8, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (21, 15, 'FORM_TYPE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (22, 16, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (23, 16, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (9, 'cb8accda-6118-4d3b-85a3-149e28960c54', 'Bricks', 'Monthly Survey of Building Materials Bricks', '074', 'Voluntary - BEIS');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (17, 'f8bb361a-e63a-11e7-80c1-9a214cf093ae', 9, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (18, 'f8bb34f8-e63a-11e7-80c1-9a214cf093ae', 9, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (24, 17, 'FORM_TYPE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (25, 18, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (26, 18, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (10, 'c23bb1c1-5202-43bb-8357-7a07c844308f', 'MWSS', 'Monthly Wages and Salaries Survey', '134', 'Statistics of Trade Act 1947');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (19, 'f8bb33c2-e63a-11e7-80c1-9a214cf093ae', 10, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (20, 'f8bb326e-e63a-11e7-80c1-9a214cf093ae', 10, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (27, 20, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (28, 20, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (11, '416b8a82-2031-4f41-b59b-95482d916ca3', 'PCS', 'Public Corporations Survey', '137', 'Voluntary Not Stated');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (21, 'f8bb2f1c-e63a-11e7-80c1-9a214cf093ae', 11, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (22, 'f8bb2df0-e63a-11e7-80c1-9a214cf093ae', 11, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (29, 22, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (30, 22, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (12, '02b9c366-7397-42f7-942a-76dc5876d86d', 'QBS', 'Quarterly Business Survey', '139', 'Statistics of Trade Act 1947');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (23, 'f8bb2cc4-e63a-11e7-80c1-9a214cf093ae', 12, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (24, 'f8bb2b84-e63a-11e7-80c1-9a214cf093ae', 12, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (31, 23, 'FORM_TYPE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (32, 24, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (33, 24, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (13, '6aa8896f-ced5-4694-800c-6cd661b0c8b2', 'ASHE', 'Annual Survey of Hours and Earnings', '141', 'Statistics of Trade Act 1947');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (25, 'f8bb27c4-e63a-11e7-80c1-9a214cf093ae', 13, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (26, 'f8bb2698-e63a-11e7-80c1-9a214cf093ae', 13, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (34, 25, 'FORM_TYPE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (35, 26, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (36, 26, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (14, '7a2c9d6c-9aaf-4cf0-a68c-1d50b3f1b296', 'NBS', 'National Balance Sheet', '199', 'Voluntary Not Stated');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (27, 'f8bb256c-e63a-11e7-80c1-9a214cf093ae', 14, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (28, 'f8bb2422-e63a-11e7-80c1-9a214cf093ae', 14, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (37, 27, 'FORM_TYPE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (38, 28, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (39, 28, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (15, '0fc6fa22-8938-43b6-81c5-f1ccca5a5494', 'OFATS', 'Outward Foreign Affiliates Statistics Survey', '225', 'Statistics of Trade Act 1947');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (29, 'f8bb2184-e63a-11e7-80c1-9a214cf093ae', 15, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (30, 'f8bb2044-e63a-11e7-80c1-9a214cf093ae', 15, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (40, 30, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (41, 30, 'REGION');",
		"INSERT INTO survey (surveypk, id, shortname, longname, surveyref, legalbasis) VALUES (16, 'a81f8a72-47e1-4fcf-a88b-0c175829e02b', 'GovERD', 'Government Research and Development Survey', '500', 'GovERD');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (31, 'f8bb1efa-e63a-11e7-80c1-9a214cf093ae', 16, 'COLLECTION_INSTRUMENT');",
		"INSERT INTO classifiertypeselector (classifiertypeselectorpk, id, surveyfk, classifiertypeselector) VALUES (32, 'f8bb1c52-e63a-11e7-80c1-9a214cf093ae', 16, 'COMMUNICATION_TEMPLATE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (42, 32, 'LEGAL_BASIS');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (43, 32, 'REGION');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (44, 21, 'FORM_TYPE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (45, 5, 'FORM_TYPE');",
		"INSERT INTO classifiertype (classifiertypepk, classifiertypeselectorfk, classifiertype) VALUES (46, 7, 'FORM_TYPE');",
	}

	return sql
}
