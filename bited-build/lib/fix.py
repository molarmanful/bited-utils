src, fsz, name = argv[1:]
fsz = int(fsz)

f = open(argv[1])
f.fontname = name
f.em = fsz * 10

{{.}}

f.selection.all()
f.correctDirection()
f.removeOverlap()
f.encoding = "UnicodeFull"
f.generate(argv[1])
