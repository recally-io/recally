// FTS5 unicode61 does not segment CJK (plan §11.1). We pre-tokenize into a
// token stream stored in the FTS table; displayed text always comes from the
// original blocks, never from this tokenized form.
//
// Strategy: latin/digit runs stay whole; CJK chars become overlapping bigrams.
// TOKENIZER_VERSION must bump whenever this changes so indexes can be rebuilt.

const CJK = /[㐀-䶿一-鿿豈-﫿]/;

export function tokenizeForFts(text: string): string {
  const out: string[] = [];
  let cjkRun: string[] = [];
  let word = "";

  const flushWord = () => {
    if (word) {
      out.push(word.toLowerCase());
      word = "";
    }
  };

  const flushCjk = () => {
    if (cjkRun.length === 1) out.push(cjkRun[0] as string);

    for (let i = 0; i + 1 < cjkRun.length; i++) out.push(cjkRun[i]! + cjkRun[i + 1]!);
    cjkRun = [];
  };

  for (const ch of text) {
    if (CJK.test(ch)) {
      flushWord();
      cjkRun.push(ch);
    } else if (/[\p{L}\p{N}_]/u.test(ch)) {
      flushCjk();
      word += ch;
    } else {
      flushWord();
      flushCjk();
    }
  }

  flushWord();
  flushCjk();

  return out.join(" ");
}

// Build an FTS5 MATCH expression from user input without ever interpolating
// raw text into SQL. Tokens become an implicit AND of quoted terms.
export function buildFtsQuery(input: string): string {
  const tokens = tokenizeForFts(input).split(/\s+/).filter(Boolean).slice(0, 20);

  if (tokens.length === 0) return "";

  return tokens.map((t) => `"${t.replace(/"/g, '""')}"`).join(" ");
}
