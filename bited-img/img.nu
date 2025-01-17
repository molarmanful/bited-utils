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
      gen_gens
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
  if $env.chars in [null false ''] { return }

  $in
  | each { char -i $in }
  | if $env.accents { str replace -r '(\p{M})' ' $1' } else { }
  | chunks 48
  | each { str join ' ' }
  | str join "\n"
  | save -f (txt_path $env.chars)
}

def gen_map []: list<int> -> nothing {
  if $env.map in [null false ''] { return }

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
      | if $env.accents { str replace -r '(\p{M})' ' $1' } else { }
      | prepend $'U+($u)_ │'
      | str join ' '
    }
  | prepend [
      '          0 1 2 3 4 5 6 7 8 9 A B C D E F'
      '        ┌────────────────────────────────'
    ]
  | str join "\n"
  | save -f (txt_path $env.map)

  '5     1 8.'
  | repeat ($kvs | length)
  | prepend ['1' '8']
  | str join "\n"
  | save -f (txt_path $env.map "clr")
}

def txt_correct [] {
  ls_txts
  | where type == file
  | get 'name'
  | each {
      let f = $in
      open $f | str replace -m '\n$' '' | save -f $f
    }
}

def gen_gens [] {
  $env.gens
  | transpose k v
  | each {
      let k = $in.k
      let v = $in.v
      $v
      | each { txt_path $in | open }
      | str join "\n"
      | save -f (txt_path $k)
    }
  | str join "\n"
}

def gen_imgs [] {
  let tmpd = mktemp -d
  mkdir ($tmpd | path join 'fonts')
  let ttf = $tmpd | path join 'fonts/tmp.ttf'

  bitsnpicas convertbitmap -f 'ttf' -o $ttf $env.src
  let conf = mk_fc $tmpd

  print 'imgs...'
  ls_txts
  | where type == file
  | get 'name'
  | par-each {
      let txt = $in
      let path = $txt | path parse
      let out = (out_path $path.stem)
      let pango = $path | mk_pango
      with-env { FONTCONFIG_FILE: $conf } {
        bash (deps_path 'magick.bash') $ttf $txt $pango $out $env.font_size $env.clrs.bg $env.clrs.fg
      }
      print $' + ($path.stem)'
      rm -f $pango
    }

  rm -f $ttf
}

def mk_fc [dir: path]: nothing -> path {
  let d = $dir | path join 'fonts.conf'
  { dir: ($dir | path join 'fonts') }
  | format pattern (deps_path 'fonts.conf' | open)
  | save -f $d
  $d
}

def mk_pango [] {
  let tmp = mktemp
  {
    tag: span
    attributes: {
      font: $env.name
      size: $'($env.font_size * 3072 // 4)'
      background: $env.clrs.bg
      foreground: $env.clrs.fg
    }
    content: ($in | mk_content)
  } | to xml | save -f $tmp
  $tmp
}

def mk_content []: record -> list {
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
