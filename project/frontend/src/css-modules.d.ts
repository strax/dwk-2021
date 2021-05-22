declare module "*.module.css" {
  declare const module: Record<string, string>
  export default module
}
