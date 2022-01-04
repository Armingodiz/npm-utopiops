class UtopiopsCli < Formula
  desc "CLI tool for utopiops.com"
  homepage "https://github.com/Armingodiz/homebrew-utopiops"
  url "https://github.com/Armingodiz/homebrew-utopiops/releases/download/first-release/utopiops-cli-0.0.1.tar.gz"
  sha256 "4183af9568340ecfdfe8b9547dc6382db64c617194119efb15737e59a6482895"

  def install
    bin.install "utopiops-cli"
  end

  test do
    system "#{bin}/utopiops-cli", "--version"
  end
end
