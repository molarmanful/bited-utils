f = open(argv[1])
{{.}}
f.selection.all()
f.correctDirection()
f.removeOverlap()
f.encoding = "UnicodeFull"
f.generate(argv[1])
