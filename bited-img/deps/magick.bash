ttf="$1"
txt="$2"
pango="$3"
out="$4"
fsz="$5"
bg="$6"
fg="$7"

magick -font "$ttf" -pointsize "$fsz" label:@"$txt" "$out".png

magick \
  -background "$bg" -fill "$fg" +antialias \
  -size "$(identify -ping -format '%wx%h' "$out".png)" \
  pango:@"$pango" \
  -bordercolor "$bg" -border "$fsz" \
  "$out".png
