class Gkn < Formula
  desc "GitHub repository management CLI for local operations"
  homepage "https://github.com/TT-AIXion/github-kanri"
  version "0.0.0-main.20260203T013101Z.42c3a83"
  url "https://github.com/TT-AIXion/github-kanri.git",
      tag: "v0.0.0-main.20260203T013101Z.42c3a83",
      revision: "42c3a83d1a690f6ad2521bf40ecfd2e4331ee1ab"
  license "MIT"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X github.com/TT-AIXion/github-kanri/cmd/gkn.Version=#{version}"
    system "go", "build", *std_go_args(ldflags: ldflags), "./cmd/gkn"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/gkn version")
  end
end
