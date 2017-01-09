# krmp.cc

Inspired by the color css declarations seen in [materialize](http://materializecss.com/color.html), the goal of this project is to provide on-the-fly css classes who's only responsibility is to manipulate `color` and `background-color` properties.

## Using

Stylsheets are generated based on parameters provided to it in the url path as well as the query string. The "base" of the pallete should appear as the first part of the path after the hostname:

```
https://krmp.cc/:hex
```

To add the stylsheet to you application, simply add a `<link />` element pointed to the `krmp.cc` host:

```html
<link href="https://krmp.cc/6e6" rel="stylesheet" type="text/css" />
```

By default, this will generate css classes for a "_palette_" with the [default settings](https://github.com/dadleyy/krmp.cc/blob/master/krmp/request_runtime.go#L10-L15). Ultimately, the code generated looks like:

```css
.bg-1 { background-color: #6aa7d9; }
.fg-1 { color: #6aa7d9; }
.bg-1.darken-1 { background-color: #649dcc; }
.fg-1.darken-1 { color: #649dcc; }
```

These rules could then be applied to your application's markup through `class` attributes:

```html
<div class="bg-1 darken-1"></div>
```

### Customization

The following query string parameters can be added in order to customize the generated palette:

| Name | Type | Notes |
| ---- | ---- | ---- |
| steps | `number` | The amount of adjustments to make to the hue. |
| shades | `boolean` | If set to `false`, the generated css will not include `lighten`/`darken` classes for shades. |
| shade_min | `number` | Used to determine the minimum brightness of shades to generate. (0 - 100) |
| shade_max | `number` | Used to determine the maxium brightness of shades to generate. (0 - 100) |
| expanded | `boolean` | If set to `true`, the generated css will _not_ be minified. |
| noconflict | `string` | If provided, this string will be added to every class declaration. *This will also be used to modify the bower package name generated during downloads.* |

### Previewing

To preview any styesheet, add `/preview` to the url path:

* [`/428aa7/preview`](https://krmp.cc/428aa7/preview)
* [`/preview?base=428aa7`](https://krmp.cc/preview?base=428aa7)

### Downloading + Packaging

Even though downloading a local copy of the generated stylsheet can be done using curl:

```
$ curl -o generated.css https://krmp.cc/7a9eb1
```

The application also provides a `/download` (as well as a `/package`) endpoint that when used with a path-based hex string will create a `.tar.gz` archive containing the generated css as well as a `bower.json` file. This allows any generated stylsheet to be listed as a dependency of a project using [bower](https://bower.io) for their asset management:

```
$ bower i https://krmp.cc/656/package?noconflict=app --save
```

Which, would be inserted into the `bower.json` file of the project as:

```json
  "dependencies": {
    "krmp-app": "https://krmp.cc/656/package?noconflict=app"
  }
```
