# imagelint

```
Usage of /bin/image-lint:
      --max-height int   縦の最大サイズ (default 1200)
      --max-width int    横の最大サイズ (default 800)
      --min-height int   縦の最小サイズ (default 50)
      --min-width int    横の最小サイズ (px) (default 50)
'glob-pattern 1' 'glob-pattern 2'
```

チェックが通れば、0が返る。

```
image-lint --max-width 800 "./**/*.png" "./**/*.jpg"
```
