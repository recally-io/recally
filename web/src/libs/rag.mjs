import * as mupdf from "./mupdf/mupdf.js";

export async function fileToDocs(file, url) {
  let data = [];
  let doc = mupdf.Document.openDocument(await file.arrayBuffer(), file.name);
  const numPages = doc.countPages();
  for (let i = 0; i < numPages; i++) {
    const page = doc.loadPage(i);
    const text = page.toStructuredText("preserve-whitespace").asText();

    const links = page.getLinks().map((link) => link.getURI());
    const metadata = {
      name: file.name,
      type: file.type,
      url: url,
      page: i,
      links: links,
    };
    console.log({ metadata: metadata, text: text });
    data.push({
      metadata: metadata,
      content: text,
    });
  }
  return data;
}
