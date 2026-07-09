// Workspace house rule: no circled/enclosed alphanumerics (U+2460..U+24FF, e.g.
// the circled digits) as position labels. They break when content re-orders
// ("as in the second") and render inconsistently. Use markdown ordered lists,
// "(1)" inline, or names. Banned everywhere including README. Mirrors the
// retired tools/doc-lint.sh.
const CIRCLED = new RegExp("[\\u2460-\\u24FF]", "gu");

export default {
  names: ["ws-no-circled"],
  description:
    "Circled/enclosed alphanumerics are not allowed - use markdown lists, (1), or names",
  tags: ["ws", "style"],
  parser: "none",
  function: (params, onError) => {
    params.lines.forEach((line, index) => {
      CIRCLED.lastIndex = 0;
      let match;
      while ((match = CIRCLED.exec(line)) !== null) {
        onError({
          lineNumber: index + 1,
          detail:
            'Disallowed circled character "' +
            match[0] +
            '" - use "1." / "(1)" / names.',
          context: line.trim().slice(0, 40),
        });
      }
    });
  },
};
