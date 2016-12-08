# Generate submission zip file.

pushd report
pdflatex report
bibtex report
pdflatex report
pdflatex report
popd

cp report/report.pdf Brennan.Final.XPathFilesystem.pdf
go build
mv dpath dpath.linux
rm Brennan.Submission.XPathFilesystem.zip

zip Brennan.Submission.XPathFilesystem.zip \
    Brennan.Final.XPathFilesystem.pdf \
    README.EECS433.md \
    GUIDE.md \
    SYNTAX.md \
    axis.go \
    dpath.nex \
    dpath.nn.go \
    dpath.y \
    error.go \
    eval_test.go \
    item.go \
    lexer_test.go \
    lib.go \
    lib_test.go \
    main.go \
    parser_test.go \
    sequence.go \
    testutil.go \
    tree.go \
    util.go \
    y.go \
    dpath.linux

rm Brennan.Final.XPathFilesystem.pdf
rm dpath.linux
