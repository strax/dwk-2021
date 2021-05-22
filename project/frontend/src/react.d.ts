import "react"

declare module "react" {
  export interface SuspenseConfig {
    timeoutMs?: number
  }

  type StartTransition = (effect: () => void | Promise<void>) => void

  export function useDeferredValue<T>(value: T, config?: SuspenseConfig): T
  export function useTransition(config?: SuspenseConfig): [boolean, StartTransition]
}
