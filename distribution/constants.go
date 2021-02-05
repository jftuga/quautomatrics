package distribution

const JsonDistributionTemplate string = `{
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
    "subject": "__EMAILSUBJECT__"
  },
  "surveyLink": {
    "surveyId": "__SURVEYID__",
    "expirationDate": "__EXPIRATIONDATE__",
    "type": "Individual"
  },
  "sendDate": "__SENDDATE__"
}`

