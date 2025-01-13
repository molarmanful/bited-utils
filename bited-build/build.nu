export def main [cfg = 'bited-build.toml', --nerd, --release] {
  let cfg_val = if ($cfg | path type) != 'file' { return }

  open $cfg | transpose k v | each {
    let name = $in.k
    let v = $in.v

    with-env (
      deps_path 'unit.toml' | open
      | merge deep $v
      | upsert name $name
    ) {
      with-env {
        src: ({ name: $name } | format pattern $env.src)
        ttf: (out_path $'($name).ttf')
      } {
        if not ($env.out | path exists) { mkdir $env.out }
        cp $env.src $env.out

        mk_vec
        if $nerd { mk_nerd }

        [1 ...$env.xs]
        | each { into int --signed }
        | uniq
        | filter { $in > 0 }
        | each { mk_x $in }

        if $release { mk_zip }
      }
    }

    print $'($name) built.'
  }
  return
}

def mk_vec [] {
  bitsnpicas convertbitmap -f 'ttf' -o $env.ttf $env.src

  [si0 fix so1]
  | each { deps_path $'($in).py' | open }
  | insert 1 (ttfix)
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
    open $env.src
    | bited-scale -n $x
    | save -f (out_path $'($env.name)_($x)x.bdf')
    mk_rest $nm
  }
}

def mk_rest [name: string] {
  [si0 si1 fix so0]
  | each { deps_path $'($in).py' | open }
  | insert 2 (ttfix)
  | str join "\n"
  | fontforge -c $in $env.src (out_path $'($env.name).') $env.name

  bdftopcf -o (out_path $'($env.name).pcf') $env.src
}

def mk_zip [] {
  let tag = git describe --tags --abbrev=0

  ['README.md' 'LICENSE' 'AUTHORS'] | each { cp $in $env.out }
  ^zip -r (out_path $'kirsch_($tag).zip') $env.out
}

def deps_path [name: string] {
  $env.FILE_PWD | path join 'deps' $name
}

def out_path [name: string] {
  $env.out | path join $name
}

def ttfix []: nothing -> string {
  $env.sfnt
  | transpose k v
  | each {
      [
        'f.appendSFNTName('
        ([$env.sfnt_lang $in.k $in.v] | each { to json -r } | str join ',')
        ')'
      ]
      | str join
    }
  | append $env.ttfix
  | str join "\n"
}
