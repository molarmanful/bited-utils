export def main [cfg = 'bited-build.toml', --nerd, --release] {
  if ($cfg | path type) != 'file' { return }

  open $cfg | transpose k v | each {
    let name = $in.k
    let v = $in.v

    with-env (
      deps_path 'unit.toml' | open
      | merge deep $v
      | upsert name $name
    ) {
      let src = ({ name: $name } | format pattern $env.src)
      let stem = $src | path parse | get stem
      with-env {
        src: $src
        stem: $stem
        ttf: (out_path $'($stem).ttf')
      } {
        if not ($env.out_dir | path exists) { mkdir $env.out_dir }
        cp $env.src $env.out_dir

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
  nerd-font-patcher $env.ttf -out $env.out_dir --careful -c
  nerd-font-patcher $env.ttf -out $env.out_dir --careful -c -s
}

def mk_x [x = 1] {
  if $x <= 1 {
    mk_rest $env.stem $env.name
  } else {
    let name = { name: $env.name, x: $x } | format pattern $env.x_format
    let stem = $'($env.name)_($x)x'
    open $env.src
    | bited-scale -n $x --name $name
    | save -f (out_path $'($stem).bdf')
    mk_rest $stem $name
  }
}

def mk_rest [stem: string, name: string] {
  [si0 si1 fix so0]
  | each { deps_path $'($in).py' | open }
  | insert 2 (ttfix)
  | str join "\n"
  | fontforge -c $in $env.src (out_path $'($stem).') $name

  bdftopcf -o (out_path $'($stem).pcf') $env.src
}

def mk_zip [] {
  let tag = open $env.verfile

  $env.zip_includes | each { cp $in $env.out_dir }
  ^zip -r (out_path $'kirsch_($tag).zip') (out_path '*')
}

def deps_path [name: string] {
  $env.FILE_PWD | path join 'deps' $name
}

def out_path [name: string] {
  $env.out_dir | path join $name
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
