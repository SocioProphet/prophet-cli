cask "sourceos-bootstrap" do
  version "0.0.0"
  sha256 :no_check

  url "https://example.invalid/sourceos-bootstrap-0.0.0.pkg"
  name "sourceos-bootstrap"
  desc "Scaffold cask for sourceos-bootstrap distribution"
  homepage "https://example.invalid/sourceos"

  pkg "sourceos-bootstrap.pkg"
end
