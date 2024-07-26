UPDATE survey.classifiertype
SET classifier_type='COLLECTION_EXERCISE'
WHERE classifiertype.classifier_type_pk=1
  AND classifiertype.classifier_type_selector_fk=1
  AND classifiertype.classifier_type='FORM_TYPE';
INSERT INTO survey.classifiertype (classifier_type_pk,classifier_type_selector_fk, classifier_type)
VALUES (nextval('survey."classifiertype_classifiertypepk_seq"'), 1, 'RU_REF');
