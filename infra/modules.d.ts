// Ambient module declarations for non-TS imports shared by the root and app
// TypeScript programs.
declare module "*.md?raw" {
  const content: string;
  export default content;
}
