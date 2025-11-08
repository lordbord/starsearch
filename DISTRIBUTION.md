# Starsearch Cross-Platform Distribution Guide

This guide covers how starsearch is distributed across multiple package managers for Linux, macOS, and Windows.

## 📦 Current Distribution Status

| Platform | Package Manager | Status | Installation Command |
|----------|----------------|--------|---------------------|
| Linux (Arch) | AUR | ✅ Live | `yay -S starsearch` |
| Linux/macOS | Homebrew | 🔄 Setup Required | `brew install lordbord/starsearch/starsearch` |
| Windows | Chocolatey | 🔄 Setup Required | `choco install starsearch` |
| All | GitHub Releases | ✅ Automated | Manual download |

## 🚀 Quick Start for v0.1.0 Release

### Step 1: Create GitHub Release

The binaries are already built in `dist/`. To create the release:

```bash
# Using GitHub CLI
gh release create v0.1.0 \
  dist/starsearch-0.1.0-linux-amd64.tar.gz \
  dist/starsearch-0.1.0-linux-arm64.tar.gz \
  dist/starsearch-0.1.0-darwin-amd64.tar.gz \
  dist/starsearch-0.1.0-darwin-arm64.tar.gz \
  dist/starsearch-0.1.0-windows-amd64.zip \
  dist/checksums.txt \
  --title "Starsearch v0.1.0" \
  --notes-file dist/release_notes.md
```

Or create manually at: https://github.com/lordbord/starsearch/releases/new

### Step 2: Set Up Homebrew Tap

```bash
# Create new repository on GitHub: lordbord/homebrew-starsearch

# Clone and set up
git clone https://github.com/lordbord/homebrew-starsearch.git
cd homebrew-starsearch
mkdir -p Formula
cp /home/jord/Projects/starsearch/homebrew/Formula/starsearch.rb Formula/
git add Formula/starsearch.rb
git commit -m "Add starsearch formula v0.1.0"
git push origin main
```

Test installation:
```bash
brew tap lordbord/starsearch
brew install starsearch
```

### Step 3: Submit to Chocolatey

From a Windows machine:

```powershell
cd C:\path\to\starsearch\chocolatey

# Test locally
choco pack
choco install starsearch -s . -y

# If successful, publish
choco push starsearch.0.1.0.nupkg --source https://push.chocolatey.org/
```

Note: First submission requires manual review (1-2 weeks).

## 🔄 Automated Releases (Future)

A GitHub Actions workflow has been created at `.github/workflows/release.yml` that will automatically:

1. Build cross-platform binaries
2. Create compressed archives
3. Generate checksums
4. Create GitHub release

**To use for future releases:**

```bash
# Create and push a new version tag
git tag v0.2.0
git push origin v0.2.0

# GitHub Actions will automatically create the release!
```

After the automated release:
- **Homebrew**: Update formula with new version/checksums and push
- **Chocolatey**: Update nuspec and chocolateyinstall.ps1, then push

## 📁 Repository Structure

```
starsearch/
├── dist/                          # Built binaries and archives
│   ├── checksums.txt
│   └── *.tar.gz, *.zip
├── homebrew/
│   ├── Formula/
│   │   └── starsearch.rb         # Homebrew formula
│   └── README.md                  # Homebrew setup guide
├── chocolatey/
│   ├── starsearch.nuspec          # Chocolatey metadata
│   ├── tools/
│   │   ├── chocolateyinstall.ps1 # Install script
│   │   └── VERIFICATION.txt      # Verification info
│   └── README.md                  # Chocolatey setup guide
├── .github/
│   └── workflows/
│       └── release.yml            # Automated release workflow
├── PKGBUILD                       # AUR package (already published)
└── DISTRIBUTION.md                # This file
```

## 📚 Detailed Guides

- **Homebrew**: See `homebrew/README.md`
- **Chocolatey**: See `chocolatey/README.md`
- **AUR**: See `AUR_RELEASE.md`

## 🔐 Checksums (v0.1.0)

```
17d95010ca7fd60125134c28c73eb952af7078efd93f68827b33638f1e005d76  starsearch-0.1.0-darwin-amd64.tar.gz
487e33b056c37d03ec112b499119b29502a35fa41d8632820a04e3f32b73a0f4  starsearch-0.1.0-darwin-arm64.tar.gz
29a99f5c0dd28f55305fc8ead30c0ba40c7286b7370de5af1ddb3a7cfea3e3cf  starsearch-0.1.0-linux-amd64.tar.gz
d98b4c800546fb301ff703c8b8e6b32e87063207e28333d2b2545df2d32cef61  starsearch-0.1.0-linux-arm64.tar.gz
5a9fdb7c26ae8ba8652f9688a33bacad5f0faf9da71679a4327bdb6caf0e17c2  starsearch-0.1.0-windows-amd64.zip
```

## 🎯 Next Steps

1. ✅ Build cross-platform binaries (DONE)
2. 🔄 Create GitHub Release v0.1.0
3. 🔄 Create and publish Homebrew tap
4. 🔄 Submit to Chocolatey Community Repository
5. 📝 Update README.md with installation instructions

## 🔗 Useful Links

- **GitHub Repository**: https://github.com/lordbord/starsearch
- **AUR Package**: https://aur.archlinux.org/packages/starsearch
- **Homebrew Tap**: https://github.com/lordbord/homebrew-starsearch (to be created)
- **Chocolatey Package**: https://community.chocolatey.org/packages/starsearch (pending submission)

## 💡 Tips

### Testing Releases Locally

Before creating a real release, test with a pre-release:

```bash
gh release create v0.1.0-rc1 \
  --prerelease \
  --title "Starsearch v0.1.0 Release Candidate" \
  dist/*.tar.gz dist/*.zip dist/checksums.txt
```

### Security Considerations

- All binaries built with `-ldflags="-s -w"` to strip debug info and reduce size
- Checksums provided for all releases
- HTTPS downloads from GitHub
- Windows binary may trigger SmartScreen on first downloads (normal for new software)

### Platform-Specific Notes

**macOS**: Users may need to allow the app in System Preferences > Security after first run

**Windows**: PowerShell execution policy may need adjustment for Chocolatey

**Linux**: AUR package compiles from source; Homebrew provides pre-built binaries
