import * as mupdf from "mupdf";
export function loadPDF(data) {
  return new mupdf.PDFDocument(data);
}
export function drawPageAsPNG(document, pageNumber, dpi) {
  const page = document.loadPage(pageNumber);
  const zoom = dpi / 72;
  return page
    .toPixmap([zoom, 0, 0, zoom, 0, 0], mupdf.ColorSpace.DeviceRGB)
    .asPNG();
}
export function drawPageAsHTML(document, pageNumber, id) {
  return document.loadPage(pageNumber).toStructuredText().asHTML(id);
}
export function drawPageAsSVG(document, pageNumber) {
  const page = document.loadPage(pageNumber);
  const buffer = new mupdf.Buffer();
  const writer = new mupdf.DocumentWriter(buffer, "svg", "");
  const device = writer.beginPage(page.getBounds());
  page.run(device, mupdf.Matrix.identity);
  device.close();
  writer.endPage();
  return buffer.asString();
}
export function getPageText(document, pageNumber) {
  return document.loadPage(pageNumber).toStructuredText().asText();
}
export function searchPageText(
  document,
  pageNumber,
  searchString,
  maxHits = 500,
) {
  return document
    .loadPage(pageNumber)
    .toStructuredText()
    .search(searchString, maxHits);
}
