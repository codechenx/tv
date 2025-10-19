class Tv < Formula
  desc "Fast, feature-rich CSV/TSV/delimited file viewer for the command line"
  homepage "https://github.com/codechenx/tv"
  version "0.7.1"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/codechenx/tv/releases/download/v0.7.1/tv_0.7.1_Darwin_arm64.tar.gz"
      sha256 "7f449350640493640fe1b16be6aa9f4417c978b59490ca38dc65d1fdd0df6412"
    else
      url "https://github.com/codechenx/tv/releases/download/v0.7.1/tv_0.7.1_Darwin_x86_64.tar.gz"
      sha256 "9264c13e07f4a37eec112abf38144bd19f0106cbfbe91234975d9eaa77523cec"
    end
  end

  def install
    bin.install "tv"
  end

  test do
    assert_match "v0.7.1", shell_output("#{bin}/tv -v")
  end
end
