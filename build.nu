use scripts/bited-scale.nu

def main [src: path, out: path, --nerd, --release, --xs = [2 3]] {
  let name = $src | path parse | get stem
  let ttf = $out | path join $'($name).ttf'

  with-env {
    src: $src
    out: $out
    name: $name
    ttf: $ttf
    xs: $xs
    x_format: '{name}_{x}x'
  } {
    cp $src $out
    mk_vec
    if $nerd { mk_nerd }
    [1 ...$xs] | each { mk_x $in }

    if $release { mk_zip }
  }
}

def mk_vec [] {
  bitsnpicas convertbitmap -f 'ttf' -o $env.ttf $env.src

  [si0 fix so1]
  | each { open $'scripts/($in).py' }
  | str join "\n"
  | fontforge -c $in $env.ttf

  woff2_compress $env.ttf
}

def mk_nerd [] {
  nerd-font-patcher $env.ttf -out $env.out --careful -c
  nerd-font-patcher $env.ttf -out $env.out --careful -c -s
}

def mk_x [x = 1] {
  if $x <= 1 {
    mk_rest $env.name
  } else {
    let nm = { name: $env.name, x: $x } | format pattern $env.x_format
    bited-scale -x $x | save $'out/($nm).bdf'
    mk_rest $nm
  }
}

def mk_rest [name: string] {
  [si0 si1 fix so0]
  | each { open $'scripts/($in).py' }
  | str join "\n"
  | fontforge -c $in $env.src $'out/($env.name).' $env.name

  bdftopcf -o $'out/($env.name).pcf' $env.src
}

def mk_zip [] {
  let tag = git describe --tags --abbrev=0

  cp ['README.md' 'LICENSE' 'AUTHORS'] $env.out
  ^zip -r $'result/kirsch_($tag).zip' ($env.out | path join '*')
}