import { Temporal } from "proposal-temporal"

if (!Reflect.has(globalThis, "Temporal")) {
  Reflect.set(globalThis, "Temporal", Temporal)
}
