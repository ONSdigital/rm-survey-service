ALTER TABLE survey.classifiertypeselector DROP CONSTRAINT surveyfk_fkey;
ALTER TABLE survey.classifiertypeselector ADD CONSTRAINT surveyfk_fkey FOREIGN KEY (survey_fk) REFERENCES survey.survey(survey_pk);

ALTER TABLE survey.classifiertype DROP CONSTRAINT classifiertypeselectorfk_fkey;
ALTER TABLE survey.classifiertype ADD CONSTRAINT classifiertypeselectorfk_fkey FOREIGN KEY (classifier_type_selector_fk) REFERENCES survey.classifiertypeselector(classifier_type_selector_pk);
