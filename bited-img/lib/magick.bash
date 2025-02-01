magick \
  -background "{{ .Bg }}" -fill "{{ .Fg }}" +antialias \
  pango:@"{{ .Pango }}" \
  -bordercolor "{{ .Bg }}" -border "{{ .FontSize }}x" \
  -splice "x{{ .FontSize }}" \
  -strip "{{ .Out }}.png"
