ttf="$1"
txt="$2"
out="$3"
fsz="$4"
bg="$5"
fg="$6"

magick \
  -background "$bg" -fill "$fg" \
  -font "$ttf" -pointsize "$fsz" +antialias \
  label:@"$txt" \
  -bordercolor "$bg" -border "$fsz" \
  "$out".png
