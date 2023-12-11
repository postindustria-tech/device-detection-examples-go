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

function Write-Host-Colored {
    param (
        [Parameter(Mandatory=$true)]
        [string]$Object,
        [Int16]$ColorCode = 34 #DarkBlue
    )
    Write-Host "`e[${ColorCode}m$Object`e[39m"
}
function Write-Error-Colored {
    param (
        [Parameter(Mandatory=$true)]
        [string]$Object
    )
    Write-Host-Colored "::error::$Object" 31 #DarkRed
}

$failures = @()

Push-Location $ExamplesDir
try {
    Write-Host-Colored "Collecting Examples..."
    $all_examples = Get-ChildItem -Recurse -Include *.go -Exclude $ExamplesExcludeFilter -Name
    foreach ($example_file in $all_examples) {
        Write-Host $example_file
    }
    
    foreach ($example_file in $all_examples) {
        Write-Host-Colored "Starting '$example_file'..."

        go run $example_file
        $example_exit_code = $LASTEXITCODE
        Write-Host ""

        if ($example_exit_code -ne 0) {
            $failures += [IO.Path]::Combine($ExamplesDir, $example_file)
            Write-Error-Colored "'$example_file' finished with code $example_exit_code"
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
        Write-Error-Colored "'$next_test_dir' testing finished with code $test_exit_code"
    }
}

$failures_count = $failures.Length
if ($failures_count -ne 0) {
    Write-Host-Colored "Failed ($failures_count):"
    foreach ($next_failed in $failures) {
        Write-Host-Colored "- $next_failed" 31 #DarkRed
    }
    throw "Failed ($failures_count): $failures"
}
