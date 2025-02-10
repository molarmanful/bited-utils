import json

src, x, fsz, widths_json, out, name = argv[1:]
x = int(x)
fsz = int(fsz)
widths = json.loads(widths_json)

f = open(src)
f.fontname = name
f.em = int(fsz) * x * 10
for g in f.glyphs():
    name = g.glyphname if g.unicode < 0 else "U+%04X" % g.unicode
    g.width = widths[name] * x * 10

{{.}}

f.selection.all()
f.correctDirection()
f.removeOverlap()
f.encoding = "UnicodeFull"
f.generate(out, "otb")
f.generate(out + "dfont", "sbit")
