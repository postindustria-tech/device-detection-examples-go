param (
    [Parameter(Mandatory=$true)]
    [string]$RepoName,
    [Parameter(Mandatory=$true)]
    [string]$OrgName,
    [bool]$DryRun = $false
)

Push-Location $RepoName
try {
    $env:GOPROXY = "direct"
    go get -u ./... || $(throw "'go get -u' failed")
    go mod tidy || $(throw "'go mod tidy' failed")
} finally {
    Pop-Location
}
