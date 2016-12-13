const express  = require("express");
const sass     = require("node-sass");
const defaults = {base: "7DCED0", steps: 3};

const MAX_STEPS = 36;

function render(req, res) {
  let {query} = req;
  let {base, expanded} = query;
  let outputStyle = expanded ? "expanded" : "compressed";

  // attempt to use the user provided steps if they exist
  let steps = query.steps ? parseInt(query.steps, 10) : defaults.steps;

  if(base && /^[a-f0-9]{6}$/i.test(base) !== true) {
    res.status(422);
    return res.json({error: "must provide valid hex code"});
  }

  if(isNaN(steps) || parseInt(steps, 10) > MAX_STEPS) {
    res.status(422);
    return res.json({error: `step count must be between 1 and ${MAX_STEPS}`});
  }

  let config = {steps, base: base || defaults.base};

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
  @for $degree_step from 1 through ${config.steps} {
    $degrees: $degree_step * $step_amount;
    $base: adjust-hue(#${config.base}, $degrees);

    .fg-#{$degree_step} { color: $base; }
    .bg-#{$degree_step} { background-color: $base; }

    @for $brightness from 0 through 9 {
      $darker: darken($base, $brightness * 5);
      $lighter: lighten($base, $brightness * 5);
      .fg-#{$degree_step}--l#{$brightness + 1} { color: $lighter; }
      .fg-#{$degree_step}--d#{$brightness + 1} { color: $darker; }
      .bg-#{$degree_step}--l#{$brightness + 1} { background-color: $lighter; }
      .bg-#{$degree_step}--d#{$brightness + 1} { background-color: $darker; }
    }
  }`;

  sass.render({data, outputStyle}, finished)
}

function start(port) {
  let app = express();
  app.use(render);
  app.listen(port);
}

module.exports = {start};
