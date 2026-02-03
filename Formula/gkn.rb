class Gkn < Formula
  desc "GitHub repository management CLI for local operations"
  homepage "https://github.com/TT-AIXion/github-kanri"
  version "0.1.2"
  url "https://github.com/TT-AIXion/github-kanri.git",
      tag: "v0.1.2",
      revision: "d6ed3b37cf24acd9146d8743ee5fdabb6c176c78"
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
