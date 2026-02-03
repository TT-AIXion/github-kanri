class Gkn < Formula
  desc "GitHub repository management CLI for local operations"
  homepage "https://github.com/TT-AIXion/github-kanri"
  version "0.1.1"
  url "https://github.com/TT-AIXion/github-kanri.git",
      tag: "v0.1.1",
      revision: "0e1b5af4bc1470feee8bd210e0baa5a18c58f7d8"
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
