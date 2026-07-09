// Workspace house rule: no pictographic/decorative emoji in Markdown.
// Plain check/cross marks (U+2713 checkmark, U+2717 ballot X) stay allowed as
// functional table/status markers. README / README.<lang> files are exempt
// (GitHub-readme genre). Ranges mirror the retired tools/doc-lint.sh.
const EMOJI = new RegExp(
  "[\\u2600-\\u26FF\\u2700-\\u2712\\u2714-\\u2716\\u2718-\\u27BF" +
    "\\u2B00-\\u2BFF\\uFE00-\\uFE0F\\u{1F000}-\\u{1FAFF}]",
  "gu",
);

export default {
  names: ["ws-no-emoji"],
  description:
    "Pictographic/decorative emoji are not allowed in technical docs (plain check/cross marks ok; README exempt)",
  tags: ["ws", "style", "emoji"],
  parser: "none",
  function: (params, onError) => {
    if (/(^|\/)README(\.[a-z]{2})?\.md$/i.test(params.name)) return;
    params.lines.forEach((line, index) => {
      EMOJI.lastIndex = 0;
      let match;
      while ((match = EMOJI.exec(line)) !== null) {
        onError({
          lineNumber: index + 1,
          detail:
            'Disallowed emoji "' +
            match[0] +
            '" - use text or **bold** (plain check/cross marks are allowed).',
          context: line.trim().slice(0, 40),
        });
      }
    });
  },
};
