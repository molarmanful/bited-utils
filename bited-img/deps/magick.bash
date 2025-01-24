ttf="$1"
txt="$2"
pango="$3"
out="$4"
fsz="$5"
bg="$6"
fg="$7"

magick \
  -background "$bg" -fill "$fg" +antialias \
  pango:@"$pango" \
  -bordercolor "$bg" -border "$fsz" \
  "$out".png
