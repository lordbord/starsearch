# Homebrew Formula for Starsearch

This directory contains the Homebrew formula for starsearch, allowing installation via Homebrew on macOS and Linux.

## Setup Instructions

### 1. Create a Homebrew Tap Repository

Create a new GitHub repository called `homebrew-starsearch`:

```bash
# On GitHub, create a new repository: lordbord/homebrew-starsearch
```

### 2. Initialize and Push the Formula

```bash
# Clone your new tap repository
git clone https://github.com/lordbord/homebrew-starsearch.git
cd homebrew-starsearch

# Copy the formula
mkdir -p Formula
cp /path/to/starsearch/homebrew/Formula/starsearch.rb Formula/

# Commit and push
git add Formula/starsearch.rb
git commit -m "Add starsearch formula v0.1.0"
git push origin main
```

### 3. Test the Formula

Test the formula locally before publishing:

```bash
# Install from local file
brew install --build-from-source Formula/starsearch.rb

# Test the installation
starsearch --help

# Uninstall for re-testing
brew uninstall starsearch
```

### 4. Install from Tap (After Publishing)

Once pushed to GitHub, users can install with:

```bash
brew tap lordbord/starsearch
brew install starsearch
```

## Updating for New Releases

When releasing a new version:

1. **Build new binaries** and create GitHub release with new version tag

2. **Update the formula** with new version and checksums:
   ```ruby
   version "0.2.0"  # Update version
   
   # Update URLs
   url "https://github.com/lordbord/starsearch/releases/download/v0.2.0/..."
   
   # Update sha256 checksums from the new release
   sha256 "new_checksum_here"
   ```

3. **Test the updated formula**:
   ```bash
   brew reinstall --build-from-source Formula/starsearch.rb
   ```

4. **Commit and push**:
   ```bash
   git add Formula/starsearch.rb
   git commit -m "Update starsearch to v0.2.0"
   git push origin main
   ```

## Troubleshooting

### Formula doesn't install

- Verify all checksums match the released binaries
- Ensure GitHub release exists and binaries are publicly accessible
- Test download URLs manually: `curl -L <url>`

### Tap not found

- Ensure repository name follows pattern: `homebrew-<tapname>`
- Repository must be public
- Wait a few minutes after creating the repository

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)
- [Creating Taps](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
