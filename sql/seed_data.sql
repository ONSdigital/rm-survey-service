INSERT INTO survey.survey(surveyid, id, survey) VALUES (1, 'cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87', 'BRES');

INSERT INTO survey.classifiertypeselector(classifiertypeselectorid, surveyid, classifiertypeselector) VALUES (1, 1, 'COLLECTION_INSTRUMENT');
INSERT INTO survey.classifiertypeselector(classifiertypeselectorid, surveyid, classifiertypeselector) VALUES (2, 1, 'COMMUNICATION_TEMPLATE');

INSERT INTO survey.classifiertype(classifiertypeid, classifiertypeselectorid, classifiertype) VALUES (1, 1, 'COLLECTION_EXERCISE');
INSERT INTO survey.classifiertype(classifiertypeid, classifiertypeselectorid, classifiertype) VALUES (2, 1, 'RU_REF');
INSERT INTO survey.classifiertype(classifiertypeid, classifiertypeselectorid, classifiertype) VALUES (3, 2, 'LEGAL_BASIS');
