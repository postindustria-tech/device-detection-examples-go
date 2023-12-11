param (
    [Parameter(Mandatory=$true)]
    [string]$RepoName,
    [string]$ExamplesExcludeFilter = "example_base.go"
)

$ExamplesDir = [IO.Path]::Combine($RepoName, $ExamplesDir, "dd")
$TestableDirs = (
    [IO.Path]::Combine($RepoName, "uach"), 
    [IO.Path]::Combine($RepoName, "web")
)

$failures_count = 0

Push-Location $ExamplesDir
Write-Host "::group::$ExamplesDir" -ForegroundColor "DarkYellow"
try {
    Write-Host "::group::Collect Examples" -ForegroundColor "DarkYellow"
    $all_examples = Get-ChildItem -Recurse -Include *.go -Exclude $ExamplesExcludeFilter -Name
    foreach ($example_file in $all_examples) {
        Write-Host $example_file
    }
    Write-Host "::endgroup::" -ForegroundColor "DarkYellow"
    
    foreach ($example_file in $all_examples) {
        Write-Host "::group::$example_file" -ForegroundColor "DarkYellow"
        
        go run $example_file
        $example_exit_code = $LASTEXITCODE

        Write-Host ""
        Write-Host "'$example_file' finished with code $example_exit_code" -ForegroundColor ($example_exit_code -eq 0 ? "DarkGreen" : "DarkRed")
        Write-Host "::endgroup::" -ForegroundColor "DarkYellow"
        Write-Host ""

        $failures_count += ($example_exit_code -eq 0) ? 0 : 1
    }
} finally {
    Write-Host "::endgroup::" -ForegroundColor "DarkYellow"
    Pop-Location
}

Write-Host "--------------------------"

foreach ($next_test_dir in $TestableDirs) {
    Push-Location $next_test_dir
    Write-Host "::group::$next_test_dir" -ForegroundColor "DarkYellow"
    try {
        go test
        $test_exit_code = $LASTEXITCODE

        Write-Host ""
        Write-Host "testing finished with code $test_exit_code" -ForegroundColor ($test_exit_code -eq 0 ? "DarkGreen" : "DarkRed")

        $failures_count += ($example_exit_code -eq 0) ? 0 : 1
    } finally {
        Write-Host "::endgroup::" -ForegroundColor "DarkYellow"
        Pop-Location
    }
}

Write-Host "Total failures: $failures_count"

return $failures_count
