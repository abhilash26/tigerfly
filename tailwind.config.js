// tailwind.config.js
import daisyui from "daisyui";

export default {
  content: ["./views/**/*.{tmpl,html}"],
  theme: {
    extend: {},
  },
  plugins: [daisyui],
  daisyui: {
    themes: ["cupcake"],
  },
};
