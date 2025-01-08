use std log

def main [src: path, cfg = 'bited-img.yaml'] {
  let cfg_val = if ($cfg | path type) == 'file' { open $cfg } else { {} }
  with-env (
    {
      out: 'img'
      accents: false
      txt_dir: 'txt'
      txt: {
        chars: 'chars'
        map: 'map'
        sample: 'sample'
        all: 'all'
      }
      samples: []
      bg: '#ffffff'
      fg: '#000000'
    }
    | merge deep $cfg_val
    | merge { src: $src }
  ) {
    let codes = get_codes
    $codes | gen_chars
    $codes | gen_map
    txt_correct
    gen_samples
    gen_all
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
  $in
  | each { char -i $in }
  | if $env.accents { str replace -r '(\p{M})' ' $1' } else { }
  | chunks 48
  | each { str join ' ' }
  | str join "\n"
  | save -f (txt_path $env.txt.chars)
}

def gen_map []: list<int> -> nothing {
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
  | save -f (txt_path $env.txt.map)
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

def gen_samples [] {
  $env.samples
  | each {|a| txt_path $a | open }
  | str join "\n"
  | save -f (txt_path $env.txt.sample)
}

def gen_all [] {
  (txt_path $env.txt.sample) + "\n" + (txt_path $env.txt.sample)
  | save -f (txt_path $env.txt.all)
}

def gen_imgs [] {
  let ttf = mktemp --suffix .ttf
  let tmpd = mktemp -d
  let tmp = $tmpd | path join 'tmp.ttf'
  bitsnpicas convertbitmap -f 'ttf' -o $tmp $env.src
  mv $tmp $ttf
  rm -rf tmpd

  log info 'imgs...'
  ls $env.txt_dir
  | where type == file
  | get 'name'
  | par-each {
      let stem = $in | path parse | get stem
      sh scripts/magick.sh -- $ttf 16 $in ($env.out | path join $stem) $env.bg $env.fg
      log info $' + ($stem)'
    }

  rm -f $ttf
}

def txt_path [name: string] {
  $env.txt_dir | path join $'($name).txt'
}
