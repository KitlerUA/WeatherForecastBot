.PHONY: dependencies
dependencies:
	echo "Installing dependencies"
	glide install

.PHONY: code-quality
code-quality:
	gometalinter --vendor --tests --skip=mock --exclude='_gen.go' --disable=gotype --disable=errcheck --disable=gas --disable=dupl --deadline=1500s --checkstyle --sort=linter ./... > static-analysis.xml

install-helpers:
	echo "Installing GoMetaLinter"
	go get -u github.com/alecthomas/gometalinter
	echo "Installing linters"
	gometalinter --install
	echo "Installing Glide"
	go get -u github.com/Masterminds/glide
