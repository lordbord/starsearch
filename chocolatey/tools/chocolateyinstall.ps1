$ErrorActionPreference = 'Stop'

$packageName = 'starsearch'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url = 'https://github.com/lordbord/starsearch/releases/download/v0.1.0/starsearch-0.1.0-windows-amd64.zip'
$checksum = '5a9fdb7c26ae8ba8652f9688a33bacad5f0faf9da71679a4327bdb6caf0e17c2'
$checksumType = 'sha256'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url           = $url
  checksum      = $checksum
  checksumType  = $checksumType
}

Install-ChocolateyZipPackage @packageArgs
