/// <reference types="node" resolution-mode="require"/>
import * as mupdf from "mupdf";
export declare function loadPDF(
  data: Buffer | ArrayBuffer | Uint8Array,
): mupdf.PDFDocument;
export declare function drawPageAsPNG(
  document: mupdf.PDFDocument,
  pageNumber: number,
  dpi: number,
): Uint8Array;
export declare function drawPageAsHTML(
  document: mupdf.PDFDocument,
  pageNumber: number,
  id: number,
): string;
export declare function drawPageAsSVG(
  document: mupdf.PDFDocument,
  pageNumber: number,
): string;
export declare function getPageText(
  document: mupdf.PDFDocument,
  pageNumber: number,
): string;
export declare function searchPageText(
  document: mupdf.PDFDocument,
  pageNumber: number,
  searchString: string,
  maxHits?: number,
): [number, number, number, number, number, number, number, number][][];
