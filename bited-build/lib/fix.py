src, fsz = argv[1:]
fsz = int(fsz)

f = open(argv[1])
f.em = fsz * 10

{{.}}

f.selection.all()
f.correctDirection()
f.removeOverlap()
f.encoding = "UnicodeFull"
f.generate(argv[1])
