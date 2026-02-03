class Gkn < Formula
  desc "GitHub repository management CLI for local operations"
  homepage "https://github.com/TT-AIXion/github-kanri"
  version "0.0.0-main.20260203T020046Z.95b18f0"
  url "https://github.com/TT-AIXion/github-kanri.git",
      tag: "v0.0.0-main.20260203T020046Z.95b18f0",
      revision: "95b18f084c643a5f0839b3e0eef8580a51639d1f"
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
