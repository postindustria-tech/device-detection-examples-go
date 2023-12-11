param (
    [Parameter(Mandatory=$true)]
    [string]$RepoName,
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
        [string]$Location,
        [Int32]$ExitCode
    )
    if ($ExitCode -eq 0) {
        return (Add-Color "'$Location' finished with code $ExitCode" $DarkGreen)
    }
    return "::error::$(Add-Color "'$(Add-Color $Location $DarkYellow)' finished with code $(Add-Color $ExitCode $DarkYellow)" $Red)"
}

$failures = @()

Push-Location ([IO.Path]::Combine($RepoName, $ExamplesDir))
try {
    Write-Host (Add-Color "Collecting Examples...")
    $all_examples = Get-ChildItem -Recurse -Include *.go -Exclude $ExamplesExcludeFilter -Name
    foreach ($example_file in $all_examples) {
        Write-Host $example_file
    }
    
    foreach ($example_file in $all_examples) {
        Write-Host (Add-Color "Starting '$example_file'...")

        go run $example_file
        $example_exit_code = $LASTEXITCODE
        Write-Host ""

        if ($example_exit_code -ne 0) {
            $failures += [IO.Path]::Combine($ExamplesDir, (Add-Color $example_file $DarkYellow))
        }
        Write-Host (Build-Exit-Code-Message $example_file $example_exit_code)
    }
} finally {
    Pop-Location
}

$integrationTestResults = New-Item -ItemType directory -Path $RepoName/test-results/integration -Force

foreach ($next_test_dir in $TestableDirs) {
    Push-Location ([IO.Path]::Combine($RepoName, $next_test_dir))
    Write-Host (Add-Color "Testing '$next_test_dir'...")
    try {
        go test | go-junit-report -set-exit-code -iocopy -out $integrationTestResults/$next_test_dir.xml
        $test_exit_code = $LASTEXITCODE
        Write-Host ""
    } finally {
        Pop-Location
    }
    
    if ($test_exit_code -ne 0) {
        $failures += (Add-Color $next_test_dir $DarkYellow)
    }
    Write-Host (Build-Exit-Code-Message $next_test_dir $test_exit_code)
}

$failures_count = $failures.Length
if ($failures_count -ne 0) {
    Write-Host (Add-Color "Failed ($failures_count):")
    foreach ($next_failed in $failures) {
        Write-Host (Add-Color "- $next_failed" $DarkRed)
    }
    Write-Host ""
    throw (Add-Color "Failed ($(Add-Color $failures_count $DarkYellow)): $failures" $Red)
}
