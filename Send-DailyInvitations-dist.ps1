<#
Send-DailyInvitations.ps1
2021-02-05
-John Taylor

This is a simple example of how to create and send a different distribution
each weekday, Monday through Friday.

#>

$configFile = "quautomatrics_config.json"
$outputFile = "example-distribution.json"
$sendDate = "_NOW_" # send emails immediately
$expirationDate = "_DAYS:7_T23:59:59Z" # expire after 7 days

$libraryName = "(to be filled in)"
$surveyName = "(to be filled in)"
$mailingListName = "(to be filled in)"


function main($options) {
    # replace contacts with CSV filename given on cmd-line
    $new_contacts = $options[0]
    echo "new contacts: $new_contacts"
    if( -not (test-path $new_contacts) ) {
        write-error "File not found: $new_contacts"
        return
    }
    $cmd = ".\quautomatrics.exe replaceContacts -m '$mailingListName' -c '$new_contacts'"
    echo $cmd
    Invoke-Expression $cmd

    # create distribution
    rm -force -erroraction silentlycontinue $outputFile
    $today = (get-date).DayOfWeek
    $emailSubject = ""
    $emailMessage = ""
    switch($today)
    {
        "Monday"    { $emailMessage = "Survey-Monday";    $emailSubject = "Your email subject for Monday" }
        "Tuesday"   { $emailMessage = "Survey-Tuesday";   $emailSubject = "Your email subject for Tuesday" }
        "Wednesday" { $emailMessage = "Survey-Wednesday"; $emailSubject = "Your email subject for Wednesday" }
        "Thursday"  { $emailMessage = "Survey-Thursday";  $emailSubject = "Your email subject for Thursday" }
        "Friday"    { $emailMessage = "Survey-Friday";    $emailSubject = "Your email subject for Friday" }
    }

    $cmd = ".\quautomatrics.exe createDistribution -c '$configFile' -o '$outputFile' -l '$libraryName' -s '$surveyName' -m '$emailMessage' -n '$mailingListName' -j '$emailSubject' -d '$sendDate' -e '$expirationDate'"
    echo $cmd
    Invoke-Expression $cmd

    # upload newly created distribution
    if (test-path $outputFile) {
        $cmd = ".\quautomatrics.exe uploadDistribution -d '$outputFile'"
        echo $cmd
        Invoke-Expression $cmd
    }
}

main $args
