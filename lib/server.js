const express  = require("express");
const sass     = require("node-sass");
const defaults = {base: "414141", steps: 10};

function render(req, res) {
  let {query} = req;
  let {base, steps} = query;

  if(base && /^[a-f0-9]{6}$/i.test(base) !== true) {
    res.status(422);
    return res.json({error: "must provide valid hex code"});
  }

  if(steps && /^[0-9]+$/i.test(steps) !== true) {
    res.status(422);
    return res.json({error: "must provide valid step count"});
  }

  if(360 < parseInt(steps, 10)) {
    res.status(422);
    return res.json({error: "step count exceeds maximum (360)"});
  }

  let config = {steps: steps || defaults.steps, base: base || defaults.base};

  function finished(error, result) {
    if(error) {
      res.status(422);
      return res.json({error});
    }

    res.set("Content-Type", "text/css");
    res.send(result.css);
  }

  let data = `
  $step_amount: 360 / ${steps};
  @for $degree_step from 0 through ${config.steps} {
    $degrees: $degree_step * $step_amount;
    $base: adjust-hue(#${config.base}, $degrees);

    .fg-#{$degree_step + 1} {
      color: $base;
    }

    .bg-#{$degree_step + 1} {
      background-color: $base;
    }

    @for $brightness from 0 through 9 {
      $darker: darken($base, $brightness * 5);
      $lighter: lighten($base, $brightness * 5);
      .fg-#{$degree_step + 1}--l#{$brightness + 1} { color: $lighter; }
      .fg-#{$degree_step + 1}--d#{$brightness + 1} { color: $darker; }
      .bg-#{$degree_step + 1}--l#{$brightness + 1} { background-color: $lighter; }
      .bg-#{$degree_step + 1}--d#{$brightness + 1} { background-color: $darker; }
    }
  }`;

  sass.render({data}, finished)
}

function start(port) {
  let app = express();
  app.use(render);
  app.listen(port);
}

module.exports = {start};
