export def main [src: path, cfg = 'bited-img.toml'] {
  let cfg_val = if ($cfg | path type) == 'file' { open $cfg } else { {} }
  with-env (
    deps_path 'bited-img.toml' | open
    | merge deep $cfg_val
    | upsert src $src
  ) {
    let codes = get_codes
    $codes | gen_chars
    $codes | gen_map
    txt_correct
    gen_gens
    gen_imgs
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

  $in
  | group-by { $in // 16 }
  | transpose k v
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
}

def txt_correct [] {
  ls $env.txt_dir
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
  let ttf = mktemp --suffix .ttf
  let tmpd = mktemp -d
  let tmp = $tmpd | path join 'tmp.ttf'

  bitsnpicas convertbitmap -f 'ttf' -o $tmp $env.src
  mv $tmp $ttf
  rm -rf tmpd

  print 'imgs...'
  ls $env.txt_dir
  | where type == file
  | get 'name'
  | par-each {
      let stem = $in | path parse | get stem
      let out = (out_path $stem)
      bash (deps_path 'magick.bash') $ttf $in $out $env.font_size $env.bg $env.fg
      print $' + ($stem)'
    }

  rm -f $ttf
}

def deps_path [name: string]: nothing -> path {
  $env.FILE_PWD | path join 'deps' $name
}

def txt_path [name: string]: nothing -> path {
  $env.txt_dir | path join $'($name).txt'
}

def out_path [name: string]: nothing -> path {
  $env.out_dir | path join $name
}
