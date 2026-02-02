class Gkn < Formula
  desc "GitHub repository management CLI for local operations"
  homepage "https://github.com/TT-AIXion/github-kanri"
  version "0.0.0-main.20260202T131858Z.7f23932"
  url "https://github.com/TT-AIXion/github-kanri.git",
      tag: "v0.0.0-main.20260202T131858Z.7f23932",
      revision: "2d3ee34ce0228ada0b878cfd9794444dc6b17ee5"
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
