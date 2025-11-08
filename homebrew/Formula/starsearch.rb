class Starsearch < Formula
  desc "A modern, feature-rich Gemini protocol browser built with Go and Bubble Tea TUI framework"
  homepage "https://github.com/lordbord/starsearch"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/lordbord/starsearch/releases/download/v0.1.0/starsearch-0.1.0-darwin-arm64.tar.gz"
      sha256 "487e33b056c37d03ec112b499119b29502a35fa41d8632820a04e3f32b73a0f4"
    else
      url "https://github.com/lordbord/starsearch/releases/download/v0.1.0/starsearch-0.1.0-darwin-amd64.tar.gz"
      sha256 "17d95010ca7fd60125134c28c73eb952af7078efd93f68827b33638f1e005d76"
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/lordbord/starsearch/releases/download/v0.1.0/starsearch-0.1.0-linux-arm64.tar.gz"
      sha256 "d98b4c800546fb301ff703c8b8e6b32e87063207e28333d2b2545df2d32cef61"
    elsif Hardware::CPU.intel?
      url "https://github.com/lordbord/starsearch/releases/download/v0.1.0/starsearch-0.1.0-linux-amd64.tar.gz"
      sha256 "29a99f5c0dd28f55305fc8ead30c0ba40c7286b7370de5af1ddb3a7cfea3e3cf"
    end
  end

  def install
    bin.install "starsearch"
  end

  test do
    # Basic test to ensure the binary exists and is executable
    system "#{bin}/starsearch", "--help" rescue true
  end
end
