use std repeat

export def main [cfg = 'bited-img.toml'] {
  if ($cfg | path type) != 'file' { return }

  open $cfg | transpose k v | each {
    let name = $in.k
    let v = $in.v

    with-env (
      deps_path 'unit.toml' | open
      | merge deep $v
      | upsert name $name
      | update src {|unit| { name: $name } | format pattern $unit.src }
    ) {

      let codes = get_codes
      $codes | gen_chars
      $codes | gen_map
      txt_correct
      gen_imgs

    }
  }
}

def get_codes []: nothing -> list<int> {
  open $env.src
  | lines
  | find -r '^\s*ENCODING\s+[^-]'
  | each { split words | get 1 | into int }
}

def gen_chars []: list<int> -> nothing {
  if $env.chars.out in [null false ''] { return }

  $in
  | each { char -i $in }
  | chunks $env.chars.width
  | each { str join ' ' }
  | str join "\n"
  | if $env.hide_accents { str replace -r -a '(\pM)' '.' } else { }
  | save -f (txt_path $env.chars.out)
}

def gen_map []: list<int> -> nothing {
  if $env.map.out in [null false ''] { return }

  let kvs = $in | group-by { $in // 16 } | transpose k v

  $kvs
  | each {
      let k = $in.k | into int
      let v = $in.v
      let u = $k
        | fmt
        | get upperhex
        | str substring 2..
        | fill -w 4 -a 'r' -c '0'

      $k * 16
      | $in..($in + 15)
      | each { if $in in $v { char -i $in } else { ' ' } }
      | prepend $'U+($u)_ │'
      | str join ' '
    }
  | prepend [
      '          0 1 2 3 4 5 6 7 8 9 A B C D E F'
      '        ┌────────────────────────────────'
    ]
  | str join "\n"
  | save -f (txt_path $env.map.out)

  $'($env.map.label_clrs.1)     ($env.map.label_clrs.0) ($env.map.border_clr).'
  | repeat ($kvs | length)
  | prepend [$env.map.label_clrs.0 $env.map.border_clr]
  | str join "\n"
  | if $env.hide_accents { str replace -r -a '(\pM)' '.' } else { }
  | save -f (txt_path $env.map.out "clr")
}

def txt_correct [] {
  ls_txts | each {
    let f = $in
    open $f | str replace -m '\n$' '' | save -f $f
  }
}

def gen_imgs [] {
  let tmpd = mktemp -d
  mkdir ($tmpd | path join 'fonts')
  let txtd = $tmpd | path join 'txts'
  mkdir $txtd
  let font = $tmpd | path join 'fonts/tmp.ttf'

  bitsnpicas convertbitmap -f 'ttf' -o $font $env.src
  let conf = gen_fc $tmpd

  print 'imgs...'
  ls_txts | par-each { gen_img $txtd $conf $font }

  print 'gens...'
  gen_gens $txtd $conf $font

  rm -rf $tmpd
}

def gen_gens [txtd: path, conf: path, font: path] {
  $env.gens
  | transpose k v
  | each {
      let k = $in.k
      let v = $in.v
      $v
      | each { txt_path $in | open }
      | str join "\n"
      | save -f (txt_path $k)

      $v
      | each {|x| $txtd | path join $x | open }
      | str join "\n"
      | save -f ($txtd | path join $k)
      txt_path $k | gen_img $txtd $conf $font --gen
    }
}

def gen_img [txtd: path, conf: path, font: path, --gen] {
  let txt = $in
  let path = $txt | path parse
  if ($path.stem in $env.gens) == $gen {
    let out = (out_path $path.stem)
    let pango = ($txtd | path join $path.stem)
    if not $gen { $path | gen_pango $pango }

    with-env { FONTCONFIG_FILE: $conf } {
      bash (deps_path 'magick.bash') $font $txt $pango $out $env.font_size $env.clrs.bg $env.clrs.fg
    }
    print $' + ($path.stem)'
  }
}

def gen_fc [tmpd: path]: nothing -> path {
  let d = $tmpd | path join 'fonts.conf'
  { dir: ($tmpd | path join 'fonts') }
  | format pattern (deps_path 'fonts.conf' | open)
  | save -f $d
  $d
}

def gen_pango [out: path] {
  {
    tag: span
    attributes: {
      font: $env.name
      size: $'($env.font_size * 3072 // 4)'
      background: $env.clrs.bg
      foreground: $env.clrs.fg
    }
    content: ($in | gen_content)
  } | to xml | save -f $out
}

def gen_content []: record -> list {
  let path = $in
  let txt = $path | path join
  let clr = $path | upsert extension 'clr' | path join
  if ($clr | path exists) {
    let clrs = $env.clrs.base | str join "\n"
    bited-pangogo --txt $txt --clr $clr --clrs $clrs | from json
  } else {
    [($txt | open)]
  }
}

def ls_txts [] {
  ls ...(glob ($env.txt_dir | path join '*.txt'))
  | where type == file
  | get 'name'
}

def deps_path [name: string]: nothing -> path {
  $env.FILE_PWD | path join 'deps' $name
}

def txt_path [name: string, ext = 'txt']: nothing -> path {
  $env.txt_dir | path join $'($name).($ext)'
}

def out_path [name: string]: nothing -> path {
  $env.out_dir | path join $name
}
