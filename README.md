# smeago
A Golang tool to generate sitemap xml for Server Side Rendering apps.

![all text](http://orig14.deviantart.net/9c7b/f/2012/268/9/d/doodle__cute__gollum_by_agathexu-d5fu2mf.jpg)

## Install

```
go get github.com/dbalduini/smeago
```

## Example usage

```
smeago -p 3000 -o "public/" -loc "http://example.com"
```

### Params

```
-h the host name to crawl
-p the host port to crawl
-loc the host to be prefixed with the paths in the sitemap
   example: -loc http://example.com
   <url>
    <loc>
      http://example.com/foo/bar
    </loc>
   </url>
-o the relative output directory for the sitemap.xml file
```
