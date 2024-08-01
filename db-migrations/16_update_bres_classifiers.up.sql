UPDATE survey.classifiertype
SET classifier_type='FORM_TYPE'
WHERE classifiertype.classifier_type_pk=1
  AND classifiertype.classifier_type_selector_fk=1
  AND classifiertype.classifier_type='COLLECTION_EXERCISE';
DELETE FROM survey.classifiertype
WHERE classifiertype.classifier_type_pk=2
  AND classifiertype.classifier_type_selector_fk=1
  AND classifiertype.classifier_type='RU_REF';
