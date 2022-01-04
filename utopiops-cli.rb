class UtopiopsCli < Formula
  desc "CLI tool for utopiops.com"
  homepage "https://github.com/Armingodiz/homebrew-utopiops"
  url "https://github.com/Armingodiz/homebrew-utopiops/releases/download/first-release/utopiops-cli-0.0.1.tar.gz"
  sha256 "9cf3b4ea854858a910bb13436165d1c8fd96ec9cfa5748cb41bb1783d7aaf9e0"

  def install
    bin.install "utopiops-cli"
  end

  test do
    system "#{bin}/utopiops-cli", "--version"
  end
end
