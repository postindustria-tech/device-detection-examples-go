param (
    [Parameter(Mandatory=$true)]
    [string]$RepoName
)

# ./go/run-unit-tests.ps1 -RepoName $RepoName

$RefScript = [IO.Path]::Combine($RepoName, "ci", "integration-tests.ps1")
Push-Location $CIDir
try {
    pwsh $RefScript -RepoName $RepoName
} finally {
    Pop-Location
}
