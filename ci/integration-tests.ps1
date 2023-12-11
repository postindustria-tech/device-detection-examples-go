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

$DarkBlue = 34
$DarkRed = 31

function Make-Colorful {
    param (
        [Parameter(Mandatory=$true)]
        [string]$Object,
        [Int16]$ColorCode = $DarkBlue
    )
    return "`e[${ColorCode}m$Object`e[39m"
}

$failures = @()

Push-Location $ExamplesDir
try {
    Write-Host (Make-Colorful "Collecting Examples...")
    $all_examples = Get-ChildItem -Recurse -Include *.go -Exclude $ExamplesExcludeFilter -Name
    foreach ($example_file in $all_examples) {
        Write-Host $example_file
    }
    
    foreach ($example_file in $all_examples) {
        Write-Host (Make-Colorful "Starting '$example_file'...")

        go run $example_file
        $example_exit_code = $LASTEXITCODE
        Write-Host ""

        if ($example_exit_code -ne 0) {
            $failures += [IO.Path]::Combine($ExamplesDir, $example_file)
            Write-Host "::error::"(Make-Colorful "'$example_file' finished with code $example_exit_code" $DarkRed)
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
        Write-Host ""
    } finally {
        Pop-Location
    }
    
    if ($test_exit_code -ne 0) {
        $failures += $next_test_dir
        Write-Host "::error::"(Make-Colorful "'$next_test_dir' testing finished with code $test_exit_code")
    }
}

$failures_count = $failures.Length
if ($failures_count -ne 0) {
    Write-Host (Make-Colorful "Failed ($failures_count):")
    foreach ($next_failed in $failures) {
        Write-Host (Make-Colorful "- $next_failed" $DarkRed)
    }
    Write-Host ""
    throw "Failed ($failures_count): $failures"
}
