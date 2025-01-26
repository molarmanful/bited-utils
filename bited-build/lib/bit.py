f = open(argv[1])
f.fontname = argv[3]
{{.}}
f.selection.all()
f.correctDirection()
f.removeOverlap()
f.encoding = "UnicodeFull"
f.generate(argv[2], "otb")
f.generate(argv[2] + "dfont", "sbit")
