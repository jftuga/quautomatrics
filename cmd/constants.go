package cmd

const pgmName string = "quautomatrics"
const pgmVersion string = "1.3.0"
const pgmURL string = "https://github.com/jftuga/quautomatrics"

// used by createDistribution
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
