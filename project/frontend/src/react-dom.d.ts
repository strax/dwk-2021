import "react-dom"

declare module "react-dom" {
  export const createRoot = unstable_createRoot
}
