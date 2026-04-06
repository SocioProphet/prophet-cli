class SourceosBootstrap < Formula
  desc "SourceOS bootstrap engine CLI"
  homepage "https://example.invalid/sourceos"
  url "https://example.invalid/sourceos-bootstrap-0.0.0.tar.gz"
  sha256 "TODO_PLACEHOLDER"
  license "Apache-2.0"

  # Source lineage:
  # Built from sourceos-sdk/cmd/sourceos-bootstrap.
  # ProphetCLI is a façade wrapper and not the engine implementation home.

  def install
    bin.install "sourceos-bootstrap"
  end

  test do
    system "#{bin}/sourceos-bootstrap", "doctor"
  end
end
