class Gkn < Formula
  desc "GitHub repository management CLI for local operations"
  homepage "https://github.com/TT-AIXion/github-kanri"
  version "0.0.0-main.20260202T130646Z.7e8e2b1"
  url "https://github.com/TT-AIXion/github-kanri.git",
      tag: "v0.0.0-main.20260202T130646Z.7e8e2b1",
      revision: "b5289a0e9cddcd0f2f54503aa35631af43c6b4b1"
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
