SELECT classifiertype 
FROM survey.classifiertype
INNER JOIN survey.survey ON classifiertype.surveyid = survey.surveyid
WHERE survey = 'BRES'
ORDER BY classifiertype ASC;