param (
    [Parameter(Mandatory=$true)]
    [string]$RepoName,
    [Parameter(Mandatory=$true)]
    [string]$OrgName,
    [string]$ExamplesExcludeFilter = "example_base.go"
)

$ExamplesDir = "dd"
$TestableDirs = (
    "uach", 
    "web"
)

$DarkRed = 31
$DarkGreen = 32
$DarkYellow = 33
$DarkBlue = 34
$Red = 91

function Add-Color {
    param (
        [Parameter(Mandatory=$true)]
        [string]$Object,
        [Int16]$ColorCode = $DarkBlue
    )
    $ResetColor = "`e[39m"
    $ResetColorRegex = "\W*39m"
    
    $new_color_code = "`e[${ColorCode}m"
    $stripped_object = "$Object" -replace "$ResetColorRegex$",'' -replace "$ResetColorRegex",$new_color_code
    return "$new_color_code$stripped_object$ResetColor"
}
function Build-Exit-Code-Message {
    param (
        [Parameter(Mandatory=$true)]
        [string]$Location,
        [Parameter(Mandatory=$true)]
        [Int32]$ExitCode,
        [bool]$NoColor = $false
    )
    if ($NoColor) {
        return "'$Location' finished with code $ExitCode"
    }
    if ($ExitCode -eq 0) {
        return (Add-Color "'$Location' finished with code $ExitCode" $DarkGreen)
    }
    return "::error::$(Add-Color "'$(Add-Color $Location $DarkYellow)' finished with code $(Add-Color $ExitCode $DarkYellow)" $Red)"
}

$failed_locations = @()

$integrationTestResults = New-Item -ItemType directory -Path ([IO.Path]::Combine($RepoName, "test-results", "integration")) -Force

Push-Location ([IO.Path]::Combine($RepoName, $ExamplesDir))
try {
    Write-Host (Add-Color "Collecting Examples...")
    $all_examples = Get-ChildItem -Recurse -Include *.go -Exclude $ExamplesExcludeFilter -Name
    foreach ($example_file in $all_examples) {
        Write-Host $example_file
    }

    $example_class_name = "github.com/$OrgName/$RepoName/$ExamplesDir"
    $test_report_xml = New-Object -TypeName System.Xml.XmlDocument
    $test_suite = $test_report_xml.CreateElement("testsuite")
    $test_suite.SetAttribute("name", $example_class_name)
    $test_suite.SetAttribute("tests", $all_examples.Length)

    $total_examples_time = 0
    
    foreach ($example_file in $all_examples) {
        Write-Host (Add-Color "Starting '$example_file'...")

        $exec_time = (Measure-Command { 
            go run $example_file
        } | Select-Object TotalSeconds).TotalSeconds
        $example_exit_code = $LASTEXITCODE
        $total_examples_time += $exec_time

        Write-Host ""
        Write-Host (Build-Exit-Code-Message $example_file $example_exit_code)

        $test_case = $test_report_xml.CreateElement("testcase")
        $test_case.SetAttribute("classname", $example_class_name)
        $test_case.SetAttribute("name", $example_file)
        $test_case.SetAttribute("time", $exec_time)

        if ($example_exit_code -ne 0) {
            $failed_locations += ([IO.Path]::Combine($ExamplesDir, (Add-Color $example_file $DarkYellow)))

            $test_failure = $test_report_xml.CreateElement("failure")
            $test_failure.SetAttribute("type", "error")
            $failure_message = (Build-Exit-Code-Message $example_file $example_exit_code -NoColor $true)
            $test_failure.SetAttribute("message", $failure_message)
            $test_failure.InnerText = $failure_message
            $test_case.AppendChild($test_failure)
        }
        $test_suite.AppendChild($test_case)
    }
    $test_suite.SetAttribute("failures", $failed_locations.Length)
    $test_suite.SetAttribute("time", $total_examples_time)
    $test_report_xml.AppendChild($test_suite)
} finally {
    Pop-Location
    $test_report_xml.Save(([IO.Path]::Combine($integrationTestResults, "examples.xml")))
}

foreach ($next_test_dir in $TestableDirs) {
    Push-Location ([IO.Path]::Combine($RepoName, $next_test_dir))
    Write-Host (Add-Color "Testing '$next_test_dir'...")
    try {
        if (Get-Command go-junit-report) {
            $next_results_file = ([IO.Path]::Combine($integrationTestResults, "$next_test_dir.xml"))
            go test -v 2>&1 | go-junit-report -set-exit-code -iocopy -out $next_results_file
            Write-Host (Add-Color "Dumping report:")
            Get-Content $next_results_file
        } else {
            go test
        }
        $test_exit_code = $LASTEXITCODE
        Write-Host ""
    } finally {
        Pop-Location
    }
    
    if ($test_exit_code -ne 0) {
        $failed_locations += (Add-Color $next_test_dir $DarkYellow)
    }
    Write-Host (Build-Exit-Code-Message $next_test_dir $test_exit_code)
}

$failed_locations_count = $failed_locations.Length
if ($failed_locations_count -ne 0) {
    Write-Host (Add-Color "Failed ($failed_locations_count):")
    foreach ($next_failed in $failed_locations) {
        Write-Host (Add-Color "- $next_failed" $DarkRed)
    }
    Write-Host ""
    throw (Add-Color "Failed ($(Add-Color $failed_locations_count $DarkYellow)): $failed_locations" $Red)
}
