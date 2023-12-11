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

$GroupMarkerColor = "DarkBlue"
$failures = @()

Push-Location $ExamplesDir
Write-Host "::group::$ExamplesDir" -ForegroundColor $GroupMarkerColor
try {
    Write-Host "::group::Collect Examples" -ForegroundColor $GroupMarkerColor
    $all_examples = Get-ChildItem -Recurse -Include *.go -Exclude $ExamplesExcludeFilter -Name
    foreach ($example_file in $all_examples) {
        Write-Host $example_file
    }
    Write-Host "::endgroup::" -ForegroundColor $GroupMarkerColor
    
    foreach ($example_file in $all_examples) {
        Write-Host "::group::$example_file" -ForegroundColor $GroupMarkerColor
        
        go run $example_file
        $example_exit_code = $LASTEXITCODE

        Write-Host ""
        Write-Host "'$example_file' finished with code $example_exit_code" -ForegroundColor ($example_exit_code -eq 0 ? "DarkGreen" : "DarkRed")
        Write-Host "::endgroup::" -ForegroundColor $GroupMarkerColor
        Write-Host ""

        if ($example_exit_code -ne 0) {
            $failures += [IO.Path]::Combine($ExamplesDir, $example_file)
        }
    }
} finally {
    Write-Host "::endgroup::" -ForegroundColor $GroupMarkerColor
    Pop-Location
}

Write-Host "--------------------------"

foreach ($next_test_dir in $TestableDirs) {
    Push-Location $next_test_dir
    Write-Host "::group::$next_test_dir" -ForegroundColor $GroupMarkerColor
    try {
        go test
        $test_exit_code = $LASTEXITCODE

        Write-Host ""
        Write-Host "testing finished with code $test_exit_code" -ForegroundColor ($test_exit_code -eq 0 ? "DarkGreen" : "DarkRed")

        $failures_count += ($example_exit_code -eq 0) ? 0 : 1
        if ($test_exit_code -ne 0) {
            $failures += $next_test_dir
        }
    } finally {
        Write-Host "::endgroup::" -ForegroundColor $GroupMarkerColor
        Pop-Location
    }
}

$failures_count = $failures.Length
if ($failures_count -ne 0) {
    Write-Host "Failed ($failures_count):" -ForegroundColor $GroupMarkerColor
    foreach ($next_failed in $failures) {
        Write-Host "- $next_failed" -ForegroundColor "DarkRed"
    }
    throw "Failed ($failures_count): $failures"
}

