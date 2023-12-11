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

$failures = @()

Push-Location $ExamplesDir
try {
    $all_examples = Get-ChildItem -Recurse -Include *.go -Exclude $ExamplesExcludeFilter -Name
    foreach ($example_file in $all_examples) {
        Write-Host $example_file
    }
    
    foreach ($example_file in $all_examples) {
        
        go run $example_file
        $example_exit_code = $LASTEXITCODE

        if ($example_exit_code -ne 0) {
            $failures += [IO.Path]::Combine($ExamplesDir, $example_file)
            Write-Host "::error::'$example_file' finished with code $example_exit_code" -ForegroundColor "DarkRed"
        }
    }
} finally {
    Pop-Location
}

foreach ($next_test_dir in $TestableDirs) {
    Push-Location $next_test_dir
    try {
        go test
        $test_exit_code = $LASTEXITCODE
    } finally {
        Pop-Location
    }
    
    if ($test_exit_code -ne 0) {
        $failures += $next_test_dir
        Write-Host "::error::'$next_test_dir' testing finished with code $test_exit_code" -ForegroundColor "DarkRed"
    }
}

$failures_count = $failures.Length
if ($failures_count -ne 0) {
    Write-Host "Failed ($failures_count):" -ForegroundColor "DarkBlue"
    foreach ($next_failed in $failures) {
        Write-Host "- $next_failed" -ForegroundColor "DarkRed"
    }
    throw "Failed ($failures_count): $failures"
}

