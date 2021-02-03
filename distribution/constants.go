package distribution

const jsonDistributionTemplate string = `{
 "message": {
    "libraryId": "__LIBRARYID__",
    "messageId": "__MESSAGEID__"
  },
  "recipients": {
    "mailingListId": "__MAILINGLISTID__"
  },
  "header": {
    "fromName": "__FROMNAME__",
    "replyToEmail": "__REPLYTOEMAIL__",
    "fromEmail": "__FROMEMAIL__",
    "subject": "__SUBJECT__"
  },
  "surveyLink": {
    "surveyId": "__SURVEYID__",
    "expirationDate": "__EXPIRATIONDATE__",
    "type": "Individual"
  },
  "sendDate": "__SENDDATE__"
}`

