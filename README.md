# krmp.cc

Inspired by the color css declarations seen in [materialize](http://materializecss.com/color.html), the goal of this project is to provide on-the-fly css classes who's only responsibility is to manipulate `color` and `background-color` properties.

## Using

To generate a css stylesheet for use in your application, include a `link` tag pointed to the `krmp.cc` domain:

```html
<link href="https://krmp.cc" rel="stylesheet" type="text/css" />
```

By default, this will generate css classes for a "_palette_" with 3 hue steps and based on the `#6aa7d9` hex color. Ultimately, the code generated looks like:

```css
.bg-1 { background-color: #6aa7d9; }
.fg-1 { color: #6aa7d9; }
.bg-1.darken-1 { background-color: #649dcc; }
.fg-1.darken-1 { color: #649dcc; }
```

Which means, in your application's markup, you'd apply the colors by:

```html
<div class="bg-1 darken-1">
</div>
```

### Customization

The default parameters can be overidden by providing query parameters to the application:

```
https://krmp.cc?steps=35&base=aea
```

| Name | Type | Notes |
| ---- | ---- | ---- |
| base | `hex` | The color from which to start the palette. Should not contain leading `#`. |
| steps | `number` | The amount of adjustments to make to the hue. |
| shades | `boolean` | If set to `false`, the generated css will not include `lighten`/`darken` classes for shades. |
| shade_min | `number` | Used to determine the minimum brightness of shades to generate. (0 - 100) |
| shade_max | `number` | Used to determine the maxium brightness of shades to generate. (0 - 100) |
