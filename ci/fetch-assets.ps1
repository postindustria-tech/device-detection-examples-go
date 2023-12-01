param (
    [Parameter(Mandatory=$true)]
    [string]$RepoName
)

$assets = New-Item -ItemType Directory -Path assets -Force
$assetsDestination = "$RepoName"
$file = "51Degrees-LiteV4.1.hash"

$downloads = @{
    "51Degrees-LiteV4.1.hash" = {Invoke-WebRequest -Uri "https://github.com/51Degrees/device-detection-data/raw/main/51Degrees-LiteV4.1.hash" -OutFile $assets/$file}
    "20000 User Agents.csv" = {Invoke-WebRequest -Uri "https://media.githubusercontent.com/media/51Degrees/device-detection-data/main/20000%20User%20Agents.csv" -OutFile $assets/$file}
}

foreach ($file in $downloads.Keys) {
    if (!(Test-Path $assets/$file)) {
        Write-Output "Downloading $file"
        Invoke-Command -ScriptBlock $downloads[$file]
    } else {
        Write-Output "'$file' exists, skipping download"
    }
}

New-Item -ItemType SymbolicLink -Force -Target "$assets/51Degrees-LiteV4.1.hash" -Path "$assetsDestination/51Degrees-LiteV4.1.hash"
New-Item -ItemType SymbolicLink -Force -Target "$assets/20000 User Agents.csv" -Path "$assetsDestination/20000 User Agents.csv"
