INSERT INTO survey.survey(surveyid, id, name) VALUES (1, 'cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87', 'BRES');

INSERT INTO survey.classifiertypeselector(classifiertypeselectorid, id, surveyid, classifiertypeselector) VALUES (1, 'efa868fb-fb80-44c7-9f33-d6800a17c4da', 1, 'COLLECTION_INSTRUMENT');
INSERT INTO survey.classifiertypeselector(classifiertypeselectorid, id, surveyid, classifiertypeselector) VALUES (2, 'e119ffd6-6fc1-426c-ae81-67a96f9a71ba', 1, 'COMMUNICATION_TEMPLATE');

INSERT INTO survey.classifiertype(classifiertypeid, classifiertypeselectorid, classifiertype) VALUES (1, 1, 'COLLECTION_EXERCISE');
INSERT INTO survey.classifiertype(classifiertypeid, classifiertypeselectorid, classifiertype) VALUES (2, 1, 'RU_REF');
INSERT INTO survey.classifiertype(classifiertypeid, classifiertypeselectorid, classifiertype) VALUES (3, 2, 'LEGAL_BASIS');
