# Chocolatey Package for Starsearch

This directory contains the Chocolatey package definition for starsearch, allowing installation via Chocolatey on Windows.

## Prerequisites

- Windows machine or VM for testing
- [Chocolatey](https://chocolatey.org/install) installed
- Chocolatey account (for publishing to community repository)

## Package Structure

```
chocolatey/
├── starsearch.nuspec          # Package metadata
├── tools/
│   ├── chocolateyinstall.ps1  # Installation script
│   └── VERIFICATION.txt        # Verification instructions
└── README.md                   # This file
```

## Testing the Package Locally

### 1. Pack the Package

From the `chocolatey` directory:

```powershell
choco pack
```

This creates `starsearch.0.1.0.nupkg`

### 2. Install Locally

```powershell
# Install from local package
choco install starsearch -s . -y

# Test the installation
starsearch --help

# Uninstall for re-testing
choco uninstall starsearch -y
```

## Publishing to Chocolatey Community Repository

### 1. Create Chocolatey Account

Register at: https://community.chocolatey.org/account/Register

### 2. Get Your API Key

Find your API key at: https://community.chocolatey.org/account

### 3. Set API Key

```powershell
choco apikey --key YOUR_API_KEY --source https://push.chocolatey.org/
```

### 4. Push Package

```powershell
choco push starsearch.0.1.0.nupkg --source https://push.chocolatey.org/
```

### 5. Wait for Approval

- First-time packages require manual review (can take 1-2 weeks)
- You'll receive email updates about the review process
- Once approved, updates are typically automated

## Alternative: Host Your Own Feed

If you don't want to publish to the community repository, you can host your own feed:

### Option A: Simple File Share

```powershell
# Install from local directory
choco install starsearch -s C:\path\to\chocolatey\packages -y

# Or from network share
choco install starsearch -s \\server\share\packages -y
```

### Option B: Web-Based Feed

Host the `.nupkg` files on your own web server and add as a source:

```powershell
choco source add -n=myrepo -s="https://myserver.com/packages/"
choco install starsearch -s myrepo -y
```

## Updating for New Releases

When releasing a new version:

### 1. Update Package Files

**starsearch.nuspec**:
```xml
<version>0.2.0</version>
<releaseNotes>https://github.com/lordbord/starsearch/releases/tag/v0.2.0</releaseNotes>
```

**tools/chocolateyinstall.ps1**:
```powershell
$url = 'https://github.com/lordbord/starsearch/releases/download/v0.2.0/starsearch-0.2.0-windows-amd64.zip'
$checksum = 'NEW_CHECKSUM_HERE'
```

**tools/VERIFICATION.txt**:
- Update URL and checksum

### 2. Test, Pack, and Push

```powershell
# Test locally first
choco pack
choco install starsearch -s . -y --force

# If successful, push to Chocolatey
choco push starsearch.0.2.0.nupkg --source https://push.chocolatey.org/
```

## Troubleshooting

### Package fails to install

- Verify the download URL is accessible
- Confirm checksum matches the released binary
- Check PowerShell execution policy: `Get-ExecutionPolicy`

### Checksum mismatch

Get the correct checksum:
```powershell
Get-FileHash starsearch-0.1.0-windows-amd64.zip -Algorithm SHA256
```

### Review rejected

- Review feedback from Chocolatey moderators
- Common issues:
  - Missing or incorrect verification
  - Security concerns in install script
  - Metadata issues in nuspec

## Resources

- [Chocolatey Documentation](https://docs.chocolatey.org/)
- [Creating Packages](https://docs.chocolatey.org/en-us/create/create-packages)
- [Package Validator Rules](https://docs.chocolatey.org/en-us/community-repository/moderation/package-validator/rules/)
- [Community Repository Submission](https://docs.chocolatey.org/en-us/community-repository/)
