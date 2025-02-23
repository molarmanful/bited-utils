import json

src, x, fsz, widths_json, out, otb_name, dfont_name = argv[1:]
x = int(x)
fsz = int(fsz)
widths = json.loads(widths_json)

f = open(src)
f.em = int(fsz) * x * 10
for g in f.glyphs():
    gn = g.glyphname if g.unicode < 0 else "U+%04X" % g.unicode
    g.width = widths[gn] * x * 10

{{.}}

f.selection.all()
f.correctDirection()
f.removeOverlap()
f.encoding = "UnicodeFull"
f.fontname = otb_name
f.generate(out, "otb")
f.fontname = dfont_name
f.generate(out + "dfont", "sbit")
